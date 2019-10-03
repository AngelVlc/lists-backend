package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockedUsersService struct {
	mock.Mock
}

func (us *mockedUsersService) AddUser(u *models.User) error {
	args := us.Called(u)
	return args.Error(0)
}

func TestUsersHandler(t *testing.T) {
	testUsersSrv := new(mockedUsersService)

	testSrvProvider := new(mockedServiceProvider)

	handler := Handler{
		HandlerFunc:     UsersHandler,
		ServiceProvider: testSrvProvider,
	}

	t.Run("POST adds a new user and returns it", func(t *testing.T) {
		userDto := userDtoToCreate()

		data := userDto.ToUser()

		testSrvProvider.On("GetUsersService").Return(testUsersSrv)

		testUsersSrv.On("AddUser", &data).Return(nil).Once()

		body, _ := json.Marshal(userDto)
		request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertUsersResult(t, testSrvProvider, testUsersSrv, response.Result().StatusCode, http.StatusCreated)
	})

	t.Run("POST with invalid body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertUsersResult(t, testSrvProvider, testUsersSrv, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST without body should return 400", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertUsersResult(t, testSrvProvider, testUsersSrv, response.Result().StatusCode, http.StatusBadRequest)
	})

	t.Run("POST returns 500 when the insert fails", func(t *testing.T) {
		userDto := userDtoToCreate()

		data := userDto.ToUser()

		testSrvProvider.On("GetUsersService").Return(testUsersSrv)

		testUsersSrv.On("AddUser", &data).Return(errors.New("wadus")).Once()

		body, _ := json.Marshal(userDto)
		request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertUsersResult(t, testSrvProvider, testUsersSrv, response.Result().StatusCode, http.StatusInternalServerError)
	})

	t.Run("returns 405 when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/users", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assertUsersResult(t, testSrvProvider, testUsersSrv, response.Result().StatusCode, http.StatusMethodNotAllowed)
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

func assertUsersResult(t *testing.T, sp *mockedServiceProvider, us *mockedUsersService, got, want int) {
	t.Helper()

	assert.Equal(t, want, got, "status are not equal")

	sp.AssertExpectations(t)
	us.AssertExpectations(t)
}
