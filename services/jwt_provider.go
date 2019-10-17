package services

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

// JwtProvider is the interface which contains the methods to deal with a Jwt token
type JwtProvider interface {
	NewToken() interface{}
	GetTokenClaims(token interface{}) map[string]interface{}
	SignToken(token interface{}) (string, error)
	ParseToken(tokenString string) (interface{}, error)
	IsTokenValid(token interface{}) bool
}

// MyJwtProvider is the type used as JwtProvider
type MyJwtProvider struct {
	secret string
}

// NewMyJwtProvider returns a new MyJwtProvider
func NewMyJwtProvider(secret string) *MyJwtProvider {
	return &MyJwtProvider{secret}
}

// NewToken returns a new Jwt tooken
func (p *MyJwtProvider) NewToken() interface{} {
	return jwt.New(jwt.SigningMethodHS256)
}

func (p *MyJwtProvider) getJwtToken(token interface{}) *jwt.Token {
	jwtToken, _ := token.(*jwt.Token)
	return jwtToken
}

// GetTokenClaims returns the claims for the given token as a map
func (p *MyJwtProvider) GetTokenClaims(token interface{}) map[string]interface{} {
	return p.getJwtToken(token).Claims.(jwt.MapClaims)
}

// SignToken signs the given token
func (p *MyJwtProvider) SignToken(token interface{}) (string, error) {
	return p.getJwtToken(token).SignedString([]byte(p.secret))
}

// ParseToken parses the string and checks the signing method
func (p *MyJwtProvider) ParseToken(tokenString string) (interface{}, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secret), nil
	})
}

// IsTokenValid returns true if the given token is valid
func (p *MyJwtProvider) IsTokenValid(token interface{}) bool {
	return p.getJwtToken(token).Valid
}
