package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AngelVlc/lists-backend/models"
	"github.com/dgrijalva/jwt-go"
)

type JwtProvider interface {
	CreateToken(m map[string]interface{}) (string, error)
	ValidateToken(token string) (*models.JwtClaimsInfo, error)
}

type MyJwtProvider struct {
	secret string
}

func NewMyJwtProvider(secret string) *MyJwtProvider {
	return &MyJwtProvider{secret}
}

func (p *MyJwtProvider) CreateToken(m map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	for k, v := range m {
		claims[k] = v
	}
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(p.secret))
	if err != nil {
		return "", err
	}

	return t, nil
}

func (p *MyJwtProvider) ValidateToken(tokenString string) (*models.JwtClaimsInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(p.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		info := models.JwtClaimsInfo{
			UserName: parseStringClaim(claims["userName"]),
			ID:       parseStringClaim(claims["userId"]),
			IsAdmin:  parseBoolClaim(claims["isAdmin"]),
		}
		return &info, nil
	}
	return nil, errors.New("Invalid token")
}

func parseStringClaim(value interface{}) string {
	result, _ := value.(string)
	return result
}

func parseBoolClaim(value interface{}) bool {
	result, _ := value.(bool)
	return result
}
