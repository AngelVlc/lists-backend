package controllers

import (
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"net/http"
	"encoding/json"
)

// AuthHandler is the handler for the auth endpoints
func AuthHandler(r *http.Request, servicePrv services.ServiceProvider) handlerResult {
	switch r.Method {
	case http.MethodPost:
		return processAuthPOST(r, servicePrv)
	default:
		return okResult{nil, http.StatusMethodNotAllowed}
	}
}

func processAuthPOST(r *http.Request, servicePrv services.ServiceProvider) handlerResult {
	_, err := parseAuthBody(r)

	if err != nil {
		return errorResult{err}
	}

	return okResult{nil, http.StatusOK}
}

func parseAuthBody(r *http.Request) (models.Login, error) {
	if r.Body == nil {
		return models.Login{}, &appErrors.BadRequestError{Msg: "No body"}
	}
	decoder := json.NewDecoder(r.Body)

	var l models.Login
	err := decoder.Decode(&l)
	if err != nil {
		return models.Login{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return l, nil
}