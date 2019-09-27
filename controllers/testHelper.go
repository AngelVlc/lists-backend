package controllers

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func assertResult(t *testing.T, m *mockedStore, got, want int) {
	t.Helper()

	assert.Equal(t, want, got, "status are not equal")

	m.AssertExpectations(t)
}