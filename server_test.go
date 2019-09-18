package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockedStore) GetLists() []models.List {
	args := m.Called()
	return args.Get(0).([]models.List)
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

func TestLists(t *testing.T) {
	testObj := new(MockedStore)

	server := newServer(testObj)

	t.Run("GET returns list items", func(t *testing.T) {
		data := models.SampleListCollectionSlice()

		testObj.On("GetLists").Return(data)

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testObj.AssertExpectations(t)

		bytes, _ := json.Marshal(data)

		want := string(bytes) + "\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")

		assertStatus(t, response.Result().StatusCode, http.StatusOK)

		testObj.AssertExpectations(t)
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

		assertStatus(t, response.Result().StatusCode, http.StatusCreated)

		testObj.AssertExpectations(t)
	})

	t.Run("POST with invalid body should return 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusBadRequest)

		testObj.AssertExpectations(t)
	})

	t.Run("POST returns 400 when the insert fails", func(t *testing.T) {
		listDto := listDtoToCreate()

		testObj.ExpectedCalls = []*mock.Call{}
		testObj.On("AddList", mock.Anything).Return(errors.New("wadus"))

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusBadRequest)

		testObj.AssertExpectations(t)
	})

	t.Run("DELETE returns 400 when the remove fails", func(t *testing.T) {
		testObj.On("RemoveList", "1").Return(errors.New("wadus"))

		request, _ := http.NewRequest(http.MethodDelete, "/lists/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusBadRequest)

		testObj.AssertExpectations(t)
	})

	t.Run("DELETE removes a list", func(t *testing.T) {
		testObj.On("RemoveList", "2").Return(nil)

		request, _ := http.NewRequest(http.MethodDelete, "/lists/2", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusNoContent)

		testObj.AssertExpectations(t)
	})

	t.Run("PUT with invalid body should return 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPut, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusBadRequest)

		testObj.AssertExpectations(t)
	})

	t.Run("PUT returns 400 when the update fails", func(t *testing.T) {
		listDto := listDtoToUpdate()

		testObj.ExpectedCalls = []*mock.Call{}
		testObj.On("UpdateList", "1", mock.Anything).Return(errors.New("wadus"))

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/1", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusBadRequest)

		testObj.AssertExpectations(t)
	})

	t.Run("PUT updates a new list and returns it", func(t *testing.T) {
		listDto := listDtoToUpdate()

		data := listFromDto(&listDto)

		testObj.On("UpdateList", "2", &data).Return(nil)

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPut, "/lists/2", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusOK)

		testObj.AssertExpectations(t)
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
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
