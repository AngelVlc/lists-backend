package controllers

import (
	"encoding/json"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
	"net/http"
)

// UsersHandler is the handler for the users endpoints
func UsersHandler(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
	switch r.Method {
	case http.MethodPost:
		dto, err := parseUserBody(r)
		if err != nil {
			return errorResult{err}
		}
		userSrv := serviceProvider.GetUsersService()
		id, err := userSrv.AddUser(&dto)
		if err != nil {
			return errorResult{err}
		}
		return okResult{id, http.StatusCreated}
	default:
		return okResult{nil, http.StatusMethodNotAllowed}
	}

}

func parseUserBody(r *http.Request) (models.UserDto, error) {
	if r.Body == nil {
		return models.UserDto{}, &appErrors.BadRequestError{Msg: "No body"}
	}
	decoder := json.NewDecoder(r.Body)
	var dto models.UserDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.UserDto{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	return dto, nil
}
