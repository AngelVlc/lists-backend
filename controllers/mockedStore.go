package controllers

import (
	"github.com/stretchr/testify/mock"
	"github.com/AngelVlc/lists-backend/models"
)

type mockedStore struct {
	mock.Mock
}

func (m *mockedStore) GetLists() ([]models.GetListsResultDto, error) {
	args := m.Called()
	return args.Get(0).([]models.GetListsResultDto), args.Error(1)
}

func (m *mockedStore) AddList(l *models.List) error {
	args := m.Called(l)
	return args.Error(0)
}

func (m *mockedStore) RemoveList(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockedStore) UpdateList(id string, l *models.List) error {
	args := m.Called(id, l)
	return args.Error(0)
}

func (m *mockedStore) GetSingleList(id string) (models.List, error) {
	args := m.Called(id)
	return args.Get(0).(models.List), args.Error(1)
}

func (m *mockedStore) AddUser(u *models.User) error {
	args := m.Called(u)
	return args.Error(0)
}