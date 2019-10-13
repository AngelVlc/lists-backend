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

	router.Handle("/lists", s.getHandler(controllers.ListsHandler))
	router.Handle("/lists/", s.getHandler(controllers.ListsHandler))
	router.Handle("/users", s.getHandler(controllers.UsersHandler))
	router.Handle("/users/", s.getHandler(controllers.UsersHandler))
	router.Handle("/auth/", s.getHandler(controllers.AuthHandler))

	s.Handler = router

	return s
}

func (s *server) getHandler(handlerFunc controllers.HandlerFunc) controllers.Handler {
	return controllers.Handler{
		HandlerFunc:     handlerFunc,
		ServiceProvider: s.serviceProvider,
	}
}
