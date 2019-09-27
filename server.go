package main

import (
	"github.com/AngelVlc/lists-backend/controllers"
	"github.com/AngelVlc/lists-backend/controllers/lists"
	"github.com/AngelVlc/lists-backend/controllers/users"
	"github.com/AngelVlc/lists-backend/stores"
	"net/http"
)

type server struct {
	store stores.Store
	http.Handler
}

func newServer(store stores.Store) *server {
	s := new(server)
	s.store = store

	router := http.NewServeMux()
	router.Handle("/lists", s.getHandler(lists.Handler))
	router.Handle("/lists/", s.getHandler(lists.Handler))
	router.Handle("/users", s.getHandler(users.Handler))
	router.Handle("/users/", s.getHandler(users.Handler))

	s.Handler = router

	return s
}

func (s *server) getHandler(handlerFunc controllers.HandlerFunc) controllers.Handler {
	return controllers.Handler{
		HandlerFunc: lists.Handler,
		Store:       s.store,
	}
}
