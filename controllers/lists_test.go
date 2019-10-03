package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockedListsService struct {
	mock.Mock
}

func (us *mockedListsService) AddList(l *models.List) error {
	args := us.Called(l)
	return args.Error(0)
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

	handler := Handler{
		HandlerFunc:     ListsHandler,
		ServiceProvider: testSrvProvider,
	}

	t.Run("GET returns list items", func(t *testing.T) {
		data := models.SampleGetListsResultDto()

		testSrvProvider.On("GetListsService").Return(testListsSrv)
		testListsSrv.On("GetLists", &[]models.GetListsResultDto{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]models.GetListsResultDto)
			*arg = data
		})

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		bytes, _ := json.Marshal(data)

		want := string(bytes) + "\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusOK)
	})

	t.Run("GET returns 500 when when an unhandled error occurs", func(t *testing.T) {
		testSrvProvider.On("GetListsService").Return(testListsSrv)
		testListsSrv.On("GetLists", &[]models.GetListsResultDto{}).Return(errors.New("wadus")).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("GET WITH AN ID returns 404 when the id is not valid", func(t *testing.T) {
		id := "wadus"
		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		testSrvProvider.On("GetListsService").Return(testListsSrv)
		testListsSrv.On("GetSingleList", id, &models.List{}).Return(&stores.InvalidIDError{ID: id}).Once()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("GET WITH AN ID returns 500 when an unhandled error occurs", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testSrvProvider.On("GetListsService").Return(testListsSrv)
		testListsSrv.On("GetSingleList", id, &models.List{}).Return(errors.New("wadus")).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("GET WITH AN ID returns 500 when a query error occurs", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testListsSrv.On("GetSingleList", id, &models.List{}).Return(&stores.UnexpectedError{InternalError: errors.New("wadus"), Msg: "store error"}).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("GET WITH AN ID returns 400 when the list does not exist", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testSrvProvider.On("GetListsService").Return(testListsSrv)
		testListsSrv.On("GetSingleList", id, &models.List{}).Return(&stores.NotFoundError{ID: id, Model: "lists"}).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusNotFound)
	})

	t.Run("GET WITH AN ID returns a single list", func(t *testing.T) {
		data := models.SampleListSlice()[0]

		testSrvProvider.On("GetListsService").Return(testListsSrv)
		testListsSrv.On("GetSingleList", data.ID, &models.List{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*models.List)
			*arg = data
		})

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+data.ID, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusOK)

		bytes, _ := json.Marshal(data)

		want := string(bytes) + "\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")
	})

	t.Run("POST adds a new list and returns it", func(t *testing.T) {
		listDto := listDtoToCreate()

		data := listDto.ToList()

		testSrvProvider.On("GetListsService").Return(testListsSrv)

		testListsSrv.On("AddList", &data).Return(nil).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusCreated)
	})

	t.Run("POST with invalid body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST without body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST returns 500 when the insert fails", func(t *testing.T) {
		listDto := listDtoToCreate()

		data := listDto.ToList()

		testSrvProvider.On("GetListsService").Return(testListsSrv)

		testListsSrv.On("AddList", &data).Return(errors.New("wadus")).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("DELETE returns 500 when the remove fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testListsSrv.On("RemoveList", id).Return(errors.New("wadus")).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("DELETE returns 400 when the id is not valid", func(t *testing.T) {
		id := "wadus"

		testListsSrv.On("RemoveList", id).Return(&stores.InvalidIDError{ID: id}).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("DELETE returns 404 when the list does not exist", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testListsSrv.On("RemoveList", id).Return(&stores.NotFoundError{Model: "lists", ID: id}).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusNotFound)
	})

	t.Run("DELETE removes a list", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testListsSrv.On("RemoveList", id).Return(nil).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusNoContent)
	})

	t.Run("PUT with invalid body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("PUT without body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("PUT returns 500 when the update fails", func(t *testing.T) {
		listDto := listDtoToUpdate()
		id := bson.NewObjectId().Hex()

		testListsSrv.On("UpdateList", id, mock.Anything).Return(errors.New("wadus")).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("PUT returns 404 when the document does not exist", func(t *testing.T) {
		listDto := listDtoToUpdate()
		id := bson.NewObjectId().Hex()

		testListsSrv.On("UpdateList", id, mock.Anything).Return(&stores.NotFoundError{Model: "Lists", ID: id}).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusNotFound)
	})

	t.Run("PUT updates a new list and returns it", func(t *testing.T) {
		listDto := listDtoToUpdate()

		data := listDto.ToList()

		id := bson.NewObjectId().Hex()

		testListsSrv.On("UpdateList", id, &data).Return(nil).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusOK)
	})

	t.Run("returns 405 when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertListsResult(t, testSrvProvider, testListsSrv, response.Result().StatusCode, http.StatusMethodNotAllowed)
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

func assertListsResult(t *testing.T, sp *mockedServiceProvider, ls *mockedListsService, got, want int) {
	t.Helper()

	assert.Equal(t, want, got, "status are not equal")

	sp.AssertExpectations(t)
	ls.AssertExpectations(t)
}
