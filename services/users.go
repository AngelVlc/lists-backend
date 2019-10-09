package services

import (
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
)

// UsersService is the interface a users service must implement
type UsersService interface {
	AddUser(dto *models.UserDto) error
}

// MyUsersService is the service for the users entity
type MyUsersService struct {
	session   stores.MongoSession
	bcryptPrv *bcryptProvider
}

// NewMyUsersService returns a new users service
func NewMyUsersService(session stores.MongoSession) *MyUsersService {
	return &MyUsersService{
		session:   session,
		bcryptPrv: new(bcryptProvider),
	}
}

// AddUser  adds a user
func (s *MyUsersService) AddUser(dto *models.UserDto) error {
	user := dto.ToUser()
	return s.usersRepository().Add(&user)
}

func (s *MyUsersService) usersRepository() stores.Repository {
	return s.session.GetRepository("users")
}
