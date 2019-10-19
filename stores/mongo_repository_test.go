package stores

import (
	"errors"
	"testing"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2/bson"
)

type MockedMongoCollection struct {
	mock.Mock
}

func (m *MockedMongoCollection) Find(doc interface{}, query interface{}, selector interface{}) error {
	args := m.Called(doc, query, selector)
	return args.Error(0)
}

func (m *MockedMongoCollection) GetOne(doc interface{}, query interface{}, selector interface{}) error {
	args := m.Called(doc, query, selector)
	return args.Error(0)
}

func (m *MockedMongoCollection) FindOne(doc interface{}, query interface{}, selector interface{}) error {
	args := m.Called(doc, query, selector)
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

func (m *MockedMongoCollection) Name() string {
	args := m.Called()
	return args.String(0)
}

func TestUpdate(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	repository := MongoRepository{testMongoCollection}

	t.Run("Update() returns an unexpected error when the update fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		l := models.SampleList()

		testMongoCollection.On("Update", id, &l).Return(errors.New("wadus")).Once()

		err := repository.Update(id, &l)

		assert.IsType(t, &appErrors.UnexpectedError{}, err)

		assertFailedOperation(t, testMongoCollection, err, "Error updating the database")
	})

	t.Run("Update() returns a not found error when document does not exits", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		l := models.SampleList()
		testMongoCollection.On("Update", id, &l).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		err := repository.Update(id, &l)

		assert.IsType(t, &appErrors.NotFoundError{}, err)

		assertFailedOperation(t, testMongoCollection, err, "document not found")
	})

	t.Run("Update() updates a list", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		l := models.SampleList()

		testMongoCollection.On("Update", id, &l).Return(nil).Once()

		err := repository.Update(id, &l)

		assertSuccededOperation(t, testMongoCollection, err)
	})
}

func TestRemove(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	repository := MongoRepository{testMongoCollection}

	t.Run("Remove() returns an unexpected error when the remove fails", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", id).Return(errors.New("wadus")).Once()

		err := repository.Remove(id)

		assert.IsType(t, &appErrors.UnexpectedError{}, err)

		assertFailedOperation(t, testMongoCollection, err, "Error removing from the database")
	})

	t.Run("Remove() returns a not found error when document does not exits", func(t *testing.T) {
		id := bson.NewObjectId().Hex()
		testMongoCollection.On("Remove", id).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		err := repository.Remove(id)

		assert.IsType(t, &appErrors.NotFoundError{}, err)

		assertFailedOperation(t, testMongoCollection, err, "document not found")
	})

	t.Run("Remove() removes a list", func(t *testing.T) {
		oidHex := bson.NewObjectId().Hex()

		testMongoCollection.On("Remove", oidHex).Return(nil).Once()

		err := repository.Remove(oidHex)

		assertSuccededOperation(t, testMongoCollection, err)
	})
}

func TestAdd(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	repository := MongoRepository{testMongoCollection}

	t.Run("Add() returns an unexpected error when the insert fails", func(t *testing.T) {
		l := models.SampleList()
		testMongoCollection.On("Insert", &l).Return(errors.New("wadus")).Once()

		id, err := repository.Add(&l)

		assert.Empty(t, id)
		assert.IsType(t, &appErrors.UnexpectedError{}, err)

		assertFailedOperation(t, testMongoCollection, err, "Error inserting in the database")
	})

	t.Run("Add() adds a new list", func(t *testing.T) {
		l := models.SampleList()

		testMongoCollection.On("Insert", &l).Return(nil).Once()
		id, err := repository.Add(&l)

		assert.NotEmpty(t, id)

		assertSuccededOperation(t, testMongoCollection, err)
	})
}

func TestGet(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	repository := MongoRepository{testMongoCollection}

	t.Run("Get() returns the collection items", func(t *testing.T) {
		data := models.SampleGetListsResultDto()
		testMongoCollection.On("Find", &[]models.GetListsResultDto{}, nil, bson.M{"name": 1}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*[]models.GetListsResultDto)
			*arg = data
		})

		want := data
		got := []models.GetListsResultDto{}
		err := repository.Get(&got, nil, bson.M{"name": 1})

		assert.Equal(t, want, got, "they should be equal")

		assertSuccededOperation(t, testMongoCollection, err)
	})

	t.Run("Get() returns an error when the query fails", func(t *testing.T) {
		testMongoCollection.On("Find", &[]models.GetListsResultDto{}, nil, bson.M{"name": 1}).Return(errors.New("wadus")).Once()

		r := []models.GetListsResultDto{}
		err := repository.Get(&r, nil, bson.M{"name": 1})

		assertFailedOperation(t, testMongoCollection, err, "Error retrieving from the database")
	})
}

func TestGetOne(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)

	repository := MongoRepository{testMongoCollection}

	t.Run("GetOne() returns one item", func(t *testing.T) {
		data := models.List{
			ID:   "id",
			Name: "list1",
		}
		testMongoCollection.On("FindOne", &models.List{}, nil, bson.M{"name": "list1"}).Return(nil).Once().Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.List)
			*arg = data
		})

		want := data
		got := models.List{}
		err := repository.GetOne(&got, nil, bson.M{"name": "list1"})

		assert.Equal(t, want, got, "they should be equal")

		assertSuccededOperation(t, testMongoCollection, err)
	})

	t.Run("GetOne() returns an error when the query fails", func(t *testing.T) {
		testMongoCollection.On("FindOne", &models.List{}, nil, bson.M{"name": "1"}).Return(errors.New("wadus")).Once()

		r := models.List{}
		err := repository.GetOne(&r, nil, bson.M{"name": "1"})

		assert.IsType(t, &appErrors.UnexpectedError{}, err)

		assertFailedOperation(t, testMongoCollection, err, "Error retrieving from the database")
	})

	t.Run("GetOne() returns a not found error when document does not exits", func(t *testing.T) {
		testMongoCollection.On("FindOne", &models.List{}, nil, bson.M{"name": "1"}).Return(errors.New("not found")).Once()
		testMongoCollection.On("Name").Return("document").Once()

		r := models.List{}
		err := repository.GetOne(&r, nil, bson.M{"name": "1"})

		assert.IsType(t, &appErrors.NotFoundError{}, err)

		assertFailedOperation(t, testMongoCollection, err, "document not found")
	})
}

func TestIsValidID(t *testing.T) {
	testMongoCollection := new(MockedMongoCollection)
	repository := MongoRepository{testMongoCollection}

	t.Run("IsValidID() returns true if the id is valid", func(t *testing.T) {
		id := bson.NewObjectId().Hex()

		got := repository.IsValidID(id)

		want := true

		assert.Equal(t, want, got, "they should be equal")
	})

	t.Run("IsValidID() returns false if the id is not valid", func(t *testing.T) {
		got := repository.IsValidID("wadus")

		want := false

		assert.Equal(t, want, got, "they should be equal")
	})
}

func assertSuccededOperation(t *testing.T, c *MockedMongoCollection, err error) {
	c.AssertExpectations(t)

	assert.Nil(t, err)
}

func assertFailedOperation(t *testing.T, c *MockedMongoCollection, err error, errorMsg string) {
	c.AssertExpectations(t)

	assert.NotNil(t, err)

	assert.Equal(t, errorMsg, err.Error())
}

func sampleList() models.List {
	l := models.SampleListSlice()[0]
	l.ID = bson.NewObjectId().Hex()
	return l
}
