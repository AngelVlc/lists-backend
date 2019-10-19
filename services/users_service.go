package services

import (
	"fmt"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"gopkg.in/mgo.v2/bson"
)

// UsersService is the interface a users service must implement
type UsersService interface {
	AddUser(dto *models.UserDto) (string, error)
	CheckIfUserPasswordIsOk(userName string, password string) (*models.User, error)
	GetSingleUser(id string, u *models.User) error
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

// CheckIfUserPasswordIsOk returns nil if the password is correct or an error if it isn't
func (s *MyUsersService) CheckIfUserPasswordIsOk(userName string, password string) (*models.User, error) {
	foundUser, err := s.getUserByUserName(userName)
	if err != nil {
		return nil, err
	}

	if foundUser == nil {
		return nil, &appErrors.BadRequestError{Msg: "The user does not exist", InternalError: nil}
	}

	err = s.bcryptPrv.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password))
	if err != nil {
		return nil, &appErrors.BadRequestError{Msg: "Invalid password", InternalError: nil}
	}

	return foundUser, nil
}

// GetSingleUser returns a single user from its id
func (s *MyUsersService) GetSingleUser(id string, u *models.User) error {
	if !s.usersRepository().IsValidID(id) {
		return s.getInvalidIDError(id)
	}

	return s.usersRepository().GetOne(u, bson.D{{"_id", id}}, nil)
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

func (s *MyUsersService) getUserByUserName(userName string) (*models.User, error) {
	foundUsers := []models.User{}
	err := s.usersRepository().Get(&foundUsers, bson.M{"userName": userName}, nil)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error checking if user name exists", InternalError: err}
	}

	if len(foundUsers) == 0 {
		return nil, nil
	}

	return &foundUsers[0], nil
}

func (s *MyUsersService) getInvalidIDError(id string) error {
	return &appErrors.BadRequestError{Msg: fmt.Sprintf("%q is not a valid id", id), InternalError: nil}
}
