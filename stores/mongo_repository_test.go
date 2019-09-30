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

	store := NewMongoRepository(testMongoSession, testMongoCollection)

	t.Run("Update() returns an unexpected error when the update fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		l := models.SampleList()

		testMongoCollection.On("Update", id, &l).Return(errors.New("wadus")).Once()

		err := store.Update(id, &l)

		assert.IsType(t, &UnexpectedError{}, err)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error updating the database")
	})

	t.Run("Update() returns a not found error when document does not exits", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		l := models.SampleList()
		testMongoCollection.On("Update", id, &l).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		err := store.Update(id, &l)

		assert.IsType(t, &NotFoundError{}, err)

		msg := fmt.Sprintf("document with id %q not found", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("Update() returns an invalid id error when the id is not valid", func(t *testing.T) {
		id := "wadus"

		err := store.Update(id, &id)

		assert.IsType(t, &InvalidIDError{}, err)

		msg := fmt.Sprintf("%q is not a valid id", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("Update() updates a list", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		l := models.SampleList()

		testMongoCollection.On("Update", id, &l).Return(nil).Once()

		err := store.Update(id, &l)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})
}

func TestRemove(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)

	store := NewMongoRepository(testMongoSession, testMongoCollection)

	t.Run("Remove() returns an unexpected error when the remove fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", id).Return(errors.New("wadus")).Once()

		err := store.Remove(id)

		assert.IsType(t, &UnexpectedError{}, err)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error removing from the database")
	})

	t.Run("Remove() returns a not found error when document does not exits", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testMongoCollection.On("Remove", id).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		err := store.Remove(id)

		assert.IsType(t, &NotFoundError{}, err)

		msg := fmt.Sprintf("document with id %q not found", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("Remove() returns an invalid id error when the id is not valid", func(t *testing.T) {
		id := "wadus"

		err := store.Remove(id)

		assert.IsType(t, &InvalidIDError{}, err)

		msg := fmt.Sprintf("%q is not a valid id", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("Remove() removes a list", func(t *testing.T) {
		oidHex := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", oidHex).Return(nil).Once()

		err := store.Remove(oidHex)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})
}

func TestGetSingle(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)

	store := NewMongoRepository(testMongoSession, testMongoCollection)

	t.Run("GetSingle() returns an unexpected error when the remove fails", func(t *testing.T) {
		data := sampleList()
		testMongoCollection.On("FindOne", data.ID, &models.List{}).Return(errors.New("wadus")).Once()

		err := store.GetSingle(data.ID, &models.List{})

		assert.IsType(t, &UnexpectedError{}, err)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error retrieving from the database")
	})

	t.Run("GetSingle() returns a not found error when document does not exits", func(t *testing.T) {
		data := sampleList()
		testMongoCollection.On("FindOne", data.ID, &models.List{}).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		err := store.GetSingle(data.ID, &models.List{})

		assert.IsType(t, &NotFoundError{}, err)

		msg := fmt.Sprintf("document with id %q not found", data.ID)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("GetSingle() returns an invalid id error when the id is not valid", func(t *testing.T) {
		id := "wadus"

		err := store.GetSingle(id, &models.List{})

		assert.IsType(t, &InvalidIDError{}, err)

		msg := fmt.Sprintf("%q is not a valid id", id)
		assertFailedOperation(t, testMongoSession, testMongoCollection, err, msg)
	})

	t.Run("GetSingle() returns a single list", func(t *testing.T) {
		data := sampleList()
		testMongoCollection.On("FindOne", data.ID, &models.List{}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(1).(*models.List)
			*arg = data
		})

		want := data
		got := models.List{}
		err := store.GetSingle(data.ID, &got)

		assert.Equal(t, want, got, "they should be equal")

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})
}

func TestAdd(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)

	store := NewMongoRepository(testMongoSession, testMongoCollection)

	t.Run("Add() returns an unexpected error when the insert fails", func(t *testing.T) {
		l := models.SampleList()
		testMongoCollection.On("Insert", &l).Return(errors.New("wadus")).Once()

		err := store.Add(&l)

		assert.IsType(t, &UnexpectedError{}, err)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error inserting in the database")
	})

	t.Run("Add() adds a new list", func(t *testing.T) {
		l := models.SampleList()

		testMongoCollection.On("Insert", &l).Return(nil).Once()
		err := store.Add(&l)

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})
}

func TestGet(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	testMongoSession := new(MockedMongoSession)

	store := NewMongoRepository(testMongoSession, testMongoCollection)

	t.Run("Get() returns all the list items", func(t *testing.T) {
		data := models.SampleGetListsResultDto()
		testMongoCollection.On("Find", &[]models.GetListsResultDto{}, nil, bson.M{"name": 1}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]models.GetListsResultDto)
			*arg = data
		})

		want := data
		got := []models.GetListsResultDto{}
		err := store.Get(&got)

		assert.Equal(t, want, got, "they should be equal")

		assertSuccededOperation(t, testMongoSession, testMongoCollection, err)
	})

	t.Run("Get() returns an error when the query fails", func(t *testing.T) {
		testMongoCollection.On("Find", &[]models.GetListsResultDto{}, nil, bson.M{"name": 1}).Return(errors.New("wadus")).Once()

		r := []models.GetListsResultDto{}
		err := store.Get(&r)

		assertFailedOperation(t, testMongoSession, testMongoCollection, err, "Error retrieving from the database")
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

func sampleList() models.List {
	l := models.SampleListSlice()[0]
	l.ID = bson.NewObjectId().Hex()
	return l
}
