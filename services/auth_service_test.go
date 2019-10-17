package services

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	appErrors "github.com/AngelVlc/lists-backend/errors"
	"github.com/AngelVlc/lists-backend/models"
	"github.com/stretchr/testify/mock"
)

type mockedJwtProvider struct {
	mock.Mock
}

func (m *mockedJwtProvider) NewToken() interface{} {
	args := m.Called()
	return args.Get(0).(interface{})
}

func (m *mockedJwtProvider) GetTokenClaims(token interface{}) map[string]interface{} {
	args := m.Called(token)
	return args.Get(0).(map[string]interface{})
}

func (m *mockedJwtProvider) SignToken(token interface{}) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *mockedJwtProvider) ParseToken(tokenString string) (interface{}, error) {
	args := m.Called(tokenString)

	got := args.Get(0)

	if got == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(interface{}), args.Error(1)
}

func (m *mockedJwtProvider) IsTokenValid(token interface{}) bool {
	args := m.Called(token)
	return args.Bool(0)
}

func TestAuthServiceCreateToken(t *testing.T) {
	mockedJwtProvider := new(mockedJwtProvider)

	service := NewMyAuthService(mockedJwtProvider)

	u := models.User{}
	token := struct{}{}
	claims := map[string]interface{}{}
	refreshToken := struct{}{}
	refreshTokenClaims := map[string]interface{}{}

	t.Run("should return an UnexpectedError if sign token fails", func(t *testing.T) {
		mockedJwtProvider.On("NewToken").Return(token).Once()
		mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
		mockedJwtProvider.On("SignToken", token).Return("", errors.New("wadus")).Once()

		tokens, err := service.CreateTokens(&u)

		assert.Nil(t, tokens)
		assert.NotNil(t, err)
		unexpectedErr, isUnexpectedErr := err.(*appErrors.UnexpectedError)
		assert.Equal(t, true, isUnexpectedErr, "should be an unexpected error")
		assert.Equal(t, "Error creating jwt token", unexpectedErr.Error())
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return an UnexpectedError if sign refresh token fails", func(t *testing.T) {
		mockedJwtProvider.On("NewToken").Return(token).Once()
		mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
		mockedJwtProvider.On("NewToken").Return(refreshToken).Once()
		mockedJwtProvider.On("GetTokenClaims", refreshToken).Return(refreshTokenClaims).Once()
		mockedJwtProvider.On("SignToken", token).Return("token", nil).Once()
		mockedJwtProvider.On("SignToken", refreshToken).Return("", errors.New("wadus")).Once()

		tokens, err := service.CreateTokens(&u)

		assert.Nil(t, tokens)
		assert.NotNil(t, err)
		unexpectedErr, isUnexpectedErr := err.(*appErrors.UnexpectedError)
		assert.Equal(t, true, isUnexpectedErr, "should be an unexpected error")
		assert.Equal(t, "Error creating jwt refresh token", unexpectedErr.Error())
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return a signed token if no error happen", func(t *testing.T) {
		theToken := "theToken"
		theRefreshToken := "theRefreshToken"
		mockedJwtProvider.On("NewToken").Return(token).Once()
		mockedJwtProvider.On("GetTokenClaims", token).Return(claims).Once()
		mockedJwtProvider.On("NewToken").Return(refreshToken).Once()
		mockedJwtProvider.On("GetTokenClaims", refreshToken).Return(refreshTokenClaims).Once()
		mockedJwtProvider.On("SignToken", token).Return(theToken, nil).Once()
		mockedJwtProvider.On("SignToken", refreshToken).Return(theRefreshToken, nil).Once()

		tokens, err := service.CreateTokens(&u)

		want := map[string]string{
			"token":        theToken,
			"refreshToken": theRefreshToken,
		}

		assert.Equal(t, want, tokens)
		assert.Nil(t, err)

		mockedJwtProvider.AssertExpectations(t)
	})
}

func TestAuthServiceParseToken(t *testing.T) {
	mockedJwtProvider := new(mockedJwtProvider)

	service := NewMyAuthService(mockedJwtProvider)

	theToken := "theToken"

	t.Run("should return an unathorized error when jwt ParseToken() fails", func(t *testing.T) {
		mockedJwtProvider.On("ParseToken", theToken).Return(nil, errors.New("wadus")).Once()

		jwtInfo, err := service.ParseToken(theToken)

		assert.Nil(t, jwtInfo)
		assert.NotNil(t, err)
		unauthErr, isUnauthErr := err.(*appErrors.UnauthorizedError)
		assert.Equal(t, true, isUnauthErr, "should be an unauthorized error")
		assert.Equal(t, "Invalid token", unauthErr.Error())
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return an unauthorized error when the jwt IsTokenValid() return false", func(t *testing.T) {
		token := struct{}{}

		mockedJwtProvider.On("ParseToken", theToken).Return(token, nil).Once()
		mockedJwtProvider.On("IsTokenValid", token).Return(false).Once()

		jwtInfo, err := service.ParseToken(theToken)

		assert.Nil(t, jwtInfo)
		assert.NotNil(t, err)
		unauthErr, isUnauthErr := err.(*appErrors.UnauthorizedError)
		assert.Equal(t, true, isUnauthErr, "should be an unauthorized error")
		assert.Equal(t, "Invalid token", unauthErr.Error())
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return a jwt info when the token is valid", func(t *testing.T) {
		token := struct{}{}

		jwtInfo := models.JwtClaimsInfo{
			UserName: "wadus",
			IsAdmin:  true,
			ID:       "11",
		}

		mockedJwtProvider.On("ParseToken", theToken).Return(token, nil).Once()
		mockedJwtProvider.On("IsTokenValid", token).Return(true).Once()

		c := map[string]interface{}{
			"userName": "wadus",
			"isAdmin":  true,
			"userId":   "11",
		}
		mockedJwtProvider.On("GetTokenClaims", token).Return(c).Once()

		res, err := service.ParseToken(theToken)

		assert.Equal(t, &jwtInfo, res)
		assert.Nil(t, err)
		mockedJwtProvider.AssertExpectations(t)
	})
}

func TestAuthServiceParseRefreshToken(t *testing.T) {
	mockedJwtProvider := new(mockedJwtProvider)

	service := NewMyAuthService(mockedJwtProvider)

	theRefreshToken := "theRefreshToken"

	t.Run("should return an unathorized error when jwt ParseToken() fails", func(t *testing.T) {
		mockedJwtProvider.On("ParseToken", theRefreshToken).Return(nil, errors.New("wadus")).Once()

		jwtInfo, err := service.ParseRefreshToken(theRefreshToken)

		assert.Nil(t, jwtInfo)
		assert.NotNil(t, err)
		unauthErr, isUnauthErr := err.(*appErrors.UnauthorizedError)
		assert.Equal(t, true, isUnauthErr, "should be an unauthorized error")
		assert.Equal(t, "Invalid refresh token", unauthErr.Error())
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return an unauthorized error when the jwt IsTokenValid() return false", func(t *testing.T) {
		refreshToken := struct{}{}

		mockedJwtProvider.On("ParseToken", theRefreshToken).Return(refreshToken, nil).Once()
		mockedJwtProvider.On("IsTokenValid", refreshToken).Return(false).Once()

		rtInfo, err := service.ParseRefreshToken(theRefreshToken)

		assert.Nil(t, rtInfo)
		assert.NotNil(t, err)
		unauthErr, isUnauthErr := err.(*appErrors.UnauthorizedError)
		assert.Equal(t, true, isUnauthErr, "should be an unauthorized error")
		assert.Equal(t, "Invalid refresh token", unauthErr.Error())
		mockedJwtProvider.AssertExpectations(t)
	})

	t.Run("should return a refresh token claims info when the token is valid", func(t *testing.T) {
		refreshToken := struct{}{}

		rtInfo := models.RefreshTokenClaimsInfo{
			ID: "id",
		}

		mockedJwtProvider.On("ParseToken", theRefreshToken).Return(refreshToken, nil).Once()
		mockedJwtProvider.On("IsTokenValid", refreshToken).Return(true).Once()

		c := map[string]interface{}{
			"userId": rtInfo.ID,
		}
		mockedJwtProvider.On("GetTokenClaims", refreshToken).Return(c).Once()

		res, err := service.ParseRefreshToken(theRefreshToken)

		assert.Equal(t, &rtInfo, res)
		assert.Nil(t, err)
		mockedJwtProvider.AssertExpectations(t)
	})
}

func TestAuthServiceJwtProviderIntegration(t *testing.T) {
	jwtPrv := NewMyJwtProvider("theSecret")

	service := NewMyAuthService(jwtPrv)

	u := models.User{
		UserName: "wadus",
		IsAdmin:  true,
		ID:       "theId",
	}

	tokens, err := service.CreateTokens(&u)
	assert.NotNil(t, tokens)
	assert.Nil(t, err)

	jwtInfo, err := service.ParseToken(tokens["token"])
	assert.NotNil(t, jwtInfo)
	assert.Nil(t, err)

	assert.Equal(t, u.UserName, jwtInfo.UserName)
	assert.Equal(t, u.IsAdmin, jwtInfo.IsAdmin)
	assert.Equal(t, u.ID, jwtInfo.ID)

	rtClaims, err := service.ParseRefreshToken(tokens["refreshToken"])
	assert.NotNil(t, rtClaims)
	assert.Nil(t, err)

	assert.Equal(t, u.ID, rtClaims.ID)
}
