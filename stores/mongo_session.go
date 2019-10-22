package stores

import (
	"log"
	"strings"

	"gopkg.in/mgo.v2"
)

// MongoSession is the interface used to retrieve the mongo collection
type MongoSession interface {
	GetRepository(collectionName string) Repository
}

// MyMongoSession is the object used to access the mongo collection
type MyMongoSession struct {
	session      *mgo.Session
	databaseName string
}

// NewMyMongoSession returns a new MyMongoSession
func NewMyMongoSession(mongoUri string) *MyMongoSession {
	parts := strings.Split(mongoUri, "/")
	databaseName := parts[len(parts)-1]

	s, err := mgo.Dial(mongoUri)
	if err != nil {
		panic(err)
	}

	if err = s.Ping(); err != nil {
		panic(err)
	}

	log.Println("Connected with lists mongo database.")

	return &MyMongoSession{
		session:      s,
		databaseName: databaseName,
	}
}

// GetRepository returns a mongo repository for the given collection
func (s *MyMongoSession) GetRepository(collectionName string) Repository {
	ms := s.session.Copy()
	c := ms.DB(s.databaseName).C(collectionName)
	mc := NewMyMongoCollection(c)
	return &MongoRepository{mc}
}
