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
	router.Handle("/lists", controllers.Handler{lists.Handler, s.store})
	router.Handle("/lists/", controllers.Handler{lists.Handler, s.store})
	router.Handle("/users", controllers.Handler{users.Handler, s.store})
	router.Handle("/users/", controllers.Handler{users.Handler, s.store})

	s.Handler = router

	return s
}
