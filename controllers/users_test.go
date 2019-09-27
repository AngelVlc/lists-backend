package controllers

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"strings"
	"errors"
	"github.com/AngelVlc/lists-backend/models"
)


func TestUsers(t *testing.T) {
	testObj := new(mockedStore)

	handler := Handler {
		HandlerFunc: UsersHandler,
		Repository: testObj,
	}

	t.Run("POST adds a new user and returns it", func(t *testing.T) {
		userDto := userDtoToCreate()

		data := userDto.ToUser()

		testObj.On("AddUser", &data).Return(nil).Once()

		body, _ := json.Marshal(userDto)
		request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusCreated)
	})

	t.Run("POST with invalid body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST without body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST returns 500 when the insert fails", func(t *testing.T) {
		userDto := userDtoToCreate()

		data := userDto.ToUser()

		testObj.On("AddUser", &data).Return(errors.New("wadus")).Once()

		body, _ := json.Marshal(userDto)
		request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("returns 405 when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/users", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertResult(t, testObj, response.Result().StatusCode, http.StatusMethodNotAllowed)
	})
}

func userDtoToCreate() models.UserDto {
	return models.UserDto{
		UserName:           "newUser1",
		NewPassword:        "password",
		ConfirmNewPassword: "password",
		IsAdmin:            true,
	}
}
