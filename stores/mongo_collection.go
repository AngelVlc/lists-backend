package stores

import (
	"github.com/AngelVlc/lists-backend/models"
	"gopkg.in/mgo.v2"
)

// MongoCollection is an interface which contains the methods used by the mongo collection
type MongoCollection interface {
	FindAll() []models.List
	Insert(l *models.List) error
}

// MyMongoCollection contains the methods used by the mongo collection
type MyMongoCollection struct {
	collection *mgo.Collection
}

// NewMyMongoCollection returns a new MyMongoCollection
func NewMyMongoCollection(c *mgo.Collection) *MyMongoCollection {
	return &MyMongoCollection{c}
}

// FindAll queries the mongo collection and returns every document
func (c *MyMongoCollection) FindAll() []models.List {
	r := []models.List{}
	c.collection.Find(nil).All(&r)
	return r
}

// Insert adds a new document to the collection
func (c *MyMongoCollection) Insert(l *models.List) error {
	return c.collection.Insert(l)
}
