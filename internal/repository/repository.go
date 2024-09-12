package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/ilholbea/viego-art-backend/pkg/utils"
	"github.com/ilholbea/viego-art-backend/types"
)

const (
	NonAdminUserType = "non-admin"
)

type UserRepository interface {
	RetrieveUser(ctx context.Context, receivedUser *types.User) (*types.User, error)
	CreateUser(ctx context.Context, user *types.User) error
	DeleteUser(ctx context.Context, id int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (ar *userRepository) RetrieveUser(ctx context.Context, receivedUser *types.User) (*types.User, error) {
	foundUser := &types.User{}
	query := `SELECT id, username, email, password_hash, user_type, created_at FROM users usr WHERE usr.email = $1`
	err := ar.db.QueryRowContext(ctx, query, receivedUser.Email).Scan(
		&foundUser.ID,
		&foundUser.Username,
		&foundUser.Email,
		&foundUser.Password,
		&foundUser.Type,
		&foundUser.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if !utils.ComparePasswords(foundUser.Password, receivedUser.Password) {
		return nil, fmt.Errorf(`user "%s" does not match provided password`, receivedUser.Email)
	}

	return foundUser, nil
}

func (ar *userRepository) CreateUser(ctx context.Context, user *types.User) error {
	query := `INSERT INTO users (id, username, email, password_hash, user_type) VALUES ($1, $2, $3, $4, $5)`

	// Hash the password and update the User object
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	// Execute query with context
	_, err = ar.db.ExecContext(ctx, query, uuid.New(), user.Username, user.Email, hashedPassword, NonAdminUserType)
	if err != nil {
		return err
	}

	return nil
}

func (ar *userRepository) DeleteUser(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
