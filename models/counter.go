package models

// Counter is the model used for store a counter in the database
type Counter struct {
	ID    string `bson:"_id"`
	Name  string `bson:"name"`
	Value int    `bson:"value"`
}
