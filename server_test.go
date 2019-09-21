package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/AngelVlc/lists-backend/models"
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

func (m *MockedStore) GetLists() []models.GetListsResultDto {
	args := m.Called()
	return args.Get(0).([]models.GetListsResultDto)
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

		testObj.On("GetLists").Return(data)

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		bytes, _ := json.Marshal(data)

		want := string(bytes) + "\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusOK)
	})

	t.Run("GET WITH AN ID returns 400 when the id is not valid", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/lists/wadus", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("GET WITH AN ID returns 400 when the query fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testObj.On("GetSingleList", id).Return(models.List{}, errors.New("wadus"))

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("GET WITH AN ID returns a single list", func(t *testing.T) {
		data := models.SampleListSlice()[0]

		testObj.On("GetSingleList", data.ID.Hex()).Return(data, nil)

		request, _ := http.NewRequest(http.MethodGet, "/lists/"+data.ID.Hex(), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusOK)

		bytes, _ := json.Marshal(data)

		want := string(bytes) + "\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")
	})

	t.Run("POST adds a new list and returns it", func(t *testing.T) {
		listDto := listDtoToCreate()

		data := listFromDto(&listDto)

		testObj.On("AddList", &data).Return(nil)

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusCreated)
	})

	t.Run("POST with invalid body should return 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST returns 400 when the insert fails", func(t *testing.T) {
		listDto := listDtoToCreate()
		listDto.Name += " with error"

		data := listFromDto(&listDto)

		testObj.On("AddList", &data).Return(errors.New("wadus"))

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("DELETE returns 400 when the remove fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testObj.On("RemoveList", id).Return(errors.New("wadus"))

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("DELETE removes a list", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testObj.On("RemoveList", id).Return(nil)

		request, _ := http.NewRequest(http.MethodDelete, "/lists/"+id, nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusNoContent)
	})

	t.Run("PUT with invalid body should return 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("PUT returns 400 when the update fails", func(t *testing.T) {
		listDto := listDtoToUpdate()
		id := bson.NewObjectId().Hex()

		testObj.On("UpdateList", id, mock.Anything).Return(errors.New("wadus"))

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("PUT updates a new list and returns it", func(t *testing.T) {
		listDto := listDtoToUpdate()

		data := listFromDto(&listDto)

		id := bson.NewObjectId().Hex()

		testObj.On("UpdateList", id, &data).Return(nil)

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/"+id, bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusOK)
	})

	t.Run("returns 405 when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, testObj, response.Result().StatusCode, http.StatusMethodNotAllowed)
	})
}

func assertStatus(t *testing.T, m *MockedStore, got, want int) {
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
