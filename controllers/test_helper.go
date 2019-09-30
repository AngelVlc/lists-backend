package controllers

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedStore struct {
	mock.Mock
}

func (m *mockedStore) Get(doc interface{}) error {
	args := m.Called(doc)
	return args.Error(0)
}

func (m *mockedStore) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *mockedStore) Update(id string, doc interface{}) error {
	args := m.Called(id, doc)
	return args.Error(0)
}

func (m *mockedStore) GetSingle(id string, item interface{}) error {
	args := m.Called(id, item)
	return args.Error(0)
}

func (m *mockedStore) Add(doc interface{}) error {
	args := m.Called(doc)
	return args.Error(0)
}

func assertResult(t *testing.T, m *mockedStore, got, want int) {
	t.Helper()

	assert.Equal(t, want, got, "status are not equal")

	m.AssertExpectations(t)
}