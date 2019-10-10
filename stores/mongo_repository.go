package stores

import (
	"fmt"
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"gopkg.in/mgo.v2/bson"
	"reflect"
)

// MongoRepository is the store which uses mongo db
type MongoRepository struct {
	mongoCollection MongoCollection
}

// Get returns the lists collection
func (s *MongoRepository) Get(doc interface{}, query interface{}, selector interface{}) error {
	if err := s.mongoCollection.Find(doc, query, selector); err != nil {
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
func (s *MongoRepository) Update(id string, doc interface{}) error {
	if err := s.isValidID(id); err != nil {
		return err
	}

	reflect.ValueOf(doc).Elem().FieldByName("ID").SetString(id)

	if err := s.mongoCollection.Update(id, doc); err != nil {
		if err.Error() == "not found" {
			return &appErrors.NotFoundError{
				ID:    id,
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
func (s *MongoRepository) Remove(id string) error {
	if err := s.isValidID(id); err != nil {
		return err
	}

	if err := s.mongoCollection.Remove(id); err != nil {
		if err.Error() == "not found" {
			return &appErrors.NotFoundError{
				ID:    id,
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

// GetSingle returns a single document
func (s *MongoRepository) GetSingle(id string, doc interface{}) error {
	if err := s.isValidID(id); err != nil {
		return err
	}

	if err := s.mongoCollection.FindOne(id, doc); err != nil {
		if err.Error() == "not found" {
			return &appErrors.NotFoundError{
				ID:    id,
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

func (s *MongoRepository) isValidID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return &appErrors.BadRequestError{Msg: fmt.Sprintf("%q is not a valid id", id), InternalError: nil}
	}

	return nil
}
