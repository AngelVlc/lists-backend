package stores

import (
	"fmt"
	"github.com/AngelVlc/lists-backend/models"
)

// Store is the interface which a store must implement
type Store interface {
	GetLists() ([]models.GetListsResultDto, error)
	GetSingleList(string) (models.List, error)
	AddList(l *models.List) error
	RemoveList(string) error
	UpdateList(string, *models.List) error
	AddUser(u *models.User) error
}

// NotFoundError happens when the document does not exist in the store
type NotFoundError struct {
	ID    string
	Model string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%v with id %q not found", e.Model, e.ID)
}

// UnexpectedError is used for unexpected errors
type UnexpectedError struct {
	Msg           string
	InternalError error
}

func (e *UnexpectedError) Error() string {
	return e.Msg
}

// InvalidIDError happens when an id is not valid
type InvalidIDError struct {
	ID string
}

func (e *InvalidIDError) Error() string {
	return fmt.Sprintf("%q is not a valid id", e.ID)
}
