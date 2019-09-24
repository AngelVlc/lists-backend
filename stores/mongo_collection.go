package stores

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoCollection is an interface which contains the methods used by the mongo collection
// for testing purposes
type MongoCollection interface {
	FindAll(doc interface{}) error
	FindOne(id string, doc interface{}) error
	Insert(doc interface{}) error
	Remove(id string) error
	Update(id string, doc interface{}) error
	DropCollection() error
	Name() string
}

// MyMongoCollection implements the MongoCollection interface
type MyMongoCollection struct {
	collection *mgo.Collection
}

// NewMyMongoCollection returns a new MyMongoCollection
func NewMyMongoCollection(c *mgo.Collection) *MyMongoCollection {
	return &MyMongoCollection{c}
}

// FindAll returns all documents
func (c *MyMongoCollection) FindAll(doc interface{}) error {
	return c.collection.Find(nil).Select(bson.M{"name": 1}).All(doc)
}

// FindOne returns a single document
func (c *MyMongoCollection) FindOne(id string, doc interface{}) error {
	if err := c.collection.FindId(id).One(doc); err != nil {
		return err
	}

	return nil
}

// Insert adds a new document
func (c *MyMongoCollection) Insert(doc interface{}) error {
	return c.collection.Insert(doc)
}

// Remove removes a document
func (c *MyMongoCollection) Remove(id string) error {
	return c.collection.RemoveId(id)
}

// Update updates a list
func (c *MyMongoCollection) Update(id string, doc interface{}) error {
	return c.collection.UpdateId(id, doc)
}

// DropCollection drops the collection
func (c *MyMongoCollection) DropCollection() error {
	return c.collection.DropCollection()
}

// Name returns the mongo collection name
func (c *MyMongoCollection) Name() string {
	return c.collection.Name
}
