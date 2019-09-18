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

func TestLists(t *testing.T) {
	testObj := new(MockedStore)

	server := newServer(testObj)

	t.Run("GET returns list items", func(t *testing.T) {
		data := models.SampleListCollectionSlice()

		testObj.On("GetLists").Return(data)

		log.Println(data)

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testObj.AssertExpectations(t)

		bytes, _ := json.Marshal(data)

		want := string(bytes) + "\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")

		assertStatus(t, response.Result().StatusCode, http.StatusOK)
	})

	t.Run("POST adds a new list and returns it", func(t *testing.T) {
		listDto := models.ListDto{
			Name: "new list",
			Items: []models.Item{
				models.Item{
					Title:       "title",
					Description: "desc",
				},
			},
		}

		data := models.List{
			Name:  listDto.Name,
			Items: listDto.Items,
		}

		testObj.On("AddList", &data).Return(nil)

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusCreated)
	})

	t.Run("POST with invalid body should return 404", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/lists", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST returns 400 when the insert fails", func(t *testing.T) {
		listDto := models.ListDto{
			Name: "new list",
			Items: []models.Item{
				models.Item{
					Title:       "title",
					Description: "desc",
				},
			},
		}

		testObj.ExpectedCalls = []*mock.Call{}
		testObj.On("AddList", mock.Anything).Return(errors.New("wadus"))

		body, _ := json.Marshal(listDto)
		request, _ := http.NewRequest(http.MethodPost, "/lists", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("DELETE returns 400 when the remove fails", func(t *testing.T) {
		testObj.On("RemoveList", "1").Return(errors.New("wadus"))

		request, _ := http.NewRequest(http.MethodDelete, "/lists/1", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("DELETE removes a list", func(t *testing.T) {
		testObj.On("RemoveList", "2").Return(nil)

		request, _ := http.NewRequest(http.MethodDelete, "/lists/2", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Result().StatusCode, http.StatusNoContent)
	})
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
