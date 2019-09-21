package models

import (
	"gopkg.in/mgo.v2/bson"
)

// User is the model for a user
type User struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	UserName     string `json:"userName" bson:"userName"`
	PasswordHash string `json:"passwordHash" bson:"passwordHash"`
	IsAdmin      bool   `json:"isAdmin" bson:"isAdmin"`
}
