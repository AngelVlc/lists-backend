package stores

import (
	"reflect"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"gopkg.in/mgo.v2/bson"
)

// MongoRepository is the store which uses mongo db
type MongoRepository struct {
	mongoCollection MongoCollection
}

// Get returns several items from a collection
func (s *MongoRepository) Get(doc interface{}, query interface{}, selector interface{}) error {
	if err := s.mongoCollection.Find(doc, query, selector); err != nil {
		return &appErrors.UnexpectedError{
			Msg:           "Error retrieving from the database",
			InternalError: err,
		}
	}

	return nil
}

// GetOne returns a single item
func (s *MongoRepository) GetOne(doc interface{}, query interface{}, selector interface{}) error {
	if err := s.mongoCollection.FindOne(doc, query, selector); err != nil {
		if err.Error() == "not found" {
			return &appErrors.NotFoundError{
				Model: s.mongoCollection.Name(),
			}
		}
		return &appErrors.UnexpectedError{
			Msg:           "Error retrieving from the database",
			InternalError: err,
		}
	}

	return nil
}

// Add adds a new document to the collection
func (s *MongoRepository) Add(doc interface{}) (string, error) {
	id := bson.NewObjectId().Hex()
	reflect.ValueOf(doc).Elem().FieldByName("ID").SetString(id)

	if err := s.mongoCollection.Insert(doc); err != nil {
		return "", &appErrors.UnexpectedError{
			Msg:           "Error inserting in the database",
			InternalError: err,
		}
	}

	return id, nil
}

// Update updates a document
func (s *MongoRepository) Update(query interface{}, doc interface{}) error {
	if err := s.mongoCollection.Update(query, doc); err != nil {
		if err.Error() == "not found" {
			return &appErrors.NotFoundError{
				Model: s.mongoCollection.Name(),
			}
		}
		return &appErrors.UnexpectedError{
			Msg:           "Error updating the database",
			InternalError: err,
		}
	}

	return nil
}

// Remove removes a document from the collection
func (s *MongoRepository) Remove(query interface{}) error {
	if err := s.mongoCollection.Remove(query); err != nil {
		if err.Error() == "not found" {
			return &appErrors.NotFoundError{
				Model: s.mongoCollection.Name(),
			}
		}
		return &appErrors.UnexpectedError{
			Msg:           "Error removing from the database",
			InternalError: err,
		}
	}

	return nil
}

// IsValidID returns true if the id is valid
func (s *MongoRepository) IsValidID(id string) bool {
	return bson.IsObjectIdHex(id)
}
