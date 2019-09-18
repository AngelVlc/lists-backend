package stores

import (
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockedMongoCollection struct {
	mock.Mock
}

func (m *MockedMongoCollection) FindAll() []models.List {
	args := m.Called()
	return args.Get(0).([]models.List)
}

func (m *MockedMongoCollection) Insert(l *models.List) error {
	args := m.Called(l)
	return args.Error(0)
}

type MockedMongoSession struct {
	mock.Mock
}

func (m *MockedMongoSession) Collection() MongoCollection {
	args := m.Called()
	return args.Get(0).(MongoCollection)
}

func TestStore(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)
	testMongoSession.On("Collection").Return(testMongoCollection)

	store := NewMongoStore(testMongoSession)

	t.Run("GetLists() returns all the list items", func(t *testing.T) {
		data := models.SampleListCollectionSlice()
		testMongoCollection.On("FindAll").Return(data)

		want := data
		got := store.GetLists()

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.Equal(t, want, got, "they should be equal")
	})

	t.Run("AddList() adds a new list", func(t *testing.T) {
		l := models.SampleList()

		testMongoCollection.On("Insert", &l).Return(nil)
		err := store.AddList(&l)

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.Nil(t, err)
	})

	t.Run("AddList() returns an error when the insert fails", func(t *testing.T) {
		l := models.SampleList()

		testMongoCollection.On("Insert", &l).Return(errors.New("wadus"))
		err := store.AddList(&l)

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.NotNil(t, err)

		assert.Equal(t, "Error inserting in the database", err.Error())
	})
}

func assertMocksExpectations(s *MockedMongoSession, c *MockedMongoCollection, t *testing.T) {
	s.AssertExpectations(t)
	c.AssertExpectations(t)
}
