package services

import (
	"golang.org/x/crypto/bcrypt"
)

// BcryptProvider is the interface which contains the methods used to use bcrypt
type BcryptProvider interface {
	GenerateFromPassword(password []byte, cost int) ([]byte, error)
	CompareHashAndPassword(hashedPassword, password []byte) error
}

// MyBcryptProvider is the implementation for BcryptProvider and uses the real bcrypt package
type MyBcryptProvider struct{}

// NewMyBcryptProvider returns a new MyBryptProvider
func NewMyBcryptProvider() *MyBcryptProvider {
	return new(MyBcryptProvider)
}

// GenerateFromPassword generates a hashed password
func (b *MyBcryptProvider) GenerateFromPassword(password []byte, cost int) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

// CompareHashAndPassword checks if the given hashed password and the password matches
func (b *MyBcryptProvider) CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}
