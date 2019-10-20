package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
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

	checkAdminUser(sp)

	server := newServer(sp)

	log.Printf("Listening on port %v ...\n", port)

	if err := http.ListenAndServe(addr, server); err != nil {
		log.Fatalf("could not listen on port %v %v", port, err)
	}
}

func checkAdminUser(sp services.ServiceProvider) {
	us := sp.GetUsersService()

	u := models.User{}
	err := us.GetUserByUserName("admin", &u)

	if err == nil {
		log.Printf("Admin user already exists")
		return
	}

	if _, ok := err.(*appErrors.NotFoundError); ok {
		log.Printf("Admin user does not exist")

		n := models.UserDto{
			UserName:           "admin",
			NewPassword:        "admin",
			ConfirmNewPassword: "admin",
			IsAdmin:            true,
		}
		_, err = us.AddUser(&n)

		if err != nil {
			log.Fatalf("error creating admin user: %v", err)
		}

		log.Printf("Created admin user")
	}
}
