package services

import (
	"github.com/AngelVlc/lists-backend/stores"
)

type ServiceProvider interface {
	GetUsersService() UsersService
	GetListsService() ListsService
}

type MyServiceProvider struct {
	session   stores.MongoSession
	bcryptPrv BcryptProvider
}

func NewMyServiceProvider(s stores.MongoSession, bp BcryptProvider) *MyServiceProvider {
	return &MyServiceProvider{
		session:   s,
		bcryptPrv: bp,
	}
}

func (sp *MyServiceProvider) GetUsersService() UsersService {
	return NewMyUsersService(sp.session, sp.bcryptPrv)
}

func (sp *MyServiceProvider) GetListsService() ListsService {
	return NewMyListsService(sp.session)
}
