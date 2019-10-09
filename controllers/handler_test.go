package controllers

import (
	"errors"
	"fmt"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	testSrvProvider := new(mockedServiceProvider)

	t.Run("Returns 200 when no error", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return okResult{nil, http.StatusOK}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: testSrvProvider,
		}

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("Return 200 with content when no error", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			obj := struct {
				Field1 string
				Field2 string
			}{Field1: "a", Field2: "b"}
			return okResult{obj, http.StatusOK}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: testSrvProvider,
		}

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		want := "{\"Field1\":\"a\",\"Field2\":\"b\"}\n"

		got := string(response.Body.String())

		assert.Equal(t, want, got, "they should be equal")

		assert.Equal(t, http.StatusOK, response.Result().StatusCode)
	})

	t.Run("Returns 500 when an unexpected error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{&appErrors.UnexpectedError{Msg: "error", InternalError: errors.New("msg")}}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: testSrvProvider,
		}

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
		assert.Equal(t, "error\n", string(response.Body.String()))
	})

	t.Run("Returns 404 when a not found error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{&appErrors.NotFoundError{ID: "id", Model: "model"}}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: testSrvProvider,
		}

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusNotFound, response.Result().StatusCode)
		assert.Equal(t, "model with id \"id\" not found\n", string(response.Body.String()))
	})

	t.Run("Returns 400 when a bad request error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{&appErrors.BadRequestError{Msg: fmt.Sprintf("%q is not a valid id", "id")}}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: testSrvProvider,
		}

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Result().StatusCode)
		assert.Equal(t, "\"id\" is not a valid id\n", string(response.Body.String()))
	})

	t.Run("Returns 500 when an unhandled error happens", func(t *testing.T) {
		f := func(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
			return errorResult{errors.New("wadus")}
		}

		handler := Handler{
			HandlerFunc:     f,
			ServiceProvider: testSrvProvider,
		}

		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		handler.ServeHTTP(response, request)

		assert.Equal(t, http.StatusInternalServerError, response.Result().StatusCode)
		assert.Equal(t, "Internal error\n", string(response.Body.String()))
	})
}
