package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/f00b455/blank-go/internal/config"
)

func TestConnect(t *testing.T) {
	tests := []struct {
		name        string
		config      *config.DatabaseConfig
		expectError bool
		description string
	}{
		{
			name: "fails with invalid host",
			config: &config.DatabaseConfig{
				Host:     "invalid-host-that-does-not-exist",
				Port:     "5432",
				User:     "testuser",
				Password: "testpass",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			expectError: true,
			description: "should fail when host is invalid",
		},
		{
			name: "fails with invalid port",
			config: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     "99999",
				User:     "testuser",
				Password: "testpass",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			expectError: true,
			description: "should fail when port is invalid",
		},
		{
			name: "fails with empty credentials",
			config: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "",
				Password: "",
				Name:     "",
				SSLMode:  "disable",
			},
			expectError: true,
			description: "should fail when credentials are empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.config)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}

func TestConnect_DSN(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	dsn := cfg.DSN()
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "user=testuser")
	assert.Contains(t, dsn, "password=testpass")
	assert.Contains(t, dsn, "dbname=testdb")
	assert.Contains(t, dsn, "sslmode=disable")
}
