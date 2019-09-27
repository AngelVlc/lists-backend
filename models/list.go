package models

// List is the model for the list
type List struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Items []Item `json:"items" bson:"items"`
}
