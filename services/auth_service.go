package services

import (
	"time"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
)

// AuthService is the interface an auth service must implement
type AuthService interface {
	CreateToken(u *models.User) (string, error)
	ParseToken(token string) (*models.JwtClaimsInfo, error)
}

// MyAuthService is the service for auth methods
type MyAuthService struct {
	jwtPrv JwtProvider
}

// NewMyAuthService returns a new auth service
func NewMyAuthService(jwtp JwtProvider) *MyAuthService {
	return &MyAuthService{jwtp}
}

// CreateToken returns a new jwt token for the given user
func (s *MyAuthService) CreateToken(u *models.User) (string, error) {
	token := s.jwtPrv.NewToken()

	claims := s.jwtPrv.GetTokenClaims(token)
	claims["userName"] = u.UserName
	claims["isAdmin"] = u.IsAdmin
	claims["userId"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := s.jwtPrv.SignToken(token)
	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	return signedToken, nil
}

// ParseToken takes a token string, parses it and if it is valid returns a JwtClaimsInfo
// with its claims values
func (s *MyAuthService) ParseToken(tokenString string) (*models.JwtClaimsInfo, error) {
	token, err := s.jwtPrv.ParseToken(tokenString)
	if err != nil {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid token", InternalError: err}
	}

	if !s.jwtPrv.IsTokenValid(token) {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid token"}
	}

	return s.jwtPrv.GetJwtInfo(token), nil
}
