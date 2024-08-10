package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/ilholbea/viego-art-backend/types"
)

const (
	UserAccount = "user"
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
	query := `INSERT INTO accounts (id, name, email, password, account_type) VALUES ($1, $2, $3, $4, $5)`

	// Execute query with context
	_, err := ar.db.ExecContext(ctx, query, uuid.New(), account.Name, account.Email, account.Password, UserAccount)
	if err != nil {
		return err
	}

	return nil
}

func (ar *accountRepository) DeleteAccount(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
