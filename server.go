package main

import (
	"encoding/json"
	"fmt"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"log"
	"net/http"
	"net/url"
)

// InvalidBodyError occurs when the body's request is invalid
type InvalidBodyError struct {
	InternalError error
}

func (e *InvalidBodyError) Error() string {
	return fmt.Sprint("Invalid body")
}

// NoBodyError occurs when the request does not have body
type NoBodyError struct{}

func (e *NoBodyError) Error() string {
	return fmt.Sprint("No body")
}

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %q", r.Method, r.URL)

	if err := fn(w, r); err != nil {
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
	}
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, msg string, internalError error) {
	if internalError != nil {
		log.Printf("%v %v", statusCode, internalError)
	} else {
		log.Printf("%v %v", statusCode, msg)
	}
	http.Error(w, msg, statusCode)
}

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

type server struct {
	store stores.Store
	http.Handler
}

func newServer(store stores.Store) *server {
	s := new(server)
	s.store = store

	router := http.NewServeMux()
	router.Handle("/lists", appHandler(s.listsHandler))
	router.Handle("/lists/", appHandler(s.listsHandler))

	s.Handler = router

	return s
}

func (s *server) listsHandler(w http.ResponseWriter, r *http.Request) error {
	listID := getListIDFromURL(r.URL)

	switch r.Method {
	case http.MethodGet:
		if listID == "" {
			r, err := s.store.GetLists()
			fmt.Println("####1")
			if err != nil {
				fmt.Println("####2")
				return err
			}
			fmt.Println("####3")
			writeOkResponse(w, http.StatusOK, r)
		} else {
			l, err := s.store.GetSingleList(listID)
			if err != nil {
				return err
			}
			writeOkResponse(w, http.StatusOK, l)
		}
	case http.MethodPost:
		l, err := parseListBody(r)
		if err != nil {
			return err
		}
		err = s.store.AddList(&l)
		if err != nil {
			return err
		}
		writeOkResponse(w, http.StatusCreated, l)
	case http.MethodDelete:
		err := s.store.RemoveList(listID)
		if err != nil {
			return err
		}
		writeOkResponse(w, http.StatusNoContent, nil)
	case http.MethodPut:
		l, err := parseListBody(r)
		if err != nil {
			return err
		}
		err = s.store.UpdateList(listID, &l)
		if err != nil {
			return err
		}
		writeOkResponse(w, http.StatusOK, l)
	default:
		writeOkResponse(w, http.StatusMethodNotAllowed, nil)
	}

	return nil
}

func getListIDFromURL(u *url.URL) string {
	var listID string

	if len(u.Path) > len("/lists") {
		listID = u.Path[len("/lists/"):]
	}

	return listID
}

func parseListBody(r *http.Request) (models.List, error) {
	if r.Body == nil {
		return models.List{}, &NoBodyError{}
	}
	decoder := json.NewDecoder(r.Body)
	var dto models.ListDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.List{}, &InvalidBodyError{InternalError: err}
	}

	l := dto.ToList()

	return l, nil
}
