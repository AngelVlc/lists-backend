package services

import (
	"github.com/AngelVlc/lists-backend/stores"
)

type ServiceProvider interface {
	GetUsersService() UsersService
	GetListsService() ListsService
}

type MyServiceProvider struct {
	session stores.MongoSession
}

func NewMyServiceProvider(s stores.MongoSession) *MyServiceProvider {
	return &MyServiceProvider{
		session: s,
	}
}

func (sp *MyServiceProvider) GetUsersService() UsersService {
	return NewMyUsersService(sp.session)
}

func (sp *MyServiceProvider) GetListsService() ListsService {
	return NewMyListsService(sp.session)
}
