package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

type mockedListsService struct {
	mock.Mock
}

func (us *mockedListsService) AddUserList(userID string, l *models.List) (string, error) {
	args := us.Called(userID, l)
	return args.String(0), args.Error(1)
}

func (us *mockedListsService) RemoveUserList(id string, userID string) error {
	args := us.Called(id, userID)
	return args.Error(0)
}

func (us *mockedListsService) UpdateUserList(id string, userID string, l *models.List) error {
	args := us.Called(id, userID, l)
	return args.Error(0)
}

func (us *mockedListsService) GetSingleUserList(i string, u string, l *models.List) error {
	args := us.Called(i, u, l)
	return args.Error(0)
}

func (us *mockedListsService) GetUserLists(u string, r *[]models.GetListsResultDto) error {
	args := us.Called(u, r)
	return args.Error(0)
}

func TestLists(t *testing.T) {
	testListsSrv := new(mockedListsService)

	testSrvProvider := new(mockedServiceProvider)

	jwtInfo := models.JwtClaimsInfo{
		UserID: "id",
	}

	t.Run("GET returns an okResult when there is no error", func(t *testing.T) {
		data := models.SampleGetListsResultDto()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		testListsSrv.On("GetUserLists", jwtInfo.UserID, &[]models.GetListsResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]models.GetListsResultDto)
			*arg = data
		})

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		request = addUserIDToContext(jwtInfo.UserID, request)

		got := ListsHandler(request, testSrvProvider)

		want := okResult{data, http.StatusOK}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("GET returns an errorResult with the service error when the query fails", func(t *testing.T) {
		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		err := errors.New("wadus")
		testListsSrv.On("GetUserLists", jwtInfo.UserID, &[]models.GetListsResultDto{}).Return(err).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		request = addUserIDToContext(jwtInfo.UserID, request)

		got := ListsHandler(request, testSrvProvider)

		errorResult, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorResult.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("GET WITH AN ID returns an errorResult with the service error when the query fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		err := errors.New("wadus")
		testListsSrv.On("GetSingleUserList", id, jwtInfo.UserID, &models.List{}).Return(errors.New("wadus")).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		request = addUserIDToContext(jwtInfo.UserID, request)

		got := ListsHandler(request, testSrvProvider)

		errorResult, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorResult.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("GET WITH AN ID returns an okResult with a single list", func(t *testing.T) {
		data := models.SampleListSlice()[0]

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		testListsSrv.On("GetSingleUserList", data.ID, jwtInfo.UserID, &models.List{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(2).(*models.List)
			*arg = data
		})

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+data.ID, nil)
		request = addUserIDToContext(jwtInfo.UserID, request)

		got := ListsHandler(request, testSrvProvider)

		want := okResult{data, http.StatusOK}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("POST returns an okResult when there is no error", func(t *testing.T) {
		listDto := listDtoToCreate()

		data := listDto.ToList()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()

		testListsSrv.On("AddUserList", jwtInfo.UserID, &data).Return("id", nil).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request = addUserIDToContext(jwtInfo.UserID, request)
		request.Header.Set("Content-type", "application/json")

		got := ListsHandler(request, testSrvProvider)
		want := okResult{"id", http.StatusCreated}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("POST with invalid body should return an errorResult with a BadRequestError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", strings.NewReader("wadus"))
		request = addUserIDToContext(jwtInfo.UserID, request)

		got := ListsHandler(request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		badReqErr, isInvalidBodyError := errorRes.err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isInvalidBodyError, "should be a bad request error")
		assert.Equal(t, "Invalid body", badReqErr.Error())

		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("POST returns an errorResult with the service error when the insert fails", func(t *testing.T) {
		listDto := listDtoToCreate()

		data := listDto.ToList()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		err := errors.New("wadus")
		testListsSrv.On("AddUserList", jwtInfo.UserID, &data).Return("", err).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request = addUserIDToContext(jwtInfo.UserID, request)
		request.Header.Set("Content-type", "application/json")

		got := ListsHandler(request, testSrvProvider)
		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorRes.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("DELETE returns an errorResult with the service erorr when the delete fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		err := errors.New("wadus")
		testListsSrv.On("RemoveUserList", id, jwtInfo.UserID).Return(err).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		request = addUserIDToContext(jwtInfo.UserID, request)

		got := ListsHandler(request, testSrvProvider)
		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorRes.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("DELETE returns an okResult when there is no error", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		testListsSrv.On("RemoveUserList", id, jwtInfo.UserID).Return(nil).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		request = addUserIDToContext(jwtInfo.UserID, request)

		got := ListsHandler(request, testSrvProvider)
		want := okResult{nil, http.StatusNoContent}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("PUT with invalid body should return an errorResult with a BadRequestError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/lists", strings.NewReader("wadus"))
		request = addUserIDToContext(jwtInfo.UserID, request)

		got := ListsHandler(request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		badReqErr, isInvalidBodyError := errorRes.err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isInvalidBodyError, "should be a bad request error")
		assert.Equal(t, "Invalid body", badReqErr.Error())

		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("PUT returns an errorResult with the service error when the update fails", func(t *testing.T) {
		listDto := listDtoToUpdate()
		id := bson.NewObjectId().Hex()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		err := errors.New("wadus")
		testListsSrv.On("UpdateUserList", id, jwtInfo.UserID, mock.Anything).Return(err).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request = addUserIDToContext(jwtInfo.UserID, request)
		request.Header.Set("Content-type", "application/json")

		got := ListsHandler(request, testSrvProvider)
		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorRes.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("PUT returns an okResult when there is no error", func(t *testing.T) {
		listDto := listDtoToUpdate()

		data := listDto.ToList()

		id := bson.NewObjectId().Hex()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		testListsSrv.On("UpdateUserList", id, jwtInfo.UserID, &data).Return(nil).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request = addUserIDToContext(jwtInfo.UserID, request)
		request.Header.Set("Content-type", "application/json")

		got := ListsHandler(request, testSrvProvider)
		want := okResult{data, http.StatusOK}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("returns and okResult with a 405 status when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/lists", nil)
		request = addUserIDToContext(jwtInfo.UserID, request)

		ListsHandler(request, testSrvProvider)

		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})
}

func listDtoToCreate() models.ListDto {
	return models.ListDto{
		Name: "new list",
		Items: []models.Item{
			models.Item{
				Title:       "title",
				Description: "desc",
			},
		},
	}
}

func listDtoToUpdate() models.ListDto {
	return models.ListDto{
		Name: "updated list",
		Items: []models.Item{
			models.Item{
				Title:       "replaced title",
				Description: "replaced desc",
			},
		},
	}
}

func assertListsExpectations(t *testing.T, sp *mockedServiceProvider, ls *mockedListsService) {
	t.Helper()

	sp.AssertExpectations(t)
	ls.AssertExpectations(t)
}

func addUserIDToContext(userID string, r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, reqContextUserKey, userID)

	return r.WithContext(ctx)
}
