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
		u := models.User{
			ID:       "1",
			UserName: "user",
		}

		mockedRepository.On("Add", &u).Return(errors.New("error"))

		err := service.AddUser(&u)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})
}
