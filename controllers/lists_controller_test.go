package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strings"
	"testing"
)

type mockedListsService struct {
	mock.Mock
}

func (us *mockedListsService) AddList(l *models.List) (string, error) {
	args := us.Called(l)
	return args.String(0), args.Error(1)
}

func (us *mockedListsService) RemoveList(id string) error {
	args := us.Called(id)
	return args.Error(0)
}

func (us *mockedListsService) UpdateList(id string, l *models.List) error {
	args := us.Called(id, l)
	return args.Error(0)
}

func (us *mockedListsService) GetSingleList(id string, l *models.List) error {
	args := us.Called(id, l)
	return args.Error(0)
}

func (us *mockedListsService) GetLists(r *[]models.GetListsResultDto) error {
	args := us.Called(r)
	return args.Error(0)
}

func TestLists(t *testing.T) {
	testListsSrv := new(mockedListsService)

	testSrvProvider := new(mockedServiceProvider)

	t.Run("GET returns an okResult when there is no error", func(t *testing.T) {
		data := models.SampleGetListsResultDto()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		testListsSrv.On("GetLists", &[]models.GetListsResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]models.GetListsResultDto)
			*arg = data
		})

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)

		got := ListsHandler(request, testSrvProvider, nil)

		want := okResult{data, http.StatusOK}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("GET returns an errorResult with the service error when the query fails", func(t *testing.T) {
		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		err := errors.New("wadus")
		testListsSrv.On("GetLists", &[]models.GetListsResultDto{}).Return(err).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)

		got := ListsHandler(request, testSrvProvider, nil)

		errorResult, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorResult.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("GET WITH AN ID returns an errorResult with the service error when the query fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		err := errors.New("wadus")
		testListsSrv.On("GetSingleList", id, &models.List{}).Return(errors.New("wadus")).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)

		got := ListsHandler(request, testSrvProvider, nil)

		errorResult, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorResult.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("GET WITH AN ID returns an okResult with a single list", func(t *testing.T) {
		data := models.SampleListSlice()[0]

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		testListsSrv.On("GetSingleList", data.ID, &models.List{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*models.List)
			*arg = data
		})

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+data.ID, nil)

		got := ListsHandler(request, testSrvProvider, nil)

		want := okResult{data, http.StatusOK}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("POST returns an okResult when there is no error", func(t *testing.T) {
		listDto := listDtoToCreate()

		data := listDto.ToList()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()

		testListsSrv.On("AddList", &data).Return("id", nil).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")

		got := ListsHandler(request, testSrvProvider, nil)
		want := okResult{"id", http.StatusCreated}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("POST with invalid body should return an errorResult with a BadRequestError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", strings.NewReader("wadus"))

		got := ListsHandler(request, testSrvProvider, nil)

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
		testListsSrv.On("AddList", &data).Return("", err).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")

		got := ListsHandler(request, testSrvProvider, nil)
		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorRes.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("DELETE returns an errorResult with the service erorr when the delete fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		err := errors.New("wadus")
		testListsSrv.On("RemoveList", id).Return(err).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)

		got := ListsHandler(request, testSrvProvider, nil)
		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorRes.err)
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("DELETE returns an okResult when there is no error", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testSrvProvider.On("GetListsService").Return(testListsSrv).Once()
		testListsSrv.On("RemoveList", id).Return(nil).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)

		got := ListsHandler(request, testSrvProvider, nil)
		want := okResult{nil, http.StatusNoContent}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("PUT with invalid body should return an errorResult with a BadRequestError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/lists", strings.NewReader("wadus"))

		got := ListsHandler(request, testSrvProvider, nil)

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
		testListsSrv.On("UpdateList", id, mock.Anything).Return(err).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")

		got := ListsHandler(request, testSrvProvider, nil)
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
		testListsSrv.On("UpdateList", id, &data).Return(nil).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")

		got := ListsHandler(request, testSrvProvider, nil)
		want := okResult{data, http.StatusOK}

		assert.Equal(t, want, got, "should be equal")
		assertListsExpectations(t, testSrvProvider, testListsSrv)
	})

	t.Run("returns and okResult with a 405 status when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/lists", nil)

		ListsHandler(request, testSrvProvider, nil)

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
