// Package stores contains the stores used in the app
package stores

import (
	"github.com/AngelVlc/lists-backend/models"
)

// Store is the interface which a store must implement
type Store interface {
	GetLists() []models.GetListsResultDto
	GetSingleList(string) (models.List, error)
	AddList(l *models.List) error
	RemoveList(string) error
	UpdateList(string, *models.List) error
}
