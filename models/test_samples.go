package models

import "gopkg.in/mgo.v2/bson"

func SampleList() List {
	l := List{
		Name: "added",
		Items: []Item{
			Item{
				Title: "newItem",
			},
		},
	}

	return l
}

func SampleListCollection() map[string]List {
	m := make(map[string]List, 2)

	m["1"] = List{
		ID:   bson.NewObjectId(),
		Name: "list1",
		Items: []Item{
			Item{
				Title:       "item11",
				Description: "this is the first item",
			},
			Item{
				Title:       "item12",
				Description: "this is the second item",
			},
		},
	}

	m["2"] = List{
		ID:   bson.NewObjectId(),
		Name: "list2",
		Items: []Item{
			Item{
				Title:       "item21",
				Description: "this is the first item",
			},
			Item{
				Title:       "item22",
				Description: "this is the second item",
			},
		},
	}

	return m
}

func SampleListCollectionSlice() []List {
	r := []List{}

	for _, v := range SampleListCollection() {
		r = append(r, v)
	}

	return r
}
