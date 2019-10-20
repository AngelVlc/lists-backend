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
	AddUserList(userID string, l *models.List) (string, error)
	RemoveUserList(id string, userID string) error
	UpdateUserList(id string, userID string, l *models.List) error
	GetSingleUserList(id string, userID string, l *models.List) error
	GetUserLists(userID string, r *[]models.GetListsResultDto) error
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

// AddUserList  adds a user
func (s *MyListsService) AddUserList(userID string, l *models.List) (string, error) {
	l.UserID = userID
	return s.listsRepository().Add(l)
}

// RemoveUserList removes a list
func (s *MyListsService) RemoveUserList(id string, userID string) error {
	if !s.listsRepository().IsValidID(id) {
		return s.getInvalidIDError(id)
	}

	return s.listsRepository().Remove(bson.D{{"_id", id}, {"userId", userID}})
}

// UpdateUserList updates an existing list
func (s *MyListsService) UpdateUserList(id string, userID string, l *models.List) error {
	if !s.listsRepository().IsValidID(id) {
		return s.getInvalidIDError(id)
	}

	l.ID = id
	l.UserID = userID

	return s.listsRepository().Update(bson.D{{"_id", id}, {"userId", userID}}, l)
}

// GetSingleUserList returns a single list from its id
func (s *MyListsService) GetSingleUserList(id string, userID string, l *models.List) error {
	if !s.listsRepository().IsValidID(id) {
		return s.getInvalidIDError(id)
	}

	return s.listsRepository().GetOne(l, bson.D{{"_id", id}, {"userId", userID}}, nil)
}

// GetUserLists returns the lists for the given user
func (s *MyListsService) GetUserLists(userID string, r *[]models.GetListsResultDto) error {
	return s.listsRepository().Get(r, bson.D{{"userId", userID}}, bson.M{"name": 1})
}

func (s *MyListsService) listsRepository() stores.Repository {
	return s.session.GetRepository("lists")
}

func (s *MyListsService) getInvalidIDError(id string) error {
	return &appErrors.BadRequestError{Msg: fmt.Sprintf("%q is not a valid id", id), InternalError: nil}
}
