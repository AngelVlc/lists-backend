package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/AngelVlc/lists-backend/services"
	"github.com/AngelVlc/lists-backend/stores"
)

func main() {
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%v", port)

	jwtSecret := os.Getenv("JWT_SECRET")

	ms := stores.NewMyMongoSession(false)

	bp := services.NewMyBcryptProvider()

	jwtp := services.NewMyJwtProvider(jwtSecret)

	sp := services.NewMyServiceProvider(ms, bp, jwtp)

	server := newServer(sp)

	log.Printf("Listening on port %v ...\n", port)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("could not listen on port %v %v", port, err)
	}
}
