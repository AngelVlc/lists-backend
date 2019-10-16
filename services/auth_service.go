package services

import (
	"time"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
)

// AuthService is the interface an auth service must implement
type AuthService interface {
	CreateTokens(u *models.User) (map[string]string, error)
	ParseToken(token string) (*models.JwtClaimsInfo, error)
	ParseRefreshToken(refreshTokenString string) (*models.RefreshTokenClaimsInfo, error)
}

// MyAuthService is the service for auth methods
type MyAuthService struct {
	jwtPrv JwtProvider
}

// NewMyAuthService returns a new auth service
func NewMyAuthService(jwtp JwtProvider) *MyAuthService {
	return &MyAuthService{jwtp}
}

// CreateTokens returns a new jwt token and a refresh token for the given user
func (s *MyAuthService) CreateTokens(u *models.User) (map[string]string, error) {
	t := s.jwtPrv.NewToken()

	tc := s.jwtPrv.GetTokenClaims(t)
	tc["userName"] = u.UserName
	tc["isAdmin"] = u.IsAdmin
	tc["userId"] = u.ID
	tc["exp"] = time.Now().Add(time.Minute * 15).Unix()

	st, err := s.jwtPrv.SignToken(t)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt token", InternalError: err}
	}

	rt := s.jwtPrv.NewToken()
	rtc := s.jwtPrv.GetTokenClaims(rt)
	rtc["userId"] = u.ID
	rtc["exp"] = time.Now().Add(time.Hour * 24).Unix()

	srt, err := s.jwtPrv.SignToken(rt)
	if err != nil {
		return nil, &appErrors.UnexpectedError{Msg: "Error creating jwt refresh token", InternalError: err}
	}

	result := map[string]string{
		"token":        st,
		"refreshToken": srt,
	}

	return result, nil
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

// ParseRefreshToken takes a refresh token string, parses it and if it is valid returns a
// RefreshTokenClaimsInfo with its claims values
func (s *MyAuthService) ParseRefreshToken(refreshTokenString string) (*models.RefreshTokenClaimsInfo, error) {
	refreshToken, err := s.jwtPrv.ParseToken(refreshTokenString)
	if err != nil {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid refresh token", InternalError: err}
	}

	if !s.jwtPrv.IsTokenValid(refreshToken) {
		return nil, &appErrors.UnauthorizedError{Msg: "Invalid refresh token"}
	}

	return s.jwtPrv.GetRefreshTokenInfo(refreshToken), nil
}
