package services

import (
	"context"
	"errors"
	"github.com/ilholbea/viego-art-backend/repository"
	"github.com/ilholbea/viego-art-backend/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"
	"strconv"
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
	// Validate the attributes of the new Account
	if ok := validateEmailFormat(account.Email); !ok {
		return errors.New("invalid email format")
	}

	// Hash the password and update the Account object
	hashedPassword, err := hashPassword(account.Password)
	if err != nil {
		return err
	}
	account.Password = hashedPassword

	log.Debugf("After password update account looks like: %v", account)

	// Call the repository to create the Account
	err = as.repo.CreateAccount(ctx, account)
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

// ComparePassword compares a hashed password with its plain-text version
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil // Return true if they match, false otherwise
}
