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

func (m *mockedRepository) Get(doc interface{}) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *mockedRepository) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockedRepository) Update(id string, doc interface{}) error {
	args := m.Called(id, doc)
	return args.Error(0)
}

func (m *mockedRepository) GetSingle(id string, item interface{}) error {
	args := m.Called(id, item)
	return args.Error(0)
}

func (m *mockedRepository) Add(doc interface{}) (string, error) {
	args := m.Called(doc)
	return args.String(0), args.Error(1)
}
