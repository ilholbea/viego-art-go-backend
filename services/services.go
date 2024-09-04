package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ilholbea/viego-art-backend/repository"
	"github.com/ilholbea/viego-art-backend/types"
	"regexp"
)

type AccountService interface {
	GetAccount(ctx context.Context, receivedAccount types.Account) (*types.Account, error)
	Create(ctx context.Context, account *types.Account) error
	Delete(ctx context.Context, id int) error
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (as *accountService) GetAccount(ctx context.Context, receivedAccount types.Account) (*types.Account, error) {
	if receivedAccount.Email == "" {
		return nil, fmt.Errorf("no email has been provided")
	}

	foundAccount, err := as.repo.RetrieveAccount(ctx, &receivedAccount)
	if err != nil {
		return nil, err
	}

	return foundAccount, nil
}

func (as *accountService) Create(ctx context.Context, account *types.Account) error {
	// Validate the attributes of the new Account
	if ok := validateEmailFormat(account.Email); !ok {
		return errors.New("invalid email format")
	}

	// Call the repository to create the Account
	err := as.repo.CreateAccount(ctx, account)
	if err != nil {
		return err
	}

	return nil
}

func (as *accountService) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

// validateEmailFormat validates the format of an email address using a regular expression
func validateEmailFormat(email string) bool {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	emailRegex := regexp.MustCompile(emailRegexPattern)

	return emailRegex.MatchString(email)
}
