package controllers

import (
	"encoding/json"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
	"net/http"
)

// UsersHandler is the handler for the users endpoints
func UsersHandler(w http.ResponseWriter, r *http.Request, serviceProvider services.ServiceProvider) handlerResult {
	switch r.Method {
	case http.MethodPost:
		u, err := parseUserBody(r)
		if err != nil {
			return errorResult{err}
		}
		userSrv := serviceProvider.GetUsersService()
		err = userSrv.AddUser(&u)
		if err != nil {
			return errorResult{err}
		}
		return okResult{u, http.StatusCreated}
	default:
		return okResult{nil, http.StatusMethodNotAllowed}
	}

}

func parseUserBody(r *http.Request) (models.User, error) {
	if r.Body == nil {
		return models.User{}, &NoBodyError{}
	}
	decoder := json.NewDecoder(r.Body)
	var dto models.UserDto
	err := decoder.Decode(&dto)
	if err != nil {
		return models.User{}, &InvalidBodyError{InternalError: err}
	}

	l := dto.ToUser()

	return l, nil
}
