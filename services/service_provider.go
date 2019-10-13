package services

import (
	"github.com/AngelVlc/lists-backend/stores"
)

type ServiceProvider interface {
	GetUsersService() UsersService
	GetListsService() ListsService
	GetAuthService() AuthService
}

type MyServiceProvider struct {
	session   stores.MongoSession
	bcryptPrv BcryptProvider
	jwtPrv    JwtProvider
}

func NewMyServiceProvider(s stores.MongoSession, bp BcryptProvider, jwtp JwtProvider) *MyServiceProvider {
	return &MyServiceProvider{
		session:   s,
		bcryptPrv: bp,
		jwtPrv:    jwtp,
	}
}

func (sp *MyServiceProvider) GetUsersService() UsersService {
	return NewMyUsersService(sp.session, sp.bcryptPrv)
}

func (sp *MyServiceProvider) GetListsService() ListsService {
	return NewMyListsService(sp.session)
}

func (sp *MyServiceProvider) GetAuthService() AuthService {
	return NewMyAuthService(sp.jwtPrv)
}
