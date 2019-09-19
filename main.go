package main

import (
	"fmt"
	"github.com/AngelVlc/lists-backend/stores"
	"log"
	"net/http"
)

func main() {
	port := 5000
	addr := fmt.Sprintf(":%v", port)

	session := stores.NewMyMongoSession("mongodb://mongo", "listsDb", "lists")
	store := stores.NewMongoStore(session)
	server := newServer(&store)

	log.Printf("Listening on port %v ...\n", port)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("could not listen on port %v %v", port, err)
	}
}
