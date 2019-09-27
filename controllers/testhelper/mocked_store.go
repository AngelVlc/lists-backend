package testhelper

import (
	"github.com/stretchr/testify/mock"
	"github.com/AngelVlc/lists-backend/models"
)

// MockedStore is the object used as a mock for a store
type MockedStore struct {
	mock.Mock
}

func (m *MockedStore) GetLists() ([]models.GetListsResultDto, error) {
	args := m.Called()
	return args.Get(0).([]models.GetListsResultDto), args.Error(1)
}

func (m *MockedStore) AddList(l *models.List) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *MockedStore) RemoveList(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockedStore) UpdateList(id string, l *models.List) error {
	args := m.Called(id, l)
	return args.Error(0)
}

func (m *MockedStore) GetSingleList(id string) (models.List, error) {
	args := m.Called(id)
	return args.Get(0).(models.List), args.Error(1)
}

func (m *MockedStore) AddUser(u *models.User) error {
	args := m.Called(u)
	return args.Error(0)
}