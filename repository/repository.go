package repository

import (
	"context"
	"database/sql"

	"github.com/ilholbea/viego-art-backend/types"
)

type AccountRepository interface {
	RetrieveAccount(ctx context.Context, email string) (types.Account, error)
	CreateAccount(ctx context.Context, account types.Account) error
	DeleteAccount(ctx context.Context, id int) error
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (ar *accountRepository) RetrieveAccount(ctx context.Context, email string) (types.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (ar *accountRepository) CreateAccount(ctx context.Context, account types.Account) error {
	//TODO implement me
	panic("implement me")
}

func (ar *accountRepository) DeleteAccount(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
