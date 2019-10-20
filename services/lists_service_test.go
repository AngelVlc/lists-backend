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

	t.Run("AddUserList() should call repository.AddList", func(t *testing.T) {
		u := "userId"
		l := models.List{
			ID:     "1",
			Name:   "list",
			UserID: u,
		}

		mockedRepository.On("Add", &l).Return("", errors.New("error")).Once()

		id, err := service.AddUserList(u, &models.List{ID: l.ID, Name: l.Name})

		assert.Empty(t, id)
		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("RemoveUserList() should call repository.Remove", func(t *testing.T) {
		mockedRepository.On("Remove", bson.D{{"_id", "id"}, {"userId", "uid"}}).Return(errors.New("error")).Once()
		mockedRepository.On("IsValidID", "id").Return(true).Once()

		err := service.RemoveUserList("id", "uid")

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("RemoveUserList() should return a badRequestError when the id is not valid", func(t *testing.T) {
		mockedRepository.On("IsValidID", "id").Return(false).Once()

		err := service.RemoveUserList("id", "uid")

		assert.NotNil(t, err)

		assert.IsType(t, &appErrors.BadRequestError{}, err)

		assert.Equal(t, fmt.Sprintf("%q is not a valid id", "id"), err.Error())

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("UpdateUserList() should call repository.Update when the id is valid", func(t *testing.T) {
		l := models.List{
			ID:   "1",
			Name: "list",
		}

		u := "userId"

		mockedRepository.On("IsValidID", l.ID).Return(true).Once()
		mockedRepository.On("Update", bson.D{{"_id", l.ID}, {"userId", u}}, &l).Return(errors.New("error")).Once()

		err := service.UpdateUserList(l.ID, u, &l)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("UpdateUserList() should return a badRequestError when the id is not valid", func(t *testing.T) {
		id := "wadus"

		l := models.List{
			ID:   id,
			Name: "list",
		}

		u := "userId"

		mockedRepository.On("IsValidID", id).Return(false).Once()

		err := service.UpdateUserList(id, u, &l)

		assert.NotNil(t, err)

		assert.IsType(t, &appErrors.BadRequestError{}, err)

		assert.Equal(t, fmt.Sprintf("%q is not a valid id", id), err.Error())

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("GetSingleUserList() should call repository.GetOne", func(t *testing.T) {
		l := models.List{}
		i := "listId"
		u := "userId"

		mockedRepository.On("GetOne", &l, bson.D{{"_id", i}, {"userId", u}}, nil).Return(errors.New("error")).Once()
		mockedRepository.On("IsValidID", i).Return(true).Once()

		err := service.GetSingleUserList(i, u, &l)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("GetSingleUserList() should return a badRequestError when the id is not valid", func(t *testing.T) {
		id := "wadus"

		mockedRepository.On("IsValidID", id).Return(false).Once()

		err := service.GetSingleUserList(id, "", &models.List{})

		assert.NotNil(t, err)

		assert.IsType(t, &appErrors.BadRequestError{}, err)

		assert.Equal(t, fmt.Sprintf("%q is not a valid id", id), err.Error())

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})

	t.Run("GetUserLists() should call repository.Get", func(t *testing.T) {
		r := []models.GetListsResultDto{}
		u := "userId"

		mockedRepository.On("Get", &r, bson.D{{"userId", u}}, bson.M{"name": 1}).Return(errors.New("error")).Once()

		err := service.GetUserLists(u, &r)

		assert.NotNil(t, err)

		mockedSession.AssertExpectations(t)
		mockedRepository.AssertExpectations(t)
	})
}
