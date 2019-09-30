package models

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

// GetListsResultDto is the struct used as result for the Get method
type GetListsResultDto struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

// UserDto is the struct used as DTO for a user
type UserDto struct {
	UserName           string
	NewPassword        string
	ConfirmNewPassword string
	IsAdmin            bool
}

// ToUser returns a User from the Dto
func (dto *UserDto) ToUser() User {
	return User{
		UserName: dto.UserName,
		IsAdmin:  dto.IsAdmin,
	}
}

// GetUsersResultDto is the struct used as result for the GetUsers method
type GetUsersResultDto struct {
	ID       string `json:"id" bson:"_id"`
	UserName string `json:"userName" bson:"userName"`
	IsAdmin  bool   `json:"isAdmin" bson:"isAdmin"`
}
