package controllers

import (
	"github.com/AngelVlc/lists-backend/services"
	"github.com/stretchr/testify/mock"
)

type mockedServiceProvider struct {
	mock.Mock
}

func (sp *mockedServiceProvider) GetUsersService() services.UsersService {
	args := sp.Called()
	return args.Get(0).(services.UsersService)
}

func (sp *mockedServiceProvider) GetListsService() services.ListsService {
	args := sp.Called()
	return args.Get(0).(services.ListsService)
}

func (sp *mockedServiceProvider) GetAuthService() services.AuthService {
	args := sp.Called()
	return args.Get(0).(services.AuthService)
}
