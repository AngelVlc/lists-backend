package stores

import (
	"gopkg.in/mgo.v2"
)

// MongoSession is the interface used to retrieve the mongo collection
type MongoSession interface {
	Collection() MongoCollection
}

// MyMongoSession is the object used to access the mongo collection
type MyMongoSession struct {
	session *mgo.Session
}

// NewMyMongoSession returns a new MyMongoSession
func NewMyMongoSession(s *mgo.Session) *MyMongoSession {
	return &MyMongoSession{s}
}

// Collection returns the mongo collection
func (s *MyMongoSession) Collection() MongoCollection {
	c := s.session.DB("lists").C("lists")
	return NewMyMongoCollection(c)
}
