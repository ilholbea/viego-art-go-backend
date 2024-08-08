package services

import (
	"context"
	"github.com/ilholbea/viego-art-backend/repository"
	"github.com/ilholbea/viego-art-backend/types"
)

type AccountService interface {
	GetAccount(ctx context.Context, email string) (types.Account, error)
	Create(ctx context.Context, account types.Account) error
	Delete(ctx context.Context, id int) error
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (as *accountService) GetAccount(ctx context.Context, email string) (types.Account, error) {
	//TODO implement me
	panic("implement me")
}

func (as *accountService) Create(ctx context.Context, account types.Account) error {
	//TODO implement me
	panic("implement me")
}

func (as *accountService) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
