package database

import (
	"database/sql"
	"fmt"

	"github.com/ilholbea/viego-art-backend/config"

	log "github.com/sirupsen/logrus"
)

func NewDatabase(dbConf *config.DatabaseConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConf.Host, dbConf.Port, dbConf.Username, dbConf.Password, dbConf.Database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Info("Connected to database")
	return db, nil
}
