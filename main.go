package main

import (
	"fmt"
	"github.com/AngelVlc/lists-backend/services"
	"github.com/AngelVlc/lists-backend/stores"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%v", port)

	ms := stores.NewMyMongoSession(false)

	bp := services.NewMyBcryptProvider()

	sp := services.NewMyServiceProvider(ms, bp)

	server := newServer(sp)

	log.Printf("Listening on port %v ...\n", port)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("could not listen on port %v %v", port, err)
	}
}
