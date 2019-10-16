package controllers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
)

// TokenHandler is the handler for the auth/token endpoint
func TokenHandler(r *http.Request, servicePrv services.ServiceProvider) handlerResult {
	if r.Method != http.MethodPost {
		return okResult{nil, http.StatusMethodNotAllowed}
	}

	l, err := parseTokenBody(r)
	if err != nil {
		return errorResult{err}
	}

	userSrv := servicePrv.GetUsersService()
	foundUser, err := userSrv.CheckIfUserPasswordIsOk(l.UserName, l.Password)
	if err != nil {
		return errorResult{err}
	}

	authSrv := servicePrv.GetAuthService()

	tokens, err := authSrv.CreateTokens(foundUser)
	if err != nil {
		return errorResult{err}
	}

	return okResult{tokens, http.StatusOK}
}

func parseTokenBody(r *http.Request) (models.Login, error) {
	decoder := json.NewDecoder(r.Body)

	var l models.Login
	err := decoder.Decode(&l)
	if err != nil {
		return models.Login{}, &appErrors.BadRequestError{Msg: "Invalid body", InternalError: err}
	}

	if len(l.UserName) == 0 {
		return models.Login{}, &appErrors.BadRequestError{Msg: "UserName is mandatory", InternalError: nil}
	}

	if len(l.Password) == 0 {
		return models.Login{}, &appErrors.BadRequestError{Msg: "Password is mandatory", InternalError: nil}
	}

	return l, nil
}
