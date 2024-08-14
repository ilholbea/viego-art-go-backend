package config_test

import (
	"os"
	"testing"

	"github.com/ilholbea/viego-art-backend/config"
	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name           string
		env            map[string]string
		expectedConfig *config.AppConfig
		expectedError  string
	}{
		{
			name: "All environment variables set",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				"DB_SCHEMA":   "public",
				"LOG_LEVEL":   "debug",
			},
			expectedConfig: &config.AppConfig{
				LogLevel: 5, //debug level
				DatabaseConfig: &config.DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
					Database: "db",
					Schema:   "public",
				},
			},
			expectedError: "",
		},
		{
			name: "Missing DB_HOST",
			env: map[string]string{
				// "DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				"DB_SCHEMA":   "public",
				"LOG_LEVEL":   "debug",
			},
			expectedConfig: nil,
			expectedError:  "DB_HOST environment variable not set",
		},
		{
			name: "Missing DB_PORT",
			env: map[string]string{
				"DB_HOST": "localhost",
				// "DB_PORT":     "5432",
				"DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				"DB_SCHEMA":   "public",
				"LOG_LEVEL":   "debug",
			},
			expectedConfig: nil,
			expectedError:  "DB_PORT environment variable not set",
		},
		{
			name: "Missing DB_USERNAME",
			env: map[string]string{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
				// "DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				"DB_SCHEMA":   "public",
				"LOG_LEVEL":   "debug",
			},
			expectedConfig: nil,
			expectedError:  "DB_USERNAME environment variable not set",
		},
		{
			name: "Missing DB_PASSWORD",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USERNAME": "user",
				// "DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				"DB_SCHEMA":   "public",
				"LOG_LEVEL":   "debug",
			},
			expectedConfig: nil,
			expectedError:  "DB_PASSWORD environment variable not set",
		},
		{
			name: "Missing DB_DATABASE",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				// "DB_DATABASE": "db",
				"DB_SCHEMA": "public",
				"LOG_LEVEL": "debug",
			},
			expectedConfig: nil,
			expectedError:  "DB_DATABASE environment variable not set",
		},
		{
			name: "Missing DB_SCHEMA",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				// "DB_SCHEMA":   "public",
				"LOG_LEVEL": "debug",
			},
			expectedConfig: nil,
			expectedError:  "DB_SCHEMA environment variable not set",
		},
		{
			name: "Invalid DB_PORT",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "not_a_number",
				"DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				"DB_SCHEMA":   "public",
				"LOG_LEVEL":   "debug",
			},
			expectedConfig: &config.AppConfig{},
			expectedError:  "unable to parse DB_PORT",
		},
		{
			name: "Invalid LOG_LEVEL",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				"DB_SCHEMA":   "public",
				"LOG_LEVEL":   "invalid",
			},
			expectedConfig: &config.AppConfig{
				LogLevel: 4, // info level (default on invalid input)
				DatabaseConfig: &config.DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
					Database: "db",
					Schema:   "public",
				},
			},
			expectedError: "",
		},
		{
			name: "Default LOG_LEVEL",
			env: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USERNAME": "user",
				"DB_PASSWORD": "pass",
				"DB_DATABASE": "db",
				"DB_SCHEMA":   "public",
				// "LOG_LEVEL" is not set
			},
			expectedConfig: &config.AppConfig{
				LogLevel: 4, // info level (default)
				DatabaseConfig: &config.DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Username: "user",
					Password: "pass",
					Database: "db",
					Schema:   "public",
				},
			},
			expectedError: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Set up environment variables
			for key, value := range test.env {
				os.Setenv(key, value)
			}

			defer func() {
				for key := range test.env {
					os.Unsetenv(key)
				}
			}()

			got, err := config.NewConfig()
			if err != nil {
				assert.EqualError(t, err, test.expectedError)
			}
			assert.Equal(t, test.expectedConfig, got)
		})
	}
}
