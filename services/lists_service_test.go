package services

import (
	"errors"
	"fmt"
	"testing"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
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
		mockedRepository.On("IsValidID", "id").Return(true).Once()

		err := service.RemoveList("id")

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("RemoveList() should return a badRequestError when the id is not valid", func(t *testing.T) {
		mockedRepository.On("IsValidID", "id").Return(false).Once()

		err := service.RemoveList("id")

		assert.NotNil(t, err)

		assert.IsType(t, &appErrors.BadRequestError{}, err)

		assert.Equal(t, fmt.Sprintf("%q is not a valid id", "id"), err.Error())

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("UpdateList() should call repository.RemoveList when the id is valid", func(t *testing.T) {
		l := models.List{
			ID:   "1",
			Name: "list",
		}

		mockedRepository.On("IsValidID", l.ID).Return(true).Once()
		mockedRepository.On("Update", l.ID, &l).Return(errors.New("error")).Once()

		err := service.UpdateList(l.ID, &l)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("UpdateList() should return a badRequestError when the id is not valid", func(t *testing.T) {
		id := "wadus"

		l := models.List{
			ID:   id,
			Name: "list",
		}

		mockedRepository.On("IsValidID", l.ID).Return(false).Once()

		err := service.UpdateList(id, &l)

		assert.NotNil(t, err)

		assert.IsType(t, &appErrors.BadRequestError{}, err)

		assert.Equal(t, fmt.Sprintf("%q is not a valid id", id), err.Error())

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("GetSingleList() should call repository.GetOne", func(t *testing.T) {
		l := models.List{}

		mockedRepository.On("GetOne", &l, bson.D{{"_id", "id"}}, nil).Return(errors.New("error")).Once()
		mockedRepository.On("IsValidID", "id").Return(true).Once()

		err := service.GetSingleList("id", &l)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("GetSingleList() should return a badRequestError when the id is not valid", func(t *testing.T) {
		id := "wadus"

		mockedRepository.On("IsValidID", id).Return(false).Once()

		err := service.GetSingleList(id, &models.List{})

		assert.NotNil(t, err)

		assert.IsType(t, &appErrors.BadRequestError{}, err)

		assert.Equal(t, fmt.Sprintf("%q is not a valid id", id), err.Error())

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("GetLists() should call repository.RemoveList", func(t *testing.T) {
		r := []models.GetListsResultDto{}

		mockedRepository.On("Get", &r, nil, bson.M{"name": 1}).Return(errors.New("error")).Once()

		err := service.GetLists(&r)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})
}
