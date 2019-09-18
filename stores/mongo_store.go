package stores

import (
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"gopkg.in/mgo.v2/bson"
	"log"
)

// MongoStore is the store which uses mongo db
type MongoStore struct {
	mongoSession MongoSession
}

// NewMongoStore returns a new MongoStore
func NewMongoStore(mongoSession MongoSession) MongoStore {
	return MongoStore{mongoSession}
}

// GetLists returns the lists collection
func (s *MongoStore) GetLists() []models.List {
	return s.mongoSession.Collection().FindAll()
}

// AddList adds a new list to the collection and returns it
func (s *MongoStore) AddList(l *models.List) error {
	l.ID = bson.NewObjectId()
	err := s.mongoSession.Collection().Insert(l)
	if err != nil {
		log.Println("Error inserting. Error: " + err.Error())
		return errors.New("Error inserting in the database")
	}

	return nil
}
