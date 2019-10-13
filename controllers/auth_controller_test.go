package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler(t *testing.T) {
	testUsersSrv := new(mockedUsersService)

	testSrvProvider := new(mockedServiceProvider)

	// t.Run("POST returns an okResult when there is no error", func(t *testing.T) {
	// 	userDto := userDtoToCreate()

	// 	testSrvProvider.On("GetUsersService").Return(testUsersSrv).Once()

	// 	testUsersSrv.On("AddUser", &userDto).Return("id", nil).Once()

	// 	body, _ := json.Marshal(userDto)
	// 	request, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	// 	request.Header.Set("Content-type", "application/json")

	// 	got := AuthHandler(request, testSrvProvider)
	// 	want := okResult{"id", http.StatusCreated}

	// 	assert.Equal(t, want, got, "should be equal")
	// 	assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	// })

	t.Run("POST with invalid body should return an errorResult with a BadRequestError", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/auth", strings.NewReader("wadus"))

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

		request, _ := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))

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

		request, _ := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))

		got := AuthHandler(request, testSrvProvider)

		errorRes, isErrorResult := got.(errorResult)
		assert.Equal(t, true, isErrorResult, "should be an error result")

		badReqErr, isInvalidBodyError := errorRes.err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isInvalidBodyError, "should be a bad request error")
		assert.Equal(t, "Password is mandatory", badReqErr.Error())

		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})

	t.Run("returns and okResult with a 405 status when the method is not GET, POST, PUT or DELETE", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/auth", nil)

		got := AuthHandler(request, testSrvProvider)

		want := okResult{nil, http.StatusMethodNotAllowed}

		assert.Equal(t, want, got, "should be equal")
		assertAuthExpectations(t, testSrvProvider, testUsersSrv)
	})
}

func assertAuthExpectations(t *testing.T, sp *mockedServiceProvider, us *mockedUsersService) {
	t.Helper()

	sp.AssertExpectations(t)
	us.AssertExpectations(t)
}
