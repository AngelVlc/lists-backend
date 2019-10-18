package controllers

import (
	"encoding/json"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
	"net/http"
)

// UsersHandler is the handler for the users endpoints
func UsersHandler(r *http.Request, servicePrv services.ServiceProvider, jwtInfo *models.JwtClaimsInfo) handlerResult {
	switch r.Method {
	case http.MethodPost:
		return processUsersGET(r, servicePrv)
	default:
		return okResult{nil, http.StatusMethodNotAllowed}
	}
}

func processUsersGET(r *http.Request, servicePrv services.ServiceProvider) handlerResult {
	dto, err := parseUserBody(r)
	if err != nil {
		return errorResult{err}
	}
	userSrv := servicePrv.GetUsersService()
	id, err := userSrv.AddUser(&dto)
	if err != nil {
		return errorResult{err}
	}
	return okResult{id, http.StatusCreated}
}

func parseUserBody(r *http.Request) (models.UserDto, error) {
	decoder := json.NewDecoder(r.Body)
	var dto models.UserDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.UserDto{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return dto, nil
}
