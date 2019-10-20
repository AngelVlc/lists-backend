package models

// JwtClaimsInfo is the struct which contains the jwt token claims
type JwtClaimsInfo struct {
	UserID       string
	UserName string
	IsAdmin  bool
}

// RefreshTokenClaimsInfo is the struct which contains the refresh token claims
type RefreshTokenClaimsInfo struct {
	UserID string
}
