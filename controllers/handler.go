package controllers

import (
	"encoding/json"
	"github.com/AngelVlc/lists-backend/services"
	"github.com/AngelVlc/lists-backend/stores"
	"log"
	"net/http"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	ServiceProvider services.ServiceProvider
}

type handlerResult interface {
	IsError() bool
}

type errorResult struct {
	err error
}

func (e errorResult) IsError() bool {
	return true
}

type okResult struct {
	content    interface{}
	statusCode int
}

func (r okResult) IsError() bool {
	return false
}

// HandlerFunc is the type for the handler functions
type HandlerFunc func(http.ResponseWriter, *http.Request, services.ServiceProvider) handlerResult

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %q", r.Method, r.URL)

	res := h.HandlerFunc(w, r, h.ServiceProvider)

	if res.IsError() {
		errorRes, _ := res.(errorResult)
		err := errorRes.err
		if unexErr, ok := err.(*stores.UnexpectedError); ok {
			writeErrorResponse(w, http.StatusInternalServerError, unexErr.Error(), unexErr.InternalError)
		} else if notFoundErr, ok := err.(*stores.NotFoundError); ok {
			writeErrorResponse(w, http.StatusNotFound, notFoundErr.Error(), nil)
		} else if badIDErr, ok := err.(*stores.InvalidIDError); ok {
			writeErrorResponse(w, http.StatusBadRequest, badIDErr.Error(), nil)
		} else if noBodyErr, ok := err.(*NoBodyError); ok {
			writeErrorResponse(w, http.StatusBadRequest, noBodyErr.Error(), nil)
		} else if badBodyErr, ok := err.(*InvalidBodyError); ok {
			writeErrorResponse(w, http.StatusBadRequest, badBodyErr.Error(), badBodyErr.InternalError)
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, "Internal error", err)
		}

	} else {
		okRes, _ := res.(okResult)
		writeOkResponse(w, okRes.statusCode, okRes.content)
	}
}

// writeErrorResponse is used when and endpoind responds with an error
func writeErrorResponse(w http.ResponseWriter, statusCode int, msg string, internalError error) {
	if internalError != nil {
		log.Printf("%v %v", statusCode, internalError)
	} else {
		log.Printf("%v %v", statusCode, msg)
	}
	http.Error(w, msg, statusCode)
}

// writeOkResponse is used when and endpoind does not respond with an error
func writeOkResponse(w http.ResponseWriter, statusCode int, content interface{}) {
	log.Println(statusCode)

	const jsonContentType = "application/json"

	if content != nil {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(content)
	} else {
		w.WriteHeader(statusCode)
	}
}
