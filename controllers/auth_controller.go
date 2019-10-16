package controllers

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/services"
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
	var action = ""
	if len(r.URL.Path) > len("/auth") {
		action = r.URL.Path[len("/auth/"):]
	}

	if action == "token" {
		l, err := parseAuthBody(r)
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

	return errorResult{&appErrors.UnexpectedError{Msg: "Not implemented", InternalError: nil}}
}

func parseAuthBody(r *http.Request) (models.Login, error) {
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
