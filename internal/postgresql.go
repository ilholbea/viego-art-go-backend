package internal

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/ilholbea/viego-art/config"
	"github.com/ilholbea/viego-art/types"
	_ "github.com/lib/pq"
)

type Postgres struct {
	database *sql.DB
}

// New creates a new connection to a database.
func New(c *config.Database) (*Postgres, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return &Postgres{database: db}, err
}

// CreateAccount creates a new account in the database. Generates UUID and hashes the password
// received from the user a second time.(the password received from UI is already hashed)
func (pg *Postgres) CreateAccount(Account *types.Account) error {
	Account.ID = generateUUID()
	password, err := HashPassword(Account.Password)
	if err != nil {
		return err
	}
	insertStatement := `INSERT INTO ACCOUNT(ID, FULL_NAME, EMAIL, HASHED_PASSWORD) VALUES($1, $2, $3, $4);`

	_, err = pg.database.Exec(insertStatement, Account.ID, Account.FullName, Account.Email, password)
	if err != nil {
		return err
	}
	return nil
}

// GetAllAccounts retrieves all existing account from the database.
func (pg *Postgres) GetAllAccounts() ([]types.Account, error) {
	var response []types.Account

	selectStatement := `SELECT * from Account`
	rows, err := pg.database.Query(selectStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedAccount types.Account
		if err = rows.Scan(&retrievedAccount.ID, &retrievedAccount.FullName, &retrievedAccount.Email, &retrievedAccount.Password, &retrievedAccount.Type); err != nil {
			return nil, err
		}
		response = append(response, retrievedAccount)
	}
	return response, nil
}

// GetAccount retrieves Account with email from the database.
func (pg *Postgres) GetAccount(email string) (*types.Account, error) {
	var response []types.Account
	selectStatement := `SELECT * from Account where email=$1`
	rows, err := pg.database.Query(selectStatement, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedAccount types.Account
		if err = rows.Scan(&retrievedAccount.ID, &retrievedAccount.FullName, &retrievedAccount.Email, &retrievedAccount.Password, &retrievedAccount.Type); err != nil {
			return nil, err
		}
		response = append(response, retrievedAccount)
	}

	if len(response) > 1 {
		return nil, fmt.Errorf("too many number of results with the email: %s", email)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("no results found with the email: %s", email)
	}

	return &response[0], nil
}

// DeleteAccount removes Account with email from database.
func (pg *Postgres) DeleteAccount(email string) error {
	found, err := pg.GetAccount(email)
	if err != nil {
		return err
	}

	if found == nil {
		return fmt.Errorf("no results found with the email: %s", email)
	}

	deleteStatement := `DELETE from Account where email=$1`
	_, err = pg.database.Exec(deleteStatement, email)

	if err != nil {
		return err
	}
	return nil
}

// UpdateAccount updates an existing Account from the database.
func (pg *Postgres) UpdateAccount(Account *types.Account) error {
	found, err := pg.GetAccount(Account.Email)
	if err != nil {
		return err
	}

	if found == nil {
		return fmt.Errorf("no results found with the id: %s", Account.ID)
	}

	checkPassword := CheckPasswordHash(Account.Password, found.Password)

	updatedPassword := found.Password
	if !checkPassword {
		updatedPassword, err = HashPassword(Account.Password)
		if err != nil {
			return err
		}
	}
	updateStatement := `UPDATE Account SET full_name=$1, hashed_password=$2`
	_, err = pg.database.Exec(updateStatement, Account.FullName, updatedPassword)
	if err != nil {
		return err
	}

	return nil
}

// Login retrieves a user from the database based on the email, checks the hash and returns true
// if data is confirmed.
func (pg *Postgres) Login(email, password string) (bool, error) {
	found, err := pg.GetAccount(email)
	if err != nil {
		return false, err
	}

	// Verifies the hash of the password
	checkPassword := CheckPasswordHash(password, found.Password)

	return checkPassword, nil
}

// generateUUID used to generate UUID stored in the database as PK[ID].
func generateUUID() string {
	generatedUUID := uuid.New()
	return generatedUUID.String()
}
