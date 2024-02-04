package service

import (
	"github.com/ilholbea/viego-art/types"
)

type Database interface {
	CreateAccount(Account *types.Account) error
	GetAllAccounts() ([]types.Account, error)
	GetAccount(email string) (*types.Account, error)
	DeleteAccount(email string) error
	UpdateAccount(account *types.Account) error
	Login(email, hashedPassword string) (bool, error)
}
