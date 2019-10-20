package services

import (
	"github.com/AngelVlc/lists-backend/stores"
	"github.com/stretchr/testify/mock"
)

type mockedMongoSession struct {
	mock.Mock
}

func (m *mockedMongoSession) GetRepository(collectionName string) stores.Repository {
	args := m.Called(collectionName)
	return args.Get(0).(stores.Repository)
}

type mockedRepository struct {
	mock.Mock
}

func (m *mockedRepository) Get(doc interface{}, query interface{}, selector interface{}) error {
	args := m.Called(doc, query, selector)
	return args.Error(0)
}

func (m *mockedRepository) GetOne(doc interface{}, query interface{}, selector interface{}) error {
	args := m.Called(doc, query, selector)
	return args.Error(0)
}

func (m *mockedRepository) Remove(query interface{}) error {
	args := m.Called(query)
	return args.Error(0)
}

func (m *mockedRepository) Update(query interface{}, doc interface{}) error {
	args := m.Called(query, doc)
	return args.Error(0)
}

func (m *mockedRepository) Add(doc interface{}) (string, error) {
	args := m.Called(doc)
	return args.String(0), args.Error(1)
}

func (m *mockedRepository) IsValidID(id string) bool {
	args := m.Called(id)
	return args.Bool(0)
}
