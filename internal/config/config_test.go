package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expected    *Config
		description string
	}{
		{
			name:    "loads default configuration",
			envVars: map[string]string{},
			expected: &Config{
				Port:        "3002",
				Environment: "development",
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     "5432",
					User:     "dax_user",
					Password: "dax_password",
					Name:     "dax_db",
					SSLMode:  "disable",
				},
			},
			description: "should load defaults when no env vars are set",
		},
		{
			name: "loads custom configuration from environment",
			envVars: map[string]string{
				"PORT":        "8080",
				"ENVIRONMENT": "production",
				"DB_HOST":     "db.example.com",
				"DB_PORT":     "5433",
				"DB_USER":     "custom_user",
				"DB_PASSWORD": "custom_pass",
				"DB_NAME":     "custom_db",
				"DB_SSLMODE":  "require",
			},
			expected: &Config{
				Port:        "8080",
				Environment: "production",
				Database: DatabaseConfig{
					Host:     "db.example.com",
					Port:     "5433",
					User:     "custom_user",
					Password: "custom_pass",
					Name:     "custom_db",
					SSLMode:  "require",
				},
			},
			description: "should load custom values from environment variables",
		},
		{
			name: "loads partial configuration with some defaults",
			envVars: map[string]string{
				"PORT":    "9000",
				"DB_HOST": "remote-db",
			},
			expected: &Config{
				Port:        "9000",
				Environment: "development",
				Database: DatabaseConfig{
					Host:     "remote-db",
					Port:     "5432",
					User:     "dax_user",
					Password: "dax_password",
					Name:     "dax_db",
					SSLMode:  "disable",
				},
			},
			description: "should mix custom values with defaults",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment
			os.Clearenv()

			// Set test environment variables
			for key, value := range tt.envVars {
				_ = os.Setenv(key, value)
			}

			// Load configuration
			cfg := Load()

			// Assertions
			assert.Equal(t, tt.expected.Port, cfg.Port)
			assert.Equal(t, tt.expected.Environment, cfg.Environment)
			assert.Equal(t, tt.expected.Database.Host, cfg.Database.Host)
			assert.Equal(t, tt.expected.Database.Port, cfg.Database.Port)
			assert.Equal(t, tt.expected.Database.User, cfg.Database.User)
			assert.Equal(t, tt.expected.Database.Password, cfg.Database.Password)
			assert.Equal(t, tt.expected.Database.Name, cfg.Database.Name)
			assert.Equal(t, tt.expected.Database.SSLMode, cfg.Database.SSLMode)
		})
	}
}

func TestDatabaseConfig_DSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DatabaseConfig
		expected string
	}{
		{
			name: "generates correct DSN",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "testuser",
				Password: "testpass",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable",
		},
		{
			name: "generates DSN with SSL mode required",
			config: DatabaseConfig{
				Host:     "prod.db.com",
				Port:     "5433",
				User:     "produser",
				Password: "prodpass",
				Name:     "proddb",
				SSLMode:  "require",
			},
			expected: "host=prod.db.com port=5433 user=produser password=prodpass dbname=proddb sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.DSN()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "returns environment value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "returns default value when env not set",
			key:          "UNSET_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			if tt.envValue != "" {
				_ = os.Setenv(tt.key, tt.envValue)
			}
			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}
