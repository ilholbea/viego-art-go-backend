package config

import (
	"fmt"
	"os"
	"strconv"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Schema   string
}

type AppConfig struct {
	DatabaseConfig *DatabaseConfig
}

func NewConfig() (*AppConfig, error) {

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbSchema := os.Getenv("DB_SCHEMA")

	if dbHost == "" || dbPort == "" || dbUsername == "" || dbPassword == "" || dbDatabase == "" || dbSchema == "" {
		return &AppConfig{}, fmt.Errorf("unable to initialize app config")
	}

	port, err := strconv.Atoi(dbPort)
	if err != nil {
		return &AppConfig{}, fmt.Errorf("unable to parse db port")
	}

	return &AppConfig{
		DatabaseConfig: &DatabaseConfig{
			Host:     dbHost,
			Port:     port,
			Username: dbUsername,
			Password: dbPassword,
			Database: dbDatabase,
			Schema:   dbSchema,
		},
	}, nil
}
