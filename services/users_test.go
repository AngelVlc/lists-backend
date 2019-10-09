package services

import (
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserService(t *testing.T) {
	mockedSession := new(mockedMongoSession)
	service := NewMyUsersService(mockedSession)

	mockedRepository := new(mockedRepository)

	mockedSession.On("GetRepository", "users").Return(mockedRepository)

	t.Run("AddUser() should call repository.AddUser", func(t *testing.T) {
		dto := models.UserDto{
			UserName: "user",
		}

		u := dto.ToUser()

		mockedRepository.On("Add", &u).Return(errors.New("error"))

		err := service.AddUser(&dto)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})
}
