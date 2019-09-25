package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

type MockedStore struct {
	mock.Mock
}

func (m *MockedStore) GetLists() ([]models.GetListsResultDto, error) {
	args := m.Called()
	return args.Get(0).([]models.GetListsResultDto), args.Error(1)
}

func (m *MockedStore) AddList(l *models.List) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *MockedStore) RemoveList(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockedStore) UpdateList(id string, l *models.List) error {
	args := m.Called(id, l)
	return args.Error(0)
}

func (m *MockedStore) GetSingleList(id string) (models.List, error) {
	args := m.Called(id)
	return args.Get(0).(models.List), args.Error(1)
}

func TestLists(t *testing.T) {
	testObj := new(MockedStore)

	server := newServer(testObj)

	t.Run("GET returns list items", func(t *testing.T) {
		data := models.SampleGetListsResultDto()

		testObj.On("GetLists").Return(data, nil).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		bytes, _ := json.Marshal(data)

		want := string(bytes) + "\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")

		assertResult(t, testObj, response.Result().StatusCode, http.StatusOK)
	})

	t.Run("GET returns 500 when when an unhandled error occurs", func(t *testing.T) {
		testObj.On("GetLists").Return([]models.GetListsResultDto{}, errors.New("wadus")).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("GET WITH AN ID returns 404 when the id is not valid", func(t *testing.T) {
		id := "wadus"
		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		testObj.On("GetSingleList", id).Return(models.List{}, &stores.InvalidIDError{ID: id}).Once()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("GET WITH AN ID returns 500 when an unhandled error occurs", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testObj.On("GetSingleList", id).Return(models.List{}, errors.New("wadus")).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("GET WITH AN ID returns 500 when a query error occurs", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testObj.On("GetSingleList", id).Return(models.List{}, &stores.UnexpectedError{InternalError: errors.New("wadus"), Msg: "store error"}).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("GET WITH AN ID returns 400 when the list does not exist", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testObj.On("GetSingleList", id).Return(models.List{}, &stores.NotFoundError{ID: id, Model: "lists"}).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusNotFound)
	})

	t.Run("GET WITH AN ID returns a single list", func(t *testing.T) {
		data := models.SampleListSlice()[0]

		testObj.On("GetSingleList", data.ID).Return(data, nil).Once()

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+data.ID, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusOK)

		bytes, _ := json.Marshal(data)

		want := string(bytes) + "\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")
	})

	t.Run("POST adds a new list and returns it", func(t *testing.T) {
		listDto := listDtoToCreate()

		data := listFromDto(&listDto)

		testObj.On("AddList", &data).Return(nil).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusCreated)
	})

	t.Run("POST with invalid body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST without body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST returns 500 when the insert fails", func(t *testing.T) {
		listDto := listDtoToCreate()
		listDto.Name += " with error"

		data := listFromDto(&listDto)

		testObj.On("AddList", &data).Return(errors.New("wadus")).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("DELETE returns 500 when the remove fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testObj.On("RemoveList", id).Return(errors.New("wadus")).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("DELETE returns 400 when the id is not valid", func(t *testing.T) {
		id := "wadus"

		testObj.On("RemoveList", id).Return(&stores.InvalidIDError{ID: id}).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("DELETE returns 404 when the list does not exist", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testObj.On("RemoveList", id).Return(&stores.NotFoundError{Model: "lists", ID: id}).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusNotFound)
	})

	t.Run("DELETE removes a list", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testObj.On("RemoveList", id).Return(nil).Once()

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusNoContent)
	})

	t.Run("PUT with invalid body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("PUT without body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("PUT returns 500 when the update fails", func(t *testing.T) {
		listDto := listDtoToUpdate()
		id := bson.NewObjectId().Hex()

		testObj.On("UpdateList", id, mock.Anything).Return(errors.New("wadus")).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("PUT returns 404 when the document does not exist", func(t *testing.T) {
		listDto := listDtoToUpdate()
		id := bson.NewObjectId().Hex()

		testObj.On("UpdateList", id, mock.Anything).Return(&stores.NotFoundError{Model: "Lists", ID: id}).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusNotFound)
	})

	t.Run("PUT updates a new list and returns it", func(t *testing.T) {
		listDto := listDtoToUpdate()

		data := listFromDto(&listDto)

		id := bson.NewObjectId().Hex()

		testObj.On("UpdateList", id, &data).Return(nil).Once()

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusOK)
	})

	t.Run("returns 405 when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusMethodNotAllowed)
	})
}

func assertResult(t *testing.T, m *MockedStore, got, want int) {
	t.Helper()

	assert.Equal(t, want, got, "status are not equal")

	m.AssertExpectations(t)
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

func listFromDto(listDto *models.ListDto) models.List {
	return models.List{
		Name:  listDto.Name,
		Items: listDto.Items,
	}
}
