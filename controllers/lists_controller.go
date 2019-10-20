package controllers

import (
	"encoding/json"
	"net/http"
	"net/url"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
)

// ListsHandler is the handler for the lists endpoints
func ListsHandler(r *http.Request, servicePrv services.ServiceProvider, jwtInfo *models.JwtClaimsInfo) handlerResult {
	switch r.Method {
	case http.MethodGet:
		return processListsGET(r, servicePrv, jwtInfo.ID)
	case http.MethodPost:
		return processListsPOST(r, servicePrv, jwtInfo.ID)
	case http.MethodDelete:
		return processListsDELETE(r, servicePrv, jwtInfo.ID)
	case http.MethodPut:
		return processListsPUT(r, servicePrv, jwtInfo.ID)
	default:
		return okResult{nil, http.StatusMethodNotAllowed}
	}
}

func processListsGET(r *http.Request, servicePrv services.ServiceProvider, userID string) handlerResult {
	listID := getListIDFromURL(r.URL)

	listSrv := servicePrv.GetListsService()
	if listID == "" {
		r := []models.GetListsResultDto{}
		err := listSrv.GetUserLists(userID, &r)
		if err != nil {
			return errorResult{err}
		}
		return okResult{r, http.StatusOK}
	}
	l := models.List{}
	err := listSrv.GetSingleUserList(listID, userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func processListsPOST(r *http.Request, servicePrv services.ServiceProvider, userID string) handlerResult {
	l, err := parseListBody(r)

	if err != nil {
		return errorResult{err}
	}
	listSrv := servicePrv.GetListsService()

	id, err := listSrv.AddUserList(userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func processListsPUT(r *http.Request, servicePrv services.ServiceProvider, userID string) handlerResult {
	listID := getListIDFromURL(r.URL)

	l, err := parseListBody(r)
	if err != nil {
		return errorResult{err}
	}
	listSrv := servicePrv.GetListsService()
	err = listSrv.UpdateUserList(listID, userID, &l)
	if err != nil {
		return errorResult{err}
	}
	return okResult{l, http.StatusOK}
}

func processListsDELETE(r *http.Request, servicePrv services.ServiceProvider, userID string) handlerResult {
	listID := getListIDFromURL(r.URL)

	listSrv := servicePrv.GetListsService()
	err := listSrv.RemoveUserList(listID, userID)
	if err != nil {
		return errorResult{err}
	}
	return okResult{nil, http.StatusNoContent}
}

func getListIDFromURL(u *url.URL) string {
	var listID string

	if len(u.Path) > len("/lists") {
		listID = u.Path[len("/lists/"):]
	}

	return listID
}

func parseListBody(r *http.Request) (models.List, error) {
	decoder := json.NewDecoder(r.Body)
	var dto models.ListDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.List{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	l := dto.ToList()

	return l, nil
}
