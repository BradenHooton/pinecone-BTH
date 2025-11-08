package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Success(t *testing.T) {
	// ARRANGE: Set environment variables
	os.Setenv("DATABASE_URL", "postgres://test")
	os.Setenv("JWT_SECRET", "test-secret")
	os.Setenv("JWT_EXPIRY_HOURS", "24")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("SERVER_HOST", "0.0.0.0")
	os.Setenv("UPLOAD_DIR", "./uploads")
	os.Setenv("MAX_UPLOAD_SIZE_MB", "5")
	os.Setenv("ALLOWED_ORIGINS", "http://localhost:5173")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("USDA_API_KEY", "test-key")
	os.Setenv("USDA_API_BASE_URL", "https://api.test.gov")

	defer func() {
		// Cleanup
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("JWT_EXPIRY_HOURS")
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("SERVER_HOST")
		os.Unsetenv("UPLOAD_DIR")
		os.Unsetenv("MAX_UPLOAD_SIZE_MB")
		os.Unsetenv("ALLOWED_ORIGINS")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("USDA_API_KEY")
		os.Unsetenv("USDA_API_BASE_URL")
	}()

	// ACT
	cfg, err := Load()

	// ASSERT
	require.NoError(t, err)
	assert.Equal(t, "postgres://test", cfg.DatabaseURL)
	assert.Equal(t, "test-secret", cfg.JWTSecret)
	assert.Equal(t, 24, cfg.JWTExpiryHours)
	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "0.0.0.0", cfg.ServerHost)
	assert.Equal(t, "./uploads", cfg.UploadDir)
	assert.Equal(t, 5, cfg.MaxUploadSizeMB)
	assert.Equal(t, "http://localhost:5173", cfg.AllowedOrigins)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "test-key", cfg.USDAAPIKey)
	assert.Equal(t, "https://api.test.gov", cfg.USDAAPIBaseURL)
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	// ARRANGE: Unset required environment variables
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")

	// ACT
	_, err := Load()

	// ASSERT
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}

func TestValidate_ValidConfig(t *testing.T) {
	// ARRANGE
	cfg := &Config{
		DatabaseURL:     "postgres://test",
		JWTSecret:       "test-secret",
		JWTExpiryHours:  24,
		ServerPort:      "8080",
		ServerHost:      "0.0.0.0",
		UploadDir:       "./uploads",
		MaxUploadSizeMB: 5,
		AllowedOrigins:  "http://localhost:5173",
		LogLevel:        "info",
		USDAAPIKey:      "test-key",
		USDAAPIBaseURL:  "https://api.test.gov",
	}

	// ACT
	err := cfg.Validate()

	// ASSERT
	assert.NoError(t, err)
}

func TestValidate_InvalidJWTExpiryHours(t *testing.T) {
	// ARRANGE
	cfg := &Config{
		DatabaseURL:     "postgres://test",
		JWTSecret:       "test-secret",
		JWTExpiryHours:  0, // Invalid
		ServerPort:      "8080",
		ServerHost:      "0.0.0.0",
		UploadDir:       "./uploads",
		MaxUploadSizeMB: 5,
		AllowedOrigins:  "http://localhost:5173",
		LogLevel:        "info",
	}

	// ACT
	err := cfg.Validate()

	// ASSERT
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_EXPIRY_HOURS")
}

func TestValidate_InvalidMaxUploadSize(t *testing.T) {
	// ARRANGE
	cfg := &Config{
		DatabaseURL:     "postgres://test",
		JWTSecret:       "test-secret",
		JWTExpiryHours:  24,
		ServerPort:      "8080",
		ServerHost:      "0.0.0.0",
		UploadDir:       "./uploads",
		MaxUploadSizeMB: 0, // Invalid
		AllowedOrigins:  "http://localhost:5173",
		LogLevel:        "info",
	}

	// ACT
	err := cfg.Validate()

	// ASSERT
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MAX_UPLOAD_SIZE_MB")
}
