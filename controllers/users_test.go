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

	t.Run("POST returns an okResult when there is no error", func(t *testing.T) {
		userDto := userDtoToCreate()

		data := userDto.ToUser()

		testSrvProvider.On("GetUsersService").Return(testUsersSrv).Once()

		testUsersSrv.On("AddUser", &data).Return(nil).Once()

		body, _ := json.Marshal(userDto)
		request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		got := UsersHandler(response, request, testSrvProvider)
		want := okResult{data, http.StatusCreated}

		assert.Equal(t, want, got, "should be equal")
		assertUsersExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST with invalid body should return an errorResult with an InvalidBodyError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		got := UsersHandler(response, request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		_, isInvalidBodyError := errorRes.err.(*InvalidBodyError)
		assert.Equal(t, true, isInvalidBodyError, "should be an invalid body error")
		assertUsersExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST without body should return an errorResult with a NoBodyError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", nil)
		response := httptest.NewRecorder()

		got := UsersHandler(response, request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		_, isNoBodyError := errorRes.err.(*NoBodyError)
		assert.Equal(t, true, isNoBodyError, "should be an invalid body error")
		assertUsersExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST returns an errorResult with the service error when the insert fails", func(t *testing.T) {
		userDto := userDtoToCreate()

		data := userDto.ToUser()

		testSrvProvider.On("GetUsersService").Return(testUsersSrv).Once()

		err := errors.New("wadus")
		testUsersSrv.On("AddUser", &data).Return(err).Once()

		body, _ := json.Marshal(userDto)
		request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		got := UsersHandler(response, request, testSrvProvider)

		errorResult, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assert.Equal(t, err, errorResult.err)
		assertUsersExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("returns and okResult with a 405 status when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/users", nil)
		response := httptest.NewRecorder()

		got := UsersHandler(response, request, testSrvProvider)

		want := okResult{nil, http.StatusMethodNotAllowed}

		assert.Equal(t, want, got, "should be equal")
		assertUsersExpectations(t, testSrvProvider, testUsersSrv)
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

func assertUsersExpectations(t *testing.T, sp *mockedServiceProvider, us *mockedUsersService) {
	t.Helper()

	sp.AssertExpectations(t)
	us.AssertExpectations(t)
}
