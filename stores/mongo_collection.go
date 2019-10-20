package stores

import (
	"gopkg.in/mgo.v2"
)

// MongoCollection is an interface which contains the methods used by the mongo collection
// for testing purposes
type MongoCollection interface {
	Find(doc interface{}, query interface{}, selector interface{}) error
	FindOne(doc interface{}, query interface{}, selector interface{}) error
	Insert(doc interface{}) error
	Remove(query interface{}) error
	Update(query interface{}, doc interface{}) error
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

// Find returns all documents
func (c *MyMongoCollection) Find(doc interface{}, query interface{}, selector interface{}) error {
	return c.collection.Find(query).Select(selector).All(doc)
}

// FindOne returns a single document
func (c *MyMongoCollection) FindOne(doc interface{}, query interface{}, selector interface{}) error {
	return c.collection.Find(query).Select(selector).One(doc)
}

// Insert adds a new document
func (c *MyMongoCollection) Insert(doc interface{}) error {
	return c.collection.Insert(doc)
}

// Remove removes a document
func (c *MyMongoCollection) Remove(query interface{}) error {
	return c.collection.Remove(query)
}

// Update updates a list
func (c *MyMongoCollection) Update(query interface{}, doc interface{}) error {
	return c.collection.Update(query, doc)
}

// Name returns the mongo collection name
func (c *MyMongoCollection) Name() string {
	return c.collection.Name
}
