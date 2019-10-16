package controllers

import (
	"errors"
	"fmt"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockedAuthService struct {
	mock.Mock
}

func (s *mockedAuthService) CreateTokens(u *models.User) (map[string]string, error) {
	args := s.Called(u)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (s *mockedAuthService) ParseToken(token string) (*models.JwtClaimsInfo, error) {
	args := s.Called(token)
	return args.Get(0).(*models.JwtClaimsInfo), args.Error(1)
}

func TestHandlerWithoutAuth(t *testing.T) {
	mockServicePrv := new(mockedServiceProvider)

	t.Run("Returns 200 when no error", func(t *testing.T) {
		f := func(r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return okResult{nil, http.StatusOK}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("Returns 200 with content when no error", func(t *testing.T) {
		f := func(r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			obj := struct {
				Field1 string
				Field2 string
			}{Field1: "a", Field2: "b"}
			return okResult{obj, http.StatusOK}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		want := "{\"Field1\":\"a\",\"Field2\":\"b\"}\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("Returns 500 when an unexpected error happens", func(t *testing.T) {
		f := func(r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{&appErrors.UnexpectedError{Msg: "error", InternalError: errors.New("msg")}}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
		assert.Equal(t, "error\n", string(response.Body.String()))
	})

	t.Run("Returns 404 when a not found error happens", func(t *testing.T) {
		f := func(r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{&appErrors.NotFoundError{ID: "id", Model: "model"}}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
		assert.Equal(t, "model with id \"id\" not found\n", string(response.Body.String()))
	})

	t.Run("Returns 400 when a bad request error happens", func(t *testing.T) {
		f := func(r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{&appErrors.BadRequestError{Msg: fmt.Sprintf("%q is not a valid id", "id")}}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
		assert.Equal(t, "\"id\" is not a valid id\n", string(response.Body.String()))
	})

	t.Run("Returns 401 when an unauthorized error happens", func(t *testing.T) {
		f := func(r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{&appErrors.UnauthorizedError{Msg: "wadus"}}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "wadus\n", string(response.Body.String()))
	})

	t.Run("Returns 500 when an unhandled error happens", func(t *testing.T) {
		f := func(r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{errors.New("wadus")}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
		assert.Equal(t, "Internal error\n", string(response.Body.String()))
	})
}

func TestHandlerWithAuth(t *testing.T) {
	mockServicePrv := new(mockedServiceProvider)

	mockAuthSvc := new(mockedAuthService)

	f := func(r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
		return errorResult{errors.New("wadus")}
	}

	t.Run("Returns 401 when the request does not have auth header", func(t *testing.T) {
		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
			RequireAuth:     true,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "No authorization header\n", string(response.Body.String()))
	})

	t.Run("Returns 401 when the request auth header is not valid", func(t *testing.T) {
		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
			RequireAuth:     true,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "bad_header")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "Invalid authorization header\n", string(response.Body.String()))
	})

	t.Run("Returns 401 when the auth token is not valid", func(t *testing.T) {
		mockServicePrv.On("GetAuthService").Return(mockAuthSvc).Once()

		mockAuthSvc.On("ParseToken", "token").Return(&models.JwtClaimsInfo{}, errors.New("wadus")).Once()

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
			RequireAuth:     true,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "Bearer token")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode)
		assert.Equal(t, "Invalid auth token\n", string(response.Body.String()))
	})

	t.Run("Returns 403 when the resource requires admin and the user is not admin", func(t *testing.T) {
		mockServicePrv.On("GetAuthService").Return(mockAuthSvc).Once()

		mockAuthSvc.On("ParseToken", "token").Return(&models.JwtClaimsInfo{IsAdmin: false}, nil).Once()

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: mockServicePrv,
			RequireAuth:     true,
			RequireAdmin:    true,
		}

		request, _ := http.NewRequest(http.MethodGet, "/wadus", nil)
		request.Header.Set("Authorization", "Bearer token")
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusForbidden, response.Result().StatusCode)
		assert.Equal(t, "Access forbidden\n", string(response.Body.String()))
	})

}
