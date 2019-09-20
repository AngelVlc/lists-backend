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

func (m *MockedMongoCollection) FindOne(id bson.ObjectId) (models.List, error) {
	args := m.Called(id)
	return args.Get(0).(models.List), args.Error(1)
}

func (m *MockedMongoCollection) Insert(l *models.List) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *MockedMongoCollection) Remove(id bson.ObjectId) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockedMongoCollection) Update(id bson.ObjectId, l *models.List) error {
	args := m.Called(id, l)
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

func TestStore(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)
	testMongoSession.On("Collection", ListsCollectionName).Return(testMongoCollection)

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
		got, err := store.GetSingleList(data.ID.Hex())

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.Equal(t, want, got, "they should be equal")

		assert.Nil(t, err)
	})

	t.Run("GetSingleList() returns an error when the id is not a valid bson object id", func(t *testing.T) {
		_, err := store.GetSingleList("wadus")

		assert.NotNil(t, err)

		assert.Equal(t, "Error getting the list from the database, invalid id", err.Error())
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

	t.Run("RemoveList() returns an error when the remove fails", func(t *testing.T) {
		oidHex := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", bson.ObjectIdHex(oidHex)).Return(errors.New("wadus"))

		err := store.RemoveList(oidHex)

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.NotNil(t, err)

		assert.Equal(t, "Error removing from the database", err.Error())
	})

	t.Run("RemoveList() returns an error when the id is not a valid bson object id", func(t *testing.T) {
		err := store.RemoveList("wadus")

		assert.NotNil(t, err)

		assert.Equal(t, "Error removing from the database, invalid id", err.Error())
	})

	t.Run("RemoveList() removes a list", func(t *testing.T) {
		oidHex := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", bson.ObjectIdHex(oidHex)).Return(nil)

		err := store.RemoveList(oidHex)

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.Nil(t, err)
	})

	t.Run("UpdateList() returns an error when the update fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testMongoCollection.On("Update", bson.ObjectIdHex(id), mock.Anything).Return(errors.New("wadus"))

		l := models.SampleList()
		err := store.UpdateList(id, &l)

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.NotNil(t, err)

		assert.Equal(t, "Error updating the database", err.Error())
	})

	t.Run("UpdateList() returns an error when the id is not a valid bson object id", func(t *testing.T) {
		l := models.SampleList()
		err := store.UpdateList("wadus", &l)

		assert.NotNil(t, err)

		assert.Equal(t, "Error updating the database, invalid id", err.Error())
	})

	t.Run("UpdateList() updates a list", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testMongoCollection.On("Update", bson.ObjectIdHex(id), mock.Anything).Return(nil)

		l := models.SampleList()
		err := store.UpdateList(id, &l)

		assertMocksExpectations(testMongoSession, testMongoCollection, t)

		assert.Nil(t, err)
	})

}

func assertMocksExpectations(s *MockedMongoSession, c *MockedMongoCollection, t *testing.T) {
	s.AssertExpectations(t)
	c.AssertExpectations(t)
}
