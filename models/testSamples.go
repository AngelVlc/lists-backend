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

func SampleListSlice() []List {
	list1 := List{
		ID:   bson.NewObjectId().Hex(),
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

	list2 := List{
		ID:   bson.NewObjectId().Hex(),
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

	return []List{list1, list2}
}

func SampleGetListsResultDto() []GetListsResultDto {
	r := []GetListsResultDto{}

	for _, v := range SampleListSlice() {
		dto := GetListsResultDto{
			ID:   v.ID,
			Name: v.Name,
		}
		r = append(r, dto)
	}

	return r
}

func SampleUser() User {
	return User{
		UserName:     "user1",
		PasswordHash: "pass",
		IsAdmin:      true,
	}
}
