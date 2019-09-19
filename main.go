package main

import (
	"fmt"
	"github.com/AngelVlc/lists-backend/stores"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

func main() {
	port := 5000
	addr := fmt.Sprintf(":%v", port)

	session := stores.NewMyMongoSession(mongoSession(), "listsDb", "lists")

	store := stores.NewMongoStore(session)
	server := newServer(&store)

	log.Printf("Listening on port %v ...\n", port)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("could not listen on port %v %v", port, err)
	}
}

func mongoSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://mongo")
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	log.Println("Connected with lists mongo database.")

	return s
}
