package stores

import (
	"errors"
	"fmt"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

type MockedMongoCollection struct {
	mock.Mock
}

func (m *MockedMongoCollection) Find(doc interface{}, query interface{}, selector interface{}) error {
	args := m.Called(doc, query, selector)
	return args.Error(0)
}

func (m *MockedMongoCollection) FindOne(id string, doc interface{}) error {
	args := m.Called(id, doc)
	return args.Error(0)
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

func (m *MockedMongoCollection) Name() string {
	args := m.Called()
	return args.String(0)
}

type MockedMongoSession struct {
	mock.Mock
}

func (m *MockedMongoSession) Collection(name string) MongoCollection {
	args := m.Called(name)
	return args.Get(0).(MongoCollection)
}

func TestUpdate(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)

	store := NewMongoStore(testMongoSession)

	t.Run("update() returns an unexpected error when the update fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testMongoCollection.On("Update", id, &id).Return(errors.New("wadus")).Once()

		err := store.update(testMongoCollection, id, &id)

		assert.IsType(t, &UnexpectedError{}, err)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error updating the database")
	})

	t.Run("update() returns a not found error when document does not exits", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testMongoCollection.On("Update", id, &id).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		err := store.update(testMongoCollection, id, &id)

		assert.IsType(t, &NotFoundError{}, err)

		msg := fmt.Sprintf("document with id %q not found", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("update() returns an invalid id error when the id is not valid", func(t *testing.T) {
		id := "wadus"

		err := store.update(testMongoCollection, id, &id)

		assert.IsType(t, &InvalidIDError{}, err)

		msg := fmt.Sprintf("%q is not a valid id", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})
}

func TestRemove(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)

	store := NewMongoStore(testMongoSession)

	t.Run("remove() returns an unexpected error when the remove fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", id).Return(errors.New("wadus")).Once()

		err := store.remove(testMongoCollection, id)

		assert.IsType(t, &UnexpectedError{}, err)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error removing from the database")
	})

	t.Run("remove() returns a not found error when document does not exits", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testMongoCollection.On("Remove", id).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		err := store.remove(testMongoCollection, id)

		assert.IsType(t, &NotFoundError{}, err)

		msg := fmt.Sprintf("document with id %q not found", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("remove() returns an invalid id error when the id is not valid", func(t *testing.T) {
		id := "wadus"

		err := store.remove(testMongoCollection, id)

		assert.IsType(t, &InvalidIDError{}, err)

		msg := fmt.Sprintf("%q is not a valid id", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})
}

func TestGetSingle(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)

	store := NewMongoStore(testMongoSession)

	t.Run("getSingle() returns an unexpected error when the remove fails", func(t *testing.T) {
		data := models.SampleListSlice()[0]
		testMongoCollection.On("FindOne", data.ID, &models.List{}).Return(errors.New("wadus")).Once()

		err := store.getSingle(testMongoCollection, data.ID, &models.List{})

		assert.IsType(t, &UnexpectedError{}, err)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error retrieving from the database")
	})

	t.Run("getSingle() returns a not found error when document does not exits", func(t *testing.T) {
		data := models.SampleListSlice()[0]
		testMongoCollection.On("FindOne", data.ID, &models.List{}).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		err := store.getSingle(testMongoCollection, data.ID, &models.List{})

		assert.IsType(t, &NotFoundError{}, err)

		msg := fmt.Sprintf("document with id %q not found", data.ID)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("getSingle() returns an invalid id error when the id is not valid", func(t *testing.T) {
		id := "wadus"

		err := store.getSingle(testMongoCollection, id, &models.List{})

		assert.IsType(t, &InvalidIDError{}, err)

		msg := fmt.Sprintf("%q is not a valid id", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})
}

func TestAdd(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)

	store := NewMongoStore(testMongoSession)

	t.Run("add() returns an unexpected error when the insert fails", func(t *testing.T) {
		l := models.SampleList()
		testMongoCollection.On("Insert", &l).Return(errors.New("wadus")).Once()

		err := store.add(testMongoCollection, &l)

		assert.IsType(t, &UnexpectedError{}, err)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error inserting in the database")
	})
}

func TestStoreForLists(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)
	testMongoSession.On("Collection", listsCollectionName).Return(testMongoCollection)

	store := NewMongoStore(testMongoSession)

	t.Run("GetLists() returns all the list items", func(t *testing.T) {
		data := models.SampleGetListsResultDto()
		testMongoCollection.On("Find", &[]models.GetListsResultDto{}, nil, bson.M{"name": 1}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]models.GetListsResultDto)
			*arg = data
		})

		want := data
		got, err := store.GetLists()

		assert.Equal(t, want, got, "they should be equal")

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("GetLists() returns an error when the query fails", func(t *testing.T) {
		testMongoCollection.On("Find", &[]models.GetListsResultDto{}, nil, bson.M{"name": 1}).Return(errors.New("wadus")).Once()

		_, err := store.GetLists()

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error retrieving from the database")
	})

	t.Run("GetSingleList() returns a single list", func(t *testing.T) {
		data := models.SampleListSlice()[0]
		testMongoCollection.On("FindOne", data.ID, &models.List{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*models.List)
			*arg = data
		})

		want := data
		got, err := store.GetSingleList(data.ID)

		assert.Equal(t, want, got, "they should be equal")

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("GetSingleList() returns an error when the query fails", func(t *testing.T) {
		data := models.SampleListSlice()[0]
		testMongoCollection.On("FindOne", data.ID, &models.List{}).Return(errors.New("wadus")).Once()

		_, err := store.GetSingleList(data.ID)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error retrieving from the database")
	})

	t.Run("AddList() adds a new list", func(t *testing.T) {
		l := models.SampleList()

		testMongoCollection.On("Insert", &l).Return(nil).Once()
		err := store.AddList(&l)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("AddList() returns an error when the insert fails", func(t *testing.T) {
		l := models.SampleList()

		testMongoCollection.On("Insert", &l).Return(errors.New("wadus")).Once()
		err := store.AddList(&l)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error inserting in the database")
	})

	t.Run("RemoveList() returns an error when the remove fails", func(t *testing.T) {
		oidHex := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", oidHex).Return(errors.New("wadus")).Once()

		err := store.RemoveList(oidHex)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error removing from the database")
	})

	t.Run("RemoveList() removes a list", func(t *testing.T) {
		oidHex := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", oidHex).Return(nil).Once()

		err := store.RemoveList(oidHex)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("UpdateList() returns an error when the update fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testMongoCollection.On("Update", id, mock.Anything).Return(errors.New("wadus")).Once()

		l := models.SampleList()
		err := store.UpdateList(id, &l)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error updating the database")
	})

	t.Run("UpdateList() updates a list", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testMongoCollection.On("Update", id, mock.Anything).Return(nil).Once()

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

		testMongoCollection.On("Insert", &u).Return(nil).Once()
		err := store.AddUser(&u)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("AddUser() returns an error when the insert fails", func(t *testing.T) {
		u := models.SampleUser()

		testMongoCollection.On("Insert", &u).Return(errors.New("wadus")).Once()
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
