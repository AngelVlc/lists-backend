package models

// User is the model for a user
type User struct {
	ID           string `json:"id" bson:"_id"`
	UserName     string `json:"userName" bson:"userName"`
	PasswordHash string `json:"passwordHash" bson:"passwordHash"`
	IsAdmin      bool   `json:"isAdmin" bson:"isAdmin"`
}
