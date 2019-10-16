package main

import (
	"github.com/AngelVlc/lists-backend/controllers"
	"github.com/AngelVlc/lists-backend/services"
	"net/http"
)

type server struct {
	serviceProvider services.ServiceProvider
	http.Handler
}

func newServer(sp services.ServiceProvider) *server {
	s := new(server)
	s.serviceProvider = sp

	router := http.NewServeMux()

	router.Handle("/lists", s.getHandler(controllers.ListsHandler, true, false))
	router.Handle("/lists/", s.getHandler(controllers.ListsHandler, true, false))
	router.Handle("/users", s.getHandler(controllers.UsersHandler, true, true))
	router.Handle("/users/", s.getHandler(controllers.UsersHandler, true, true))
	router.Handle("/auth/token", s.getHandler(controllers.TokenHandler, false, false))
	router.Handle("/auth/refreshtoken", s.getHandler(controllers.RefreshTokenHandler, false, false))

	s.Handler = router

	return s
}

func (s *server) getHandler(handlerFunc controllers.HandlerFunc, requireAuth bool, requireAdmin bool) controllers.Handler {
	return controllers.Handler{
		HandlerFunc:     handlerFunc,
		ServiceProvider: s.serviceProvider,
		RequireAuth:     requireAuth,
		RequireAdmin:    requireAdmin,
	}
}
