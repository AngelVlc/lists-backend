package controllers

import (
	"encoding/json"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
	"net/http"
	"net/url"
)

// ListsHandler is the handler for the lists endpoints
func ListsHandler(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
	listID := getListIDFromURL(r.URL)

	switch r.Method {
	case http.MethodGet:
		listSrv := serviceProvider.GetListsService()
		if listID == "" {
			r := []models.GetListsResultDto{}
			err := listSrv.GetLists(&r)
			if err != nil {
				return errorResult{err}
			}
			return okResult{r, http.StatusOK}
		}
		l := models.List{}
		err := listSrv.GetSingleList(listID, &l)
		if err != nil {
			return errorResult{err}
		}
		return okResult{l, http.StatusOK}
	case http.MethodPost:
		l, err := parseListBody(r)
		if err != nil {
			return errorResult{err}
		}
		listSrv := serviceProvider.GetListsService()
		id, err := listSrv.AddList(&l)
		if err != nil {
			return errorResult{err}
		}
		return okResult{id, http.StatusCreated}
	case http.MethodDelete:
		listSrv := serviceProvider.GetListsService()
		err := listSrv.RemoveList(listID)
		if err != nil {
			return errorResult{err}
		}
		return okResult{nil, http.StatusNoContent}
	case http.MethodPut:
		l, err := parseListBody(r)
		if err != nil {
			return errorResult{err}
		}
		listSrv := serviceProvider.GetListsService()
		err = listSrv.UpdateList(listID, &l)
		if err != nil {
			return errorResult{err}
		}
		return okResult{l, http.StatusOK}
	default:
		return okResult{nil, http.StatusMethodNotAllowed}
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
		return models.List{}, &appErrors.BadRequestError{Msg: "No body"}
	}
	decoder := json.NewDecoder(r.Body)
	var dto models.ListDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.List{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	l := dto.ToList()

	return l, nil
}
