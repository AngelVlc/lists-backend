package main

import (
	"encoding/json"
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"log"
	"net/http"
	"net/url"
)

const jsonContentType = "application/json"

type server struct {
	store stores.Store
	http.Handler
}

func (s *server) getListsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %q", r.Method, r.URL)

	listID := getListIDFromURL(r.URL)

	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
	}

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("content-type", jsonContentType)
		json.NewEncoder(w).Encode(s.store.GetLists())
	case http.MethodPost:
		l, err := parseListBody(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			err = s.store.AddList(&l)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
			} else {
				w.Header().Set("content-type", jsonContentType)
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(l)
			}
		}
	case http.MethodDelete:
		err = s.store.RemoveList(listID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func newServer(store stores.Store) *server {
	s := new(server)
	s.store = store

	router := http.NewServeMux()
	router.Handle("/lists", http.HandlerFunc(s.getListsHandler))
	router.Handle("/lists/", http.HandlerFunc(s.getListsHandler))

	s.Handler = router

	return s
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
		log.Println("Error parsing body: missing body")
		return models.List{}, errors.New("error parsing body")
	}
	decoder := json.NewDecoder(r.Body)
	var dto models.ListDto
	err := decoder.Decode(&dto)
	if err != nil {
		log.Printf("Error parsing body: %v", err)
		return models.List{}, errors.New("error parsing body")
	}

	l := models.List{
		Name:  dto.Name,
		Items: dto.Items,
	}

	return l, nil
}
