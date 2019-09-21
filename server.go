package main

import (
	"encoding/json"
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"net/url"
)

const jsonContentType = "application/json"

type server struct {
	store stores.Store
	http.Handler
}

func newServer(store stores.Store) *server {
	s := new(server)
	s.store = store

	router := http.NewServeMux()
	router.Handle("/lists", http.HandlerFunc(s.listsHandler))
	router.Handle("/lists/", http.HandlerFunc(s.listsHandler))

	s.Handler = router

	return s
}

func (s *server) listsHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %q", r.Method, r.URL)

	listID := getListIDFromURL(r.URL)

	if listID != "" && !bson.IsObjectIdHex(listID) {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if listID == "" {
			w.Header().Set("content-type", jsonContentType)
			json.NewEncoder(w).Encode(s.store.GetLists())
		} else {
			l, err := s.store.GetSingleList(listID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				w.Header().Set("content-type", jsonContentType)
				json.NewEncoder(w).Encode(l)
			}
		}
	case http.MethodPost:
		l, err := parseListBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			err = s.store.AddList(&l)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				w.Header().Set("content-type", jsonContentType)
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(l)
			}
		}
	case http.MethodDelete:
		err := s.store.RemoveList(listID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
	case http.MethodPut:
		l, err := parseListBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			err = s.store.UpdateList(listID, &l)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				w.Header().Set("content-type", jsonContentType)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(l)
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
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

	l := dto.ToList()

	return l, nil
}
