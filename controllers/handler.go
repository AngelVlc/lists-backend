package controllers

import (
	"encoding/json"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/services"
	"log"
	"net/http"
	"strings"
)

// Handler is the type used to handle the endpoints
type Handler struct {
	HandlerFunc
	ServiceProvider services.ServiceProvider
	RequireAuth     bool
	RequireAdmin    bool
}

type handlerResult interface {
	IsError() bool
}

type errorResult struct {
	err error
}

func (e errorResult) IsError() bool {
	return true
}

type okResult struct {
	content    interface{}
	statusCode int
}

func (r okResult) IsError() bool {
	return false
}

// HandlerFunc is the type for the handler functions
type HandlerFunc func(*http.Request, services.ServiceProvider) handlerResult

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%v %q", r.Method, r.URL)

	if h.RequireAuth {
		token, err := getAuthToken(r)
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		authSrv := h.ServiceProvider.GetAuthService()
		jwtInfo, err := authSrv.ValidateToken(token)

		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "Invalid auth token", err)
			return
		}

		if h.RequireAdmin && !jwtInfo.IsAdmin {
			writeErrorResponse(w, http.StatusForbidden, "Access forbidden", err)
			return
		}
	}

	res := h.HandlerFunc(r, h.ServiceProvider)

	if res.IsError() {
		errorRes, _ := res.(errorResult)
		err := errorRes.err
		if unexErr, ok := err.(*appErrors.UnexpectedError); ok {
			writeErrorResponse(w, http.StatusInternalServerError, unexErr.Error(), unexErr.InternalError)
		} else if unauthErr, ok := err.(*appErrors.UnauthorizedError); ok {
			writeErrorResponse(w, http.StatusUnauthorized, unauthErr.Error(), unauthErr.InternalError)
		} else if notFoundErr, ok := err.(*appErrors.NotFoundError); ok {
			writeErrorResponse(w, http.StatusNotFound, notFoundErr.Error(), nil)
		} else if badRequestErr, ok := err.(*appErrors.BadRequestError); ok {
			writeErrorResponse(w, http.StatusBadRequest, badRequestErr.Error(), badRequestErr.InternalError)
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, "Internal error", err)
		}
	} else {
		okRes, _ := res.(okResult)
		writeOkResponse(w, okRes.statusCode, okRes.content)
	}
}

// writeErrorResponse is used when and endpoind responds with an error
func writeErrorResponse(w http.ResponseWriter, statusCode int, msg string, internalError error) {
	if internalError != nil {
		log.Printf("%v %v [%v]", statusCode, msg, internalError)
	} else {
		log.Printf("%v %v", statusCode, msg)
	}
	http.Error(w, msg, statusCode)
}

// writeOkResponse is used when and endpoind does not respond with an error
func writeOkResponse(w http.ResponseWriter, statusCode int, content interface{}) {
	log.Println(statusCode)

	const jsonContentType = "application/json"

	if content != nil {
		w.Header().Set("content-type", jsonContentType)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(content)
	} else {
		w.WriteHeader(statusCode)
	}
}

func getAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if len(authHeader) == 0 {
		return "", &appErrors.UnauthorizedError{Msg: "No authorization header", InternalError: nil}
	}

	authHeaderParts := strings.Split(authHeader, "Bearer ")

	if len(authHeaderParts) != 2 {
		return "", &appErrors.UnauthorizedError{Msg: "Invalid authorization header", InternalError: nil}
	}

	return authHeaderParts[1], nil
}
