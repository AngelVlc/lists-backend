package main

import (
	"github.com/AngelVlc/lists-backend/services"
	"github.com/AngelVlc/lists-backend/stores"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	ms := stores.NewMyMongoSession(false)
	sp := services.NewMyServiceProvider(ms, nil, nil)
	server := newServer(sp)

	t.Run("handles /users", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/users", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode, "status are not equal")
	})

	t.Run("handles /users/id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/users/wadus", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode, "status are not equal")
	})

	t.Run("handles /lists", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode, "status are not equal")
	})

	t.Run("handles /lists/id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/lists/wadus", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Result().StatusCode, "status are not equal")
	})

	t.Run("handles /auth/token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/auth/token", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Result().StatusCode, "status are not equal")
	})

	t.Run("handles /auth/refreshtoken", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/auth/refreshtoken", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Result().StatusCode, "status are not equal")
	})

}
