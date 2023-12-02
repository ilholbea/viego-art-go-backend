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
	config   *config.Database
}

// New creates a new connection to a database.
func New(c *config.Database) (Postgres, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Database)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return Postgres{}, err
	}

	return Postgres{database: db, config: c}, err
}

// CreateAccount creates a new exercise in the database.
func (pg *Postgres) CreateAccount(Account *types.Account) error {
	Account.ID = generateUUID()
	insertStatement := `INSERT INTO ACCOUNT(ID, FULL_NAME, EMAIL, HASHED_PASSWORD) VALUES($1, $2, $3, $4);`

	_, err := pg.database.Exec(insertStatement, Account.ID, Account.FullName, Account.Email, Account.HashedPassword)
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
		if err = rows.Scan(&retrievedAccount.ID, &retrievedAccount.FullName, &retrievedAccount.Email, &retrievedAccount.HashedPassword); err != nil {
			return nil, err
		}
		response = append(response, retrievedAccount)
	}
	return response, nil
}

// GetAccount retrieves Account with id from the database.
func (pg *Postgres) GetAccount(id string) (*types.Account, error) {
	var response []types.Account
	selectStatement := `SELECT * from Account where id=$1`
	rows, err := pg.database.Query(selectStatement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedAccount types.Account
		if err = rows.Scan(&retrievedAccount.ID, &retrievedAccount.FullName, &retrievedAccount.Email, &retrievedAccount.HashedPassword); err != nil {
			return nil, err
		}
		response = append(response, retrievedAccount)
	}

	if len(response) > 1 {
		return nil, fmt.Errorf("too many number of results with the id: %s", id)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("no results found with the id: %s", id)
	}

	return &response[0], nil
}

// DeleteAccount removes Account with id from database.
func (pg *Postgres) DeleteAccount(id string) error {
	found, err := pg.GetAccount(id)
	if err != nil {
		return err
	}

	if found == nil {
		return fmt.Errorf("no results found with the id: %s", id)
	}

	deleteStatement := `DELETE from Account where id=$1`
	_, err = pg.database.Exec(deleteStatement, id)

	if err != nil {
		return err
	}
	return nil
}

// UpdateAccount updates an existing exercise from the database.
func (pg *Postgres) UpdateAccount(Account *types.Account) error {
	found, err := pg.GetAccount(Account.ID)
	if err != nil {
		return err
	}

	if found == nil {
		return fmt.Errorf("no results found with the id: %s", Account.ID)
	}

	updateStatement := `UPDATE Account SET name=$1, description=$2, type=$3, video=$4, image=$5 where id=$6`
	_, err = pg.database.Exec(updateStatement, Account.FullName, Account.Email, Account.HashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Postgres) Login(email, hashedPassword string) (*types.Account, error) {
	var response []types.Account
	selectStatement := `SELECT * from Account where EMAIL=$1 and HASHED_PASSWORD=$2`
	rows, err := pg.database.Query(selectStatement, email, hashedPassword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var retrievedAccount types.Account
		if err = rows.Scan(&retrievedAccount.ID, &retrievedAccount.FullName, &retrievedAccount.Email, &retrievedAccount.HashedPassword); err != nil {
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

// generateUUID used to generate UUID stored in the database as PK[ID].
func generateUUID() string {
	generatedUUID := uuid.New()
	return generatedUUID.String()
}
