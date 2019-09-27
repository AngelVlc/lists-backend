package stores

import (
	"github.com/AngelVlc/lists-backend/models"
	"gopkg.in/mgo.v2/bson"
)

var listsCollectionName = "lists"

func (s *MongoRepository) listsCollection() MongoCollection {
	return s.mongoSession.Collection(listsCollectionName)
}

var usersCollectionName = "users"

func (s *MongoRepository) usersCollection() MongoCollection {
	return s.mongoSession.Collection(usersCollectionName)
}

// MongoRepository is the store which uses mongo db
type MongoRepository struct {
	mongoSession MongoSession
}

// NewMongoRepository returns a new MongoRepository
func NewMongoRepository(mongoSession MongoSession) MongoRepository {
	return MongoRepository{mongoSession}
}

// GetLists returns the lists collection
func (s *MongoRepository) GetLists() ([]models.GetListsResultDto, error) {
	r := []models.GetListsResultDto{}

	if err := s.listsCollection().Find(&r, nil, bson.M{"name": 1}); err != nil {
		return []models.GetListsResultDto{}, &UnexpectedError{
			Msg:           "Error retrieving from the database",
			InternalError: err,
		}
	}

	return r, nil
}

// GetSingleList returns one list
func (s *MongoRepository) GetSingleList(id string) (models.List, error) {
	r := models.List{}

	if err := s.getSingle(s.listsCollection(), id, &r); err != nil {
		return models.List{}, err
	}

	return r, nil
}

// AddList adds a new list to the collection
func (s *MongoRepository) AddList(l *models.List) error {
	l.ID = bson.NewObjectId().Hex()

	return s.add(s.listsCollection(), l)
}

// RemoveList removes a list from the collection
func (s *MongoRepository) RemoveList(id string) error {
	return s.remove(s.listsCollection(), id)
}

// UpdateList updates a list
func (s *MongoRepository) UpdateList(id string, l *models.List) error {
	l.ID = id

	return s.update(s.listsCollection(), id, l)
}

// AddUser adds a new user
func (s *MongoRepository) AddUser(u *models.User) error {
	u.ID = bson.NewObjectId().Hex()

	return s.add(s.usersCollection(), u)
}

func (s *MongoRepository) add(c MongoCollection, doc interface{}) error {
	if err := c.Insert(doc); err != nil {
		return &UnexpectedError{
			Msg:           "Error inserting in the database",
			InternalError: err,
		}
	}

	return nil
}

func (s *MongoRepository) update(c MongoCollection, id string, doc interface{}) error {
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

func (s *MongoRepository) remove(c MongoCollection, id string) error {
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

func (s *MongoRepository) getSingle(c MongoCollection, id string, doc interface{}) error {
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

func (s *MongoRepository) isValidID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return &InvalidIDError{id}
	}

	return nil
}
