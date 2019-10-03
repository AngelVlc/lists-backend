package stores

import (
	"gopkg.in/mgo.v2"
	"log"
	"os"
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
func NewMyMongoSession(useTestConfig bool) *MyMongoSession {
	url := os.Getenv("MONGO_URL")
	var databaseName string
	if !useTestConfig {
		databaseName = os.Getenv("MONGO_DATABASE_NAME")
	} else {
		databaseName = os.Getenv("MONGO_TEST_DATABASE_NAME")
	}

	s, err := mgo.Dial(url)
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
