package stores

import (
	"gopkg.in/mgo.v2"
	"log"
	"os"
)

// MongoSession is the interface used to retrieve the mongo collection
type MongoSession interface {
	Collection(name string) MongoCollection
}

// MyMongoSession is the object used to access the mongo collection
type MyMongoSession struct {
	session        *mgo.Session
	databaseName   string
	collectionName string
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

// Collection returns the mongo collection
func (s *MyMongoSession) Collection(name string) MongoCollection {
	c := s.session.DB(s.databaseName).C(name)
	return NewMyMongoCollection(c)
}
