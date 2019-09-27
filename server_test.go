package main

import (
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
	server := newServer(nil)

	t.Run("handles /users", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/users", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Result().StatusCode, "status are not equal")
	})

	t.Run("handles /users/id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/users/wadus", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Result().StatusCode, "status are not equal")
	})

	t.Run("handles /lists", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/lists", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Result().StatusCode, "status are not equal")
	})

	t.Run("handles /lists/id", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPatch, "/lists/wadus", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusMethodNotAllowed, response.Result().StatusCode, "status are not equal")
	})

}
