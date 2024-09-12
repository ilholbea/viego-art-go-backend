package database

import (
	"github.com/ilholbea/viego-art-backend/internal/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewDatabase(t *testing.T) {
	tests := []struct {
		name        string
		dbConfig    *config.DatabaseConfig
		expectedErr string
	}{
		{
			name: "TestNewDatabase",
			dbConfig: &config.DatabaseConfig{
				Host:     "test",
				Port:     1,
				Username: "test",
				Password: "test",
				Database: "test",
				Schema:   "test",
			},
			expectedErr: "nothing",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			_, err := NewDatabase(test.dbConfig)
			if err != nil {
				assert.EqualError(t, err, test.expectedErr)
			}
		})
	}
}
