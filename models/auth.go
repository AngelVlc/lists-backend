package models

// Login is the model used for login
type Login struct {
	UserName     string `json:"userName"`
	PasswordHash string `json:"passwordHash"`
}
