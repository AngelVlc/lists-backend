package models

// Item is the model for a single list item
type Item struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}
