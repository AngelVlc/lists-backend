package controllers

import (
	"net/http"
	"net/url"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"encoding/json"
)

// ListsHandler is the handler for the lists endpoints
func ListsHandler(w http.ResponseWriter, r *http.Request, repository stores.Repository) error {
	listID := getListIDFromURL(r.URL)

	switch r.Method {
	case http.MethodGet:
		if listID == "" {
			r := []models.GetListsResultDto{}
			err := repository.Get(&r)
			if err != nil {
				return err
			}
			writeOkResponse(w, http.StatusOK, r)
		} else {
			l := models.List{}
			err := repository.GetSingle(listID, &l)
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
		err = repository.Add(&l)
		if err != nil {
			return err
		}
		writeOkResponse(w, http.StatusCreated, l)
	case http.MethodDelete:
		err := repository.Remove(listID)
		if err != nil {
			return err
		}
		writeOkResponse(w, http.StatusNoContent, nil)
	case http.MethodPut:
		l, err := parseListBody(r)
		if err != nil {
			return err
		}
		err = repository.Update(listID, &l)
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