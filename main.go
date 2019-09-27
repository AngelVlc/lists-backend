package main

import (
	"fmt"
	"github.com/AngelVlc/lists-backend/stores"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%v", port)

	session := stores.NewMyMongoSession(false)
	store := stores.NewMongoRepository(session)
	server := newServer(&store)

	log.Printf("Listening on port %v ...\n", port)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("could not listen on port %v %v", port, err)
	}
}
