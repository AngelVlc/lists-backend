package services

import (
	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
)

// AuthService is the interface an auth service must implement
type AuthService interface {
	CreateToken(u *models.User) (string, error)
	ValidateToken(token string) (*models.JwtClaimsInfo, error)
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
	claims := map[string]interface{}{
		"userName": u.UserName,
		"isAdmin":  u.IsAdmin,
		"userId":   u.ID,
	}

	t, err := s.jwtPrv.CreateToken(claims)

	if err != nil {
		return "", &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	return t, nil
}

func (s *MyAuthService) ValidateToken(token string) (*models.JwtClaimsInfo, error) {
	return s.jwtPrv.ValidateToken(token)
}
