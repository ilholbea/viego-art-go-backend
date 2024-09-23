package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilholbea/viego-art-backend/internal/repository"
	"github.com/ilholbea/viego-art-backend/pkg/utils"
	"github.com/ilholbea/viego-art-backend/types"
)

type UserService interface {
	GetUser(ctx context.Context, receivedUser types.User) (*types.User, error)
	CreateUser(ctx context.Context, user *types.User) error
	DeleteUser(ctx context.Context, id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (as *userService) GetUser(ctx context.Context, receivedUser types.User) (*types.User, error) {
	if receivedUser.Email == "" {
		return nil, fmt.Errorf("no email has been provided")
	}

	foundUser, err := as.repo.RetrieveUser(ctx, &receivedUser)
	if err != nil {
		return nil, err
	}

	return foundUser, nil
}

func (as *userService) CreateUser(ctx context.Context, user *types.User) error {
	// Validate the attributes of the new User
	if ok := utils.ValidateEmailFormat(user.Email); !ok {
		return errors.New("invalid email format")
	}

	// Call the repository to create the User
	err := as.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (as *userService) DeleteUser(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
