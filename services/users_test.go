package services

import (
	"errors"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

		hasshedPass := "hashedPass"
		mockedBcryptProvider.On("GenerateFromPassword", []byte(dto.NewPassword), bcryptCost).Return([]byte(hasshedPass), nil).Once()
		u := dto.ToUser()
		u.PasswordHash = string(hasshedPass)
		mockedRepository.On("Add", &u).Return(errors.New("error"))

		err := service.AddUser(&dto)

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

		err := service.AddUser(&dto)

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

		mockedBcryptProvider.On("GenerateFromPassword", []byte(dto.NewPassword), bcryptCost).Return([]byte(""), errors.New("wadus")).Once()

		err := service.AddUser(&dto)

		assert.NotNil(t, err)

		unexpectErr, isError := err.(*appErrors.UnexpectedError)
		assert.Equal(t, true, isError, "should be a bad request error")
		assert.Equal(t, "Error encrypting password", unexpectErr.Error())
	})
}
