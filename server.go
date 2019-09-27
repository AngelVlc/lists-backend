package main

import (
	"github.com/AngelVlc/lists-backend/controllers"
	"github.com/AngelVlc/lists-backend/stores"
	"net/http"
)

type server struct {
	store stores.Repository
	http.Handler
}

func newServer(store stores.Repository) *server {
	s := new(server)
	s.store = store

	router := http.NewServeMux()
	router.Handle("/lists", s.getHandler(controllers.ListsHandler))
	router.Handle("/lists/", s.getHandler(controllers.ListsHandler))
	router.Handle("/users", s.getHandler(controllers.UsersHandler))
	router.Handle("/users/", s.getHandler(controllers.UsersHandler))

	s.Handler = router

	return s
}

func (s *server) getHandler(handlerFunc controllers.HandlerFunc) controllers.Handler {
	return controllers.Handler{
		HandlerFunc: handlerFunc,
		Repository:  s.store,
	}
}
