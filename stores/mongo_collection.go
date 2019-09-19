package stores

import (
	"github.com/AngelVlc/lists-backend/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoCollection is an interface which contains the methods used by the mongo collection
type MongoCollection interface {
	FindAll() []models.List
	FindOne(bson.ObjectId) (models.List, error)
	Insert(*models.List) error
	Remove(bson.ObjectId) error
	Update(bson.ObjectId, *models.List) error
	DropCollection() error
}

// MyMongoCollection contains the methods used by the mongo collection
type MyMongoCollection struct {
	collection *mgo.Collection
}

// NewMyMongoCollection returns a new MyMongoCollection
func NewMyMongoCollection(c *mgo.Collection) *MyMongoCollection {
	return &MyMongoCollection{c}
}

// FindAll returns every list
func (c *MyMongoCollection) FindAll() []models.List {
	r := []models.List{}
	c.collection.Find(nil).All(&r)
	return r
}

// FindOne returns one list
func (c *MyMongoCollection) FindOne(id bson.ObjectId) (models.List, error) {
	r := models.List{}
	err := c.collection.FindId(id).One(&r)
	return r, err
}

// Insert adds a new list
func (c *MyMongoCollection) Insert(l *models.List) error {
	return c.collection.Insert(l)
}

// Remove removes a list
func (c *MyMongoCollection) Remove(id bson.ObjectId) error {
	return c.collection.RemoveId(id)
}

// Update updates a list
func (c *MyMongoCollection) Update(id bson.ObjectId, l *models.List) error {
	return c.collection.UpdateId(id, l)
}

// Drop collection drops the collection
func (c *MyMongoCollection) DropCollection() error {
	return c.collection.DropCollection()
}
