package services

import (
	"fmt"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"gopkg.in/mgo.v2/bson"
)

// ListsService is the interface a lists service must implement
type ListsService interface {
	AddList(l *models.List) (string, error)
	RemoveList(id string) error
	UpdateList(id string, l *models.List) error
	GetSingleList(id string, l *models.List) error
	GetLists(r *[]models.GetListsResultDto) error
}

// MyListsService is the service for the list entity
type MyListsService struct {
	session stores.MongoSession
}

// NewMyListsService returns a new lists service
func NewMyListsService(session stores.MongoSession) *MyListsService {
	return &MyListsService{
		session: session,
	}
}

// AddList  adds a user
func (s *MyListsService) AddList(l *models.List) (string, error) {
	return s.listsRepository().Add(l)
}

// RemoveList removes a list
func (s *MyListsService) RemoveList(id string) error {
	if !s.listsRepository().IsValidID(id) {
		return s.getInvalidIDError(id)
	}

	return s.listsRepository().Remove(id)
}

// UpdateList updates an existing list
func (s *MyListsService) UpdateList(id string, l *models.List) error {
	if !s.listsRepository().IsValidID(id) {
		return s.getInvalidIDError(id)
	}

	return s.listsRepository().Update(id, l)
}

// GetSingleList returns a single list from its id
func (s *MyListsService) GetSingleList(id string, l *models.List) error {
	if !s.listsRepository().IsValidID(id) {
		return s.getInvalidIDError(id)
	}

	return s.listsRepository().GetSingle(id, l)
}

// GetLists returns the lists
func (s *MyListsService) GetLists(r *[]models.GetListsResultDto) error {
	return s.listsRepository().Get(r, nil, bson.M{"name": 1})
}

func (s *MyListsService) listsRepository() stores.Repository {
	return s.session.GetRepository("lists")
}

func (s *MyListsService) getInvalidIDError(id string) error {
	return &appErrors.BadRequestError{Msg: fmt.Sprintf("%q is not a valid id", id), InternalError: nil}
}
