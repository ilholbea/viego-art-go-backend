package config

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	LogLevel       log.Level
}

func NewConfig() (*AppConfig, error) {
	var ll log.Level // log level
	dbHost, ok := os.LookupEnv("DB_HOST")
	if !ok {
		return nil, errors.New("DB_HOST environment variable not set")
	}

	dbPort, ok := os.LookupEnv("DB_PORT")
	if !ok {
		return nil, errors.New("DB_PORT environment variable not set")
	}

	dbUsername, ok := os.LookupEnv("DB_USERNAME")
	if !ok {
		return nil, errors.New("DB_USERNAME environment variable not set")
	}

	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		return nil, errors.New("DB_PASSWORD environment variable not set")

	}

	dbDatabase, ok := os.LookupEnv("DB_DATABASE")
	if !ok {
		return nil, errors.New("DB_DATABASE environment variable not set")
	}

	dbSchema, ok := os.LookupEnv("DB_SCHEMA")
	if !ok {
		return nil, errors.New("DB_SCHEMA environment variable not set")
	}

	logLevel, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevel = "info"
	}

	port, err := strconv.Atoi(dbPort)
	if err != nil {
		return &AppConfig{}, fmt.Errorf("unable to parse DB_PORT")
	}

	ll, err = log.ParseLevel(logLevel)
	if err != nil {
		ll = log.InfoLevel
	}

	return &AppConfig{
		LogLevel: ll,
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
