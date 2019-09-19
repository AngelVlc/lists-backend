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
	databaseName string
	collectionName string
}

// NewMyMongoSession returns a new MyMongoSession
func NewMyMongoSession(s *mgo.Session, databaseName string, collectionName string) *MyMongoSession {
	return &MyMongoSession{
		session: s,
		databaseName: databaseName,
		collectionName: collectionName,
	}
}

// Collection returns the mongo collection
func (s *MyMongoSession) Collection() MongoCollection {
	c := s.session.DB(s.databaseName).C(s.collectionName)
	return NewMyMongoCollection(c)
}
