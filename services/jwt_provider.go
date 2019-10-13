package services

import (
	"time"

	"github.com/AngelVlc/lists-backend/models"
	"github.com/dgrijalva/jwt-go"
)

type JwtProvider interface {
	CreateToken(u *models.User) (string, error)
}

type MyJwtProvider struct {
	secret string
}

func NewMyJwtProvider(secret string) *MyJwtProvider {
	return &MyJwtProvider{secret}
}

func (p *MyJwtProvider) CreateToken(u *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userName"] = u.UserName
	claims["userId"] = u.ID
	claims["isAdmin"] = u.IsAdmin
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(p.secret))
	if err != nil {
		return "", err
	}

	return t, nil
}
