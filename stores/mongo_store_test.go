package stores

import (
	"errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type MockedMongoCollection struct {
	mock.Mock
}

func (m *MockedMongoCollection) FindAll() []models.GetListsResultDto {
	args := m.Called()
	return args.Get(0).([]models.GetListsResultDto)
}

func (m *MockedMongoCollection) FindOne(id string) (models.List, error) {
	args := m.Called(id)
	return args.Get(0).(models.List), args.Error(1)
}

func (m *MockedMongoCollection) Insert(doc interface{}) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *MockedMongoCollection) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockedMongoCollection) Update(id string, doc interface{}) error {
	args := m.Called(id, doc)
	return args.Error(0)
}

func (m *MockedMongoCollection) DropCollection() error {
	args := m.Called()
	return args.Error(0)
}

type MockedMongoSession struct {
	mock.Mock
}

func (m *MockedMongoSession) Collection(name string) MongoCollection {
	args := m.Called(name)
	return args.Get(0).(MongoCollection)
}

func TestStoreForLists(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)
	testMongoSession.On("Collection", listsCollectionName).Return(testMongoCollection)

	store := NewMongoStore(testMongoSession)

	t.Run("GetLists() returns all the list items", func(t *testing.T) {
		data := models.SampleGetListsResultDto()
		testMongoCollection.On("FindAll").Return(data)

		want := data
		got := store.GetLists()

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.Equal(t, want, got, "they should be equal")
	})

	t.Run("GetSingleList() returns a single list", func(t *testing.T) {
		data := models.SampleListSlice()[0]
		testMongoCollection.On("FindOne", data.ID).Return(data, nil)

		want := data
		got, err := store.GetSingleList(data.ID)

		assert.Equal(t, want, got, "they should be equal")

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("AddList() adds a new list", func(t *testing.T) {
		l := models.SampleList()

		testMongoCollection.On("Insert", &l).Return(nil)
		err := store.AddList(&l)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("AddList() returns an error when the insert fails", func(t *testing.T) {
		l := models.SampleList()

		testMongoCollection.On("Insert", &l).Return(errors.New("wadus"))
		err := store.AddList(&l)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error inserting in the database")
	})

	t.Run("RemoveList() returns an error when the remove fails", func(t *testing.T) {
		oidHex := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", oidHex).Return(errors.New("wadus"))

		err := store.RemoveList(oidHex)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error removing from the database")
	})

	t.Run("RemoveList() removes a list", func(t *testing.T) {
		oidHex := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", oidHex).Return(nil)

		err := store.RemoveList(oidHex)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("UpdateList() returns an error when the update fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testMongoCollection.On("Update", id, mock.Anything).Return(errors.New("wadus"))

		l := models.SampleList()
		err := store.UpdateList(id, &l)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error updating the database")
	})

	t.Run("UpdateList() updates a list", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testMongoCollection.On("Update", id, mock.Anything).Return(nil)

		l := models.SampleList()
		err := store.UpdateList(id, &l)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})
}

func TestStoreForUsers(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)
	testMongoSession.On("Collection", usersCollectionName).Return(testMongoCollection)

	store := NewMongoStore(testMongoSession)

	t.Run("AddUser() adds a new list", func(t *testing.T) {
		u := models.SampleUser()

		testMongoCollection.On("Insert", &u).Return(nil)
		err := store.AddUser(&u)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("AddUser() returns an error when the insert fails", func(t *testing.T) {
		u := models.SampleUser()

		testMongoCollection.On("Insert", &u).Return(errors.New("wadus"))
		err := store.AddUser(&u)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error inserting in the database")
	})
}

func assertMocksExpectations(s *MockedMongoSession, c *MockedMongoCollection, t *testing.T) {
	s.AssertExpectations(t)
	c.AssertExpectations(t)
}

func assertSuccededOperation(t *testing.T, s *MockedMongoSession, c *MockedMongoCollection, err error) {
	assertMocksExpectations(s, c, t)

	assert.Nil(t, err)
}

func assertFailedOperation(t *testing.T, s *MockedMongoSession, c *MockedMongoCollection, err error, errorMsg string) {
	assertMocksExpectations(s, c, t)

	assert.NotNil(t, err)

	assert.Equal(t, errorMsg, err.Error())
}
