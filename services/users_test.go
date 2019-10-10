package services

import (
	"errors"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type mockedBcryptProvider struct {
	mock.Mock
}

func (m *mockedBcryptProvider) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	args := m.Called(password, cost)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *mockedBcryptProvider) CompareHashAndPassword(hashedPassword, password []byte) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func TestUserService(t *testing.T) {
	mockedSession := new(mockedMongoSession)

	mockedBcryptProvider := new(mockedBcryptProvider)

	service := NewMyUsersService(mockedSession, mockedBcryptProvider)

	mockedRepository := new(mockedRepository)

	mockedSession.On("GetRepository", "users").Return(mockedRepository)

	t.Run("AddUser() should call repository.AddUser if the dto is valid", func(t *testing.T) {
		dto := models.UserDto{
			UserName:           "user",
			NewPassword:        "pass",
			ConfirmNewPassword: "pass",
		}

		r := []models.GetUsersResultDto{}
		mockedRepository.On("Get", &r, bson.M{"userName": dto.UserName}, bson.M{"_id": 1}).Return(nil).Once()

		hasshedPass := "hashedPass"
		mockedBcryptProvider.On("GenerateFromPassword", []byte(dto.NewPassword), bcryptCost).Return([]byte(hasshedPass), nil).Once()
		u := dto.ToUser()
		u.PasswordHash = string(hasshedPass)
		mockedRepository.On("Add", &u).Return("", errors.New("error"))

		id, err := service.AddUser(&dto)

		assert.Empty(t, id)
		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
		mockedBcryptProvider.AssertExpectations(t)
	})

	t.Run("AddUser() should return a BadRequestError when the passwords does not match", func(t *testing.T) {
		dto := models.UserDto{
			UserName:           "user",
			NewPassword:        "pass",
			ConfirmNewPassword: "other",
		}

		id, err := service.AddUser(&dto)

		assert.Empty(t, id)
		assert.NotNil(t, err)

		badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isBadReqErr, "should be a bad request error")
		assert.Equal(t, "Passwords don't match", badReqErr.Error())
	})

	t.Run("AddUser() should return an UnexpectedError when bcrypt fails", func(t *testing.T) {
		dto := models.UserDto{
			UserName:           "user",
			NewPassword:        "pass",
			ConfirmNewPassword: "pass",
		}

		r := []models.GetUsersResultDto{}
		mockedRepository.On("Get", &r, bson.M{"userName": dto.UserName}, bson.M{"_id": 1}).Return(nil).Once()

		mockedBcryptProvider.On("GenerateFromPassword", []byte(dto.NewPassword), bcryptCost).Return([]byte(""), errors.New("wadus")).Once()

		id, err := service.AddUser(&dto)

		assert.Empty(t, id)
		assert.NotNil(t, err)

		unexpectErr, isError := err.(*appErrors.UnexpectedError)
		assert.Equal(t, true, isError, "should be a bad request error")
		assert.Equal(t, "Error encrypting password", unexpectErr.Error())
	})

	t.Run("AddUser() should return a BadRequestError if a user with the same name exists", func(t *testing.T) {
		dto := models.UserDto{
			UserName:           "user",
			NewPassword:        "pass",
			ConfirmNewPassword: "pass",
		}

		item := models.GetUsersResultDto{
			ID: "id",
		}
		r := []models.GetUsersResultDto{item}
		mockedRepository.On("Get", &[]models.GetUsersResultDto{}, bson.M{"userName": dto.UserName}, bson.M{"_id": 1}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]models.GetUsersResultDto)
			*arg = r
		})

		id, err := service.AddUser(&dto)

		assert.Empty(t, id)
		assert.NotNil(t, err)

		badReqErr, isBadReqErr := err.(*appErrors.BadRequestError)
		assert.Equal(t, true, isBadReqErr, "should be a bad request error")
		assert.Equal(t, "A user with the same user name already exists", badReqErr.Error())
	})
}
