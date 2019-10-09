package services

import (
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListsService(t *testing.T) {
	mockedSession := new(mockedMongoSession)
	service := NewMyListsService(mockedSession)

	mockedRepository := new(mockedRepository)

	mockedSession.On("GetRepository", "lists").Return(mockedRepository)

	t.Run("AddList() should call repository.AddList", func(t *testing.T) {
		l := models.List{
			ID:   "1",
			Name: "list",
		}

		mockedRepository.On("Add", &l).Return("", errors.New("error")).Once()

		id, err := service.AddList(&l)

		assert.Empty(t, id)
		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("RemoveList() should call repository.RemoveList", func(t *testing.T) {
		mockedRepository.On("Remove", "id").Return(errors.New("error")).Once()

		err := service.RemoveList("id")

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("UpdateList() should call repository.RemoveList", func(t *testing.T) {
		l := models.List{
			ID:   "1",
			Name: "list",
		}

		mockedRepository.On("Update", "id", &l).Return(errors.New("error")).Once()

		err := service.UpdateList("id", &l)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("GetSingleList() should call repository.RemoveList", func(t *testing.T) {
		l := models.List{}

		mockedRepository.On("GetSingle", "id", &l).Return(errors.New("error")).Once()

		err := service.GetSingleList("id", &l)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("GetLists() should call repository.RemoveList", func(t *testing.T) {
		r := []models.GetListsResultDto{}

		mockedRepository.On("Get", &r).Return(errors.New("error")).Once()

		err := service.GetLists(&r)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})
}
