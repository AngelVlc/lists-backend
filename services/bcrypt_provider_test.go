package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBryptProvider(t *testing.T) {
	prv := NewMyBcryptProvider()

	password := "the_password"

	hashedBytes, err := prv.GenerateFromPassword([]byte(password), bcryptCost)

	assert.Nil(t, err)

	err = prv.CompareHashAndPassword(hashedBytes, []byte(password))

	assert.Nil(t, err)
}
