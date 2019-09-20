// Package models contains the models for the api
package models

import (
	"gopkg.in/mgo.v2/bson"
)

// List is the model for the list
type List struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Name  string        `json:"name" bson:"name"`
	Items []Item        `json:"items" bson:"items"`
}
