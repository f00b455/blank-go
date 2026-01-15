package database

import (
	"testing"

	"github.com/f00b455/blank-go/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestConnect_InvalidDSN(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "invalid-host-that-does-not-exist",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	db, err := Connect(cfg)

	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "failed to connect to database")
}

func TestConnect_EmptyHost(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	db, err := Connect(cfg)

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestConnect_InvalidPort(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     "0",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
		SSLMode:  "disable",
	}

	db, err := Connect(cfg)

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestConnect_DSNGeneration(t *testing.T) {
	tests := []struct {
		name   string
		config *config.DatabaseConfig
	}{
		{
			name: "Standard configuration",
			config: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "testuser",
				Password: "testpass",
				Name:     "testdb",
				SSLMode:  "disable",
			},
		},
		{
			name: "With SSL mode require",
			config: &config.DatabaseConfig{
				Host:     "db.example.com",
				Port:     "5433",
				User:     "admin",
				Password: "adminpass",
				Name:     "proddb",
				SSLMode:  "require",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := tt.config.DSN()
			assert.NotEmpty(t, dsn)
			assert.Contains(t, dsn, tt.config.Host)
			assert.Contains(t, dsn, tt.config.User)
			assert.Contains(t, dsn, tt.config.Name)
			assert.Contains(t, dsn, tt.config.SSLMode)
		})
	}
}

func TestConnect_ConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    *config.DatabaseConfig
		shouldErr bool
	}{
		{
			name: "Valid config with all fields",
			config: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "user",
				Password: "pass",
				Name:     "db",
				SSLMode:  "disable",
			},
			shouldErr: true, // Will error because DB doesn't exist
		},
		{
			name: "Config with custom port",
			config: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     "5433",
				User:     "user",
				Password: "pass",
				Name:     "db",
				SSLMode:  "disable",
			},
			shouldErr: true, // Will error because DB doesn't exist
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Connect(tt.config)

			if tt.shouldErr {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}

func TestConnect_NilConfig(t *testing.T) {
	// Test that passing nil config causes a panic or error
	defer func() {
		if r := recover(); r != nil {
			// Expected behavior - function should panic with nil config
			assert.NotNil(t, r)
		}
	}()

	db, err := Connect(nil)

	// If we get here without panic, we should have an error
	if err == nil {
		t.Fatal("Expected error or panic with nil config")
	}
	assert.Error(t, err)
	assert.Nil(t, db)
}
