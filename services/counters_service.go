package services

import (
	"github.com/AngelVlc/lists-backend/models"
	"github.com/AngelVlc/lists-backend/stores"
	"gopkg.in/mgo.v2/bson"
)

// CountersService contains the methods for working with counters
type CountersService interface {
	AddCounter(name string) error
	IncrementCounter(name string) error
	ExistsCounter(name string) bool
	GetCounterValue(name string) (int, error)
}

// MyCountersService is the service for working with counters
type MyCountersService struct {
	session stores.MongoSession
}

// NewMyCountersService creates a MyCountersService
func NewMyCountersService(session stores.MongoSession) *MyCountersService {
	return &MyCountersService{
		session: session,
	}
}

// AddCounter adds a new counter
func (s *MyCountersService) AddCounter(name string) error {
	c := models.Counter{
		Name:  name,
		Value: 1,
	}
	_, err := s.countersRepository().Add(&c)
	return err
}

// ExistsCounter returns true if the counter already exist
func (s *MyCountersService) ExistsCounter(name string) bool {
	err := s.countersRepository().GetOne(&models.Counter{}, bson.D{{"name", name}}, nil)

	if err == nil {
		return true
	}

	return false
}

// GetCounterValue returns a counter's value
func (s *MyCountersService) GetCounterValue(name string) (int, error) {
	c := models.Counter{}
	err := s.countersRepository().GetOne(&c, bson.D{{"name", name}}, nil)

	if err == nil {
		return c.Value, nil
	}

	return -1, err
}

// IncrementCounter increments a counter
func (s *MyCountersService) IncrementCounter(name string) error {
	return s.countersRepository().Update(bson.D{{"name", name}}, bson.M{"$inc": bson.M{"value": 1}})
}

func (s *MyCountersService) countersRepository() stores.Repository {
	return s.session.GetRepository("counters")
}
