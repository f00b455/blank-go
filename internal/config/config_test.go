package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cfg := Load()

	assert.NotNil(t, cfg)
	assert.NotEmpty(t, cfg.Port)
	assert.NotEmpty(t, cfg.Environment)
	assert.NotEmpty(t, cfg.Database.Host)
	assert.NotEmpty(t, cfg.Database.Port)
	assert.NotEmpty(t, cfg.Database.User)
	assert.NotEmpty(t, cfg.Database.Password)
	assert.NotEmpty(t, cfg.Database.Name)
	assert.NotEmpty(t, cfg.Database.SSLMode)
}

func TestLoad_DefaultValues(t *testing.T) {
	// Ensure no env vars are set
	os.Clearenv()

	cfg := Load()

	assert.Equal(t, "3002", cfg.Port)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "dax_user", cfg.Database.User)
	assert.Equal(t, "dax_password", cfg.Database.Password)
	assert.Equal(t, "dax_db", cfg.Database.Name)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
}

func TestLoad_EnvironmentVariables(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("PORT", "8080")
	_ = os.Setenv("ENVIRONMENT", "production")
	_ = os.Setenv("DB_HOST", "prod-db.example.com")
	_ = os.Setenv("DB_PORT", "5433")
	_ = os.Setenv("DB_USER", "prod_user")
	_ = os.Setenv("DB_PASSWORD", "prod_pass")
	_ = os.Setenv("DB_NAME", "prod_db")
	_ = os.Setenv("DB_SSLMODE", "require")

	defer os.Clearenv()

	cfg := Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, "prod-db.example.com", cfg.Database.Host)
	assert.Equal(t, "5433", cfg.Database.Port)
	assert.Equal(t, "prod_user", cfg.Database.User)
	assert.Equal(t, "prod_pass", cfg.Database.Password)
	assert.Equal(t, "prod_db", cfg.Database.Name)
	assert.Equal(t, "require", cfg.Database.SSLMode)
}

func TestDatabaseConfig_DSN(t *testing.T) {
	cfg := &DatabaseConfig{
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

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "Returns env value when set",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "custom",
			expected:     "custom",
		},
		{
			name:         "Returns default when env not set",
			key:          "NONEXISTENT_KEY",
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
