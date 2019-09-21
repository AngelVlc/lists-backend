package stores

import (
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var listsCollectionName = "lists"

func (s *MongoStore) listsCollection() MongoCollection {
	return s.mongoSession.Collection(listsCollectionName)
}

var usersCollectionName = "users"

func (s *MongoStore) usersCollection() MongoCollection {
	return s.mongoSession.Collection(usersCollectionName)
}


// MongoStore is the store which uses mongo db
type MongoStore struct {
	mongoSession MongoSession
}

// NewMongoStore returns a new MongoStore
func NewMongoStore(mongoSession MongoSession) MongoStore {
	return MongoStore{mongoSession}
}

// GetLists returns the lists collection
func (s *MongoStore) GetLists() []models.GetListsResultDto {
	return s.listsCollection().FindAll()
}

// GetSingleList returns one list
func (s *MongoStore) GetSingleList(id string) (models.List, error) {
	oid := bson.ObjectIdHex(id)

	return s.listsCollection().FindOne(oid)
}

// AddList adds a new list to the collection
func (s *MongoStore) AddList(l *models.List) error {
	l.ID = bson.NewObjectId()
	err := s.listsCollection().Insert(l)
	if err != nil {
		log.Println("Error inserting. Error: " + err.Error())
		return errors.New("Error inserting in the database")
	}

	return nil
}

// RemoveList removes a list from the collection
func (s *MongoStore) RemoveList(id string) error {
	oid := bson.ObjectIdHex(id)

	if err := s.listsCollection().Remove(oid); err != nil {
		log.Println("Error removing. Error: " + err.Error())
		return errors.New("Error removing from the database")
	}

	return nil
}

// UpdateList updates a list
func (s *MongoStore) UpdateList(id string, l *models.List) error {
	oid := bson.ObjectIdHex(id)

	l.ID = oid

	if err := s.listsCollection().Update(oid, l); err != nil {
		log.Println("Error updating. Error: " + err.Error())
		return errors.New("Error updating the database")
	}

	return nil
}

// AddUser adds a new user
func (s *MongoStore) AddUser(u *models.User) error {
	u.ID = bson.NewObjectId()
	err := s.usersCollection().Insert(u)
	if err != nil {
		log.Println("Error inserting. Error: " + err.Error())
		return errors.New("Error inserting in the database")
	}

	return nil
}