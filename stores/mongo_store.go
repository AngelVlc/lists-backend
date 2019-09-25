package stores

import (
	"github.com/AngelVlc/lists-backend/models"
	"gopkg.in/mgo.v2/bson"
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
func (s *MongoStore) GetLists() ([]models.GetListsResultDto, error) {
	r := []models.GetListsResultDto{}

	if err := s.listsCollection().FindAll(&r); err != nil {
		return []models.GetListsResultDto{}, &UnexpectedError{
			Msg:           "Error retrieving from the database",
			InternalError: err,
		}
	}

	return r, nil
}

// GetSingleList returns one list
func (s *MongoStore) GetSingleList(id string) (models.List, error) {
	r := models.List{}

	if err := s.getSingle(s.listsCollection(), id, &r); err != nil {
		return models.List{}, err
	}

	return r, nil
}

// AddList adds a new list to the collection
func (s *MongoStore) AddList(l *models.List) error {
	l.ID = bson.NewObjectId().Hex()

	return s.add(s.listsCollection(), l)
}

// RemoveList removes a list from the collection
func (s *MongoStore) RemoveList(id string) error {
	return s.remove(s.listsCollection(), id)
}

// UpdateList updates a list
func (s *MongoStore) UpdateList(id string, l *models.List) error {
	l.ID = id

	return s.update(s.listsCollection(), id, l)
}

// AddUser adds a new user
func (s *MongoStore) AddUser(u *models.User) error {
	u.ID = bson.NewObjectId().Hex()

	return s.add(s.usersCollection(), u)
}

func (s *MongoStore) add(c MongoCollection, doc interface{}) error {
	if err := c.Insert(doc); err != nil {
		return &UnexpectedError{
			Msg:           "Error inserting in the database",
			InternalError: err,
		}
	}

	return nil
}

func (s *MongoStore) update(c MongoCollection, id string, doc interface{}) error {
	if err := s.isValidID(id); err != nil {
		return err
	}

	if err := c.Update(id, doc); err != nil {
		if err.Error() == "not found" {
			return &NotFoundError{
				ID:    id,
				Model: c.Name(),
			}
		}
		return &UnexpectedError{
			Msg:           "Error updating the database",
			InternalError: err,
		}
	}

	return nil
}

func (s *MongoStore) remove(c MongoCollection, id string) error {
	if err := s.isValidID(id); err != nil {
		return err
	}

	if err := c.Remove(id); err != nil {
		if err.Error() == "not found" {
			return &NotFoundError{
				ID:    id,
				Model: c.Name(),
			}
		}
		return &UnexpectedError{
			Msg:           "Error removing from the database",
			InternalError: err,
		}
	}

	return nil
}

func (s *MongoStore) getSingle(c MongoCollection, id string, doc interface{}) error {
	if err := s.isValidID(id); err != nil {
		return err
	}

	if err := c.FindOne(id, doc); err != nil {
		if err.Error() == "not found" {
			return &NotFoundError{
				ID:    id,
				Model: c.Name(),
			}
		}
		return &UnexpectedError{
			Msg:           "Error retrieving from the database",
			InternalError: err,
		}
	}

	return nil
}

func (s *MongoStore) isValidID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return &InvalidIDError{id}
	}

	return nil
}
