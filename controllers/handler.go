package controllers

import (
	"net/http"
	"log"
	"github.com/AngelVlc/lists-backend/stores"
	"encoding/json"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	Store stores.Store
}

// HandlerFunc is the type for the handler functions
type HandlerFunc func(http.ResponseWriter, *http.Request, stores.Store) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %q", r.Method, r.URL)

	if err := h.HandlerFunc(w, r, h.Store); err != nil {
		if unexErr, ok := err.(*stores.UnexpectedError); ok {
			WriteErrorResponse(w, http.StatusInternalServerError, unexErr.Error(), unexErr.InternalError)
		} else if notFoundErr, ok := err.(*stores.NotFoundError); ok {
			WriteErrorResponse(w, http.StatusNotFound, notFoundErr.Error(), nil)
		} else if badIDErr, ok := err.(*stores.InvalidIDError); ok {
			WriteErrorResponse(w, http.StatusBadRequest, badIDErr.Error(), nil)
		} else if noBodyErr, ok := err.(*NoBodyError); ok {
			WriteErrorResponse(w, http.StatusBadRequest, noBodyErr.Error(), nil)
		} else if badBodyErr, ok := err.(*InvalidBodyError); ok {
			WriteErrorResponse(w, http.StatusBadRequest, badBodyErr.Error(), badBodyErr.InternalError)
		} else {
			WriteErrorResponse(w, http.StatusInternalServerError, "Internal error", err)
		}
	}
}

// WriteErrorResponse is used when and endpoind responds with an error
func WriteErrorResponse(w http.ResponseWriter, statusCode int, msg string, internalError error) {
	if internalError != nil {
		log.Printf("%v %v", statusCode, internalError)
		} else {
			log.Printf("%v %v", statusCode, msg)
		}
		http.Error(w, msg, statusCode)
	}
	
// WriteOkResponse is used when and endpoind does not respond with an error
func WriteOkResponse(w http.ResponseWriter, statusCode int, content interface{}) {
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