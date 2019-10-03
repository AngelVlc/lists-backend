package services

import (
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
)

// UsersService is the interface a users service must implement
type UsersService interface {
	AddUser(u *models.User) error
}

// MyUsersService is the service for the users entity
type MyUsersService struct {
	session stores.MongoSession
}

// NewMyUsersService returns a new users service
func NewMyUsersService(session stores.MongoSession) *MyUsersService {
	return &MyUsersService{
		session: session,
	}
}

// AddUser  adds a user
func (s *MyUsersService) AddUser(u *models.User) error {
	return s.usersRepository().Add(u)
}

func (s *MyUsersService) usersRepository() stores.Repository {
	return s.session.GetRepository("users")
}
