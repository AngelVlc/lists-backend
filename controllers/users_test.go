package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	appErrors "github.com/AngelVlc/lists-backend/errors"
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

func (us *mockedUsersService) AddUser(dto *models.UserDto) error {
	args := us.Called(dto)
	return args.Error(0)
}

func TestUsersHandler(t *testing.T) {
	testUsersSrv := new(mockedUsersService)

	testSrvProvider := new(mockedServiceProvider)

	t.Run("POST returns an okResult when there is no error", func(t *testing.T) {
		userDto := userDtoToCreate()

		user := userDto.ToUser()

		testSrvProvider.On("GetUsersService").Return(testUsersSrv).Once()

		testUsersSrv.On("AddUser", &userDto).Return(nil).Once()

		body, _ := json.Marshal(userDto)
		request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		request.Header.Set("Content-type", "application/json")
		response := httptest.NewRecorder()

		got := UsersHandler(response, request, testSrvProvider)
		want := okResult{user, http.StatusCreated}

		assert.Equal(t, want, got, "should be equal")
		assertUsersExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST with invalid body should return an errorResult with a BadRequestError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", strings.NewReader("wadus"))
		response := httptest.NewRecorder()

		got := UsersHandler(response, request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		badReqErr, isInvalidBodyError := errorRes.err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isInvalidBodyError, "should be a bad request error")
		assert.Equal(t, "Invalid body", badReqErr.Error())

		assertUsersExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST without body should return an errorResult with a BadRequestError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/users", nil)
		response := httptest.NewRecorder()

		got := UsersHandler(response, request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		badReqErr, isNoBodyError := errorRes.err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isNoBodyError, "should be a bad request error")
		assert.Equal(t, "No body", badReqErr.Error())

		assertUsersExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST returns an errorResult with the service error when the insert fails", func(t *testing.T) {
		userDto := userDtoToCreate()

		testSrvProvider.On("GetUsersService").Return(testUsersSrv).Once()

		err := errors.New("wadus")
		testUsersSrv.On("AddUser", &userDto).Return(err).Once()

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
