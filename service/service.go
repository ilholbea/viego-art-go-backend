package service

import (
	"github.com/ilholbea/viego-art/types"
)

type Database interface {
	CreateAccount(Account *types.Account) error
	GetAllAccounts() ([]types.Account, error)
	GetAccount(id string) (*types.Account, error)
	Login(email, hashedPassword *string) (*types.Account, error)
	UpdateUser(exercise *types.Account) error
	DeleteUser(id string) error
}
