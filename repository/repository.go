package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/ilholbea/viego-art-backend/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
)

const (
	UserAccount = "user"
)

type AccountRepository interface {
	RetrieveAccount(ctx context.Context, receivedAccount *types.Account) (*types.Account, error)
	CreateAccount(ctx context.Context, account *types.Account) error
	DeleteAccount(ctx context.Context, id int) error
}

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (ar *accountRepository) RetrieveAccount(ctx context.Context, receivedAccount *types.Account) (*types.Account, error) {
	foundAccount := &types.Account{}
	query := `SELECT id, name, email, password, account_type FROM accounts act WHERE act.email = $1`
	err := ar.db.QueryRowContext(ctx, query, receivedAccount.Email).Scan(
		&foundAccount.ID,
		&foundAccount.Name,
		&foundAccount.Email,
		&foundAccount.Password,
		&foundAccount.Type,
	)
	if err != nil {
		return nil, err
	}

	if !comparePasswords(foundAccount.Password, receivedAccount.Password) {
		return nil, fmt.Errorf(`account "%s" does not match provided password`, receivedAccount.Email)
	}

	return foundAccount, nil
}

func (ar *accountRepository) CreateAccount(ctx context.Context, account *types.Account) error {
	query := `INSERT INTO accounts (id, name, email, password, account_type) VALUES ($1, $2, $3, $4, $5)`

	// Hash the password and update the Account object
	hashedPassword, err := hashPassword(account.Password)
	if err != nil {
		return err
	}
	account.Password = hashedPassword

	log.Debugf("After password update account looks like: %v", account)

	// Execute query with context
	_, err = ar.db.ExecContext(ctx, query, uuid.New(), account.Name, account.Email, account.Password, UserAccount)
	if err != nil {
		return err
	}

	return nil
}

func (ar *accountRepository) DeleteAccount(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	var cost int
	hashCost, ok := os.LookupEnv("HASH_COST")
	if !ok {
		cost = bcrypt.DefaultCost
	}

	cost, err := strconv.Atoi(hashCost)
	if err != nil {
		cost = bcrypt.DefaultCost
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// comparePasswords compares a hashed password with its plain-text version
func comparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Return true if they match, false otherwise
}
