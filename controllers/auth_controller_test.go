package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/AngelVlc/lists-backend/models"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/stretchr/testify/assert"
)

func TestAuthTokenHandler(t *testing.T) {
	testUsersSrv := new(mockedUsersService)
	testAuthSrv := new(mockedAuthService)

	testSrvProvider := new(mockedServiceProvider)

	t.Run("POST with invalid body should return an errorResult with a BadRequestError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/auth/token", strings.NewReader("wadus"))

		got := AuthHandler(request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		badReqErr, isInvalidBodyError := errorRes.err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isInvalidBodyError, "should be a bad request error")
		assert.Equal(t, "Invalid body", badReqErr.Error())

		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST without user name in body should return an errorResult with a BadRequestError", func(t *testing.T) {
		login := struct {
			Password string
		}{
			"pass",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		got := AuthHandler(request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		badReqErr, isInvalidBodyError := errorRes.err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isInvalidBodyError, "should be a bad request error")
		assert.Equal(t, "UserName is mandatory", badReqErr.Error())

		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST without pasword in body should return an errorResult with a BadRequestError", func(t *testing.T) {
		login := struct {
			UserName string
		}{
			"wadus",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		got := AuthHandler(request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		badReqErr, isInvalidBodyError := errorRes.err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isInvalidBodyError, "should be a bad request error")
		assert.Equal(t, "Password is mandatory", badReqErr.Error())

		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST returns and okResult with a 405 status when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/auth/token", nil)

		got := AuthHandler(request, testSrvProvider)

		want := okResult{nil, http.StatusMethodNotAllowed}

		assert.Equal(t, want, got, "should be equal")
		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST returns an errorResult when the CheckIfUserPasswordIsOk() returns an error", func(t *testing.T) {
		login := models.Login{
			UserName: "wadus",
			Password: "pass",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		testSrvProvider.On("GetUsersService").Return(testUsersSrv).Once()

		testUsersSrv.On("CheckIfUserPasswordIsOk", login.UserName, login.Password).Return(nil, errors.New("wadus")).Once()

		got := AuthHandler(request, testSrvProvider)

		_, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST returns an errorResult when CreateToken returns an error", func(t *testing.T) {
		login := models.Login{
			UserName: "wadus",
			Password: "pass",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		testSrvProvider.On("GetUsersService").Return(testUsersSrv).Once()

		user := models.User{
			UserName: login.UserName,
			ID:       "id",
		}
		testUsersSrv.On("CheckIfUserPasswordIsOk", login.UserName, login.Password).Return(&user, nil).Once()

		testSrvProvider.On("GetAuthService").Return(testAuthSrv).Once()

		testAuthSrv.On("CreateToken", &user).Return("", errors.New("wadus")).Once()

		got := AuthHandler(request, testSrvProvider)

		_, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("POST returns an okResult when there is no error", func(t *testing.T) {
		login := models.Login{
			UserName: "wadus",
			Password: "pass",
		}
		body, _ := json.Marshal(login)

		request, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBuffer(body))

		testSrvProvider.On("GetUsersService").Return(testUsersSrv).Once()

		user := models.User{
			UserName: login.UserName,
			ID:       "id",
		}
		testUsersSrv.On("CheckIfUserPasswordIsOk", login.UserName, login.Password).Return(&user, nil).Once()

		testSrvProvider.On("GetAuthService").Return(testAuthSrv).Once()

		token := "theToken"
		testAuthSrv.On("CreateToken", &user).Return(token, nil).Once()

		got := AuthHandler(request, testSrvProvider)

		want := okResult{token, http.StatusOK}

		assert.Equal(t, want, got, "should be equal")
		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})
}

func assertAuthExpectations(t *testing.T, sp *mockedServiceProvider, us *mockedUsersService) {
	t.Helper()

	sp.AssertExpectations(t)
	us.AssertExpectations(t)
}
