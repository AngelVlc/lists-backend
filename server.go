package main

import (
	"github.com/AngelVlc/lists-backend/controllers"
	"github.com/AngelVlc/lists-backend/stores"
	"net/http"
)

type server struct {
	session stores.MongoSession
	http.Handler
}

func newServer(session stores.MongoSession) *server {
	s := new(server)

	router := http.NewServeMux()

	listsRepo := stores.NewMongoRepository(session, session.Collection("lists"))
	usersRepo := stores.NewMongoRepository(session, session.Collection("users"))

	router.Handle("/lists", s.getHandler(controllers.ListsHandler, listsRepo))
	router.Handle("/lists/", s.getHandler(controllers.ListsHandler, listsRepo))
	router.Handle("/users", s.getHandler(controllers.UsersHandler, usersRepo))
	router.Handle("/users/", s.getHandler(controllers.UsersHandler, usersRepo))

	s.Handler = router

	return s
}

func (s *server) getHandler(handlerFunc controllers.HandlerFunc, store stores.Repository) controllers.Handler {
	return controllers.Handler{
		HandlerFunc: handlerFunc,
		Repository:  store,
	}
}
