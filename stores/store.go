// Package stores contains the stores used in the app
package stores

import (
	"github.com/AngelVlc/lists-backend/models"
)

// Store is the interface which a store must implement
type Store interface {
	GetLists() []models.List
	AddList(*models.List) error
	RemoveList(string) error
	UpdateList(string, *models.List) error
}
