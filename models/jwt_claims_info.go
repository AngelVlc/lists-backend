package models

// JwtClaimsInfo is the struct which contains the jwt claims values
type JwtClaimsInfo struct {
	ID       string
	UserName string
	IsAdmin  bool
}
