package models

import (
	"gopkg.in/mgo.v2/bson"
)

// ListDto is the struct used as DTO for a List
type ListDto struct {
	Name  string
	Items []Item
}

// ToList returns a List from the Dto
func (dto *ListDto) ToList() List {
	return List{
		Name:  dto.Name,
		Items: dto.Items,
	}
}

// GetListsResultDto is the struct used as result for the GetLists method
type GetListsResultDto struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name" bson:"name"`
}