package services

import (
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"gopkg.in/mgo.v2/bson"
)

// UsersService is the interface a users service must implement
type UsersService interface {
	AddUser(dto *models.UserDto) (string, error)
}

// MyUsersService is the service for the users entity
type MyUsersService struct {
	session   stores.MongoSession
	bcryptPrv BcryptProvider
}

var bcryptCost int = 3

// NewMyUsersService returns a new users service
func NewMyUsersService(session stores.MongoSession, bcryptPrv BcryptProvider) *MyUsersService {
	return &MyUsersService{
		session:   session,
		bcryptPrv: bcryptPrv,
	}
}

// AddUser  adds a user
func (s *MyUsersService) AddUser(dto *models.UserDto) (string, error) {
	if dto.NewPassword != dto.ConfirmNewPassword {
		return "", &appErrors.BadRequestError{Msg: "Passwords don't match", InternalError: nil}
	}

	userExists, err := s.existsUser(dto.UserName)
	if err != nil {
		return "", err
	}

	if userExists {
		return "", &appErrors.BadRequestError{Msg: "A user with the same user name already exists", InternalError: nil}
	}

	user := dto.ToUser()

	hasshedPass, err := s.bcryptPrv.GenerateFromPassword([]byte(dto.NewPassword), bcryptCost)
	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error encrypting password", InternalError: err}
	}

	user.PasswordHash = string(hasshedPass)

	return s.usersRepository().Add(&user)
}

func (s *MyUsersService) usersRepository() stores.Repository {
	return s.session.GetRepository("users")
}

func (s *MyUsersService) existsUser(userName string) (bool, error) {
	existingUsers := []models.GetUsersResultDto{}
	err := s.usersRepository().Get(&existingUsers, bson.M{"userName": userName}, bson.M{"_id": 1})
	if err != nil {
		return false, &appErrors.UnexpectedError{Msg: "Error checking if user name exists", InternalError: err}
	}

	return len(existingUsers) > 0, nil
}
