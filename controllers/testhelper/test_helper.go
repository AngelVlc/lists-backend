package testhelper

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// AssertResult is used to assert the result of a test
func AssertResult(t *testing.T, m *MockedStore, got, want int) {
	t.Helper()

	assert.Equal(t, want, got, "status are not equal")

	m.AssertExpectations(t)
}