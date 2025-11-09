package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Database
	DatabaseURL string

	// JWT
	JWTSecret      string
	JWTExpiryHours int
	SecureCookies  bool

	// USDA API
	USDAAPIKey     string
	USDAAPIBaseURL string

	// Server
	ServerPort string
	ServerHost string

	// File Uploads
	UploadDir       string
	MaxUploadSizeMB int

	// CORS
	AllowedOrigins string

	// Logging
	LogLevel string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Try to load .env.dev file (ignore error if it doesn't exist)
	_ = godotenv.Load(".env.dev")
	_ = godotenv.Load(".env")

	cfg := &Config{
		DatabaseURL:     os.Getenv("DATABASE_URL"),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		USDAAPIKey:      os.Getenv("USDA_API_KEY"),
		USDAAPIBaseURL:  getEnvWithDefault("USDA_API_BASE_URL", "https://api.nal.usda.gov/fdc/v1"),
		ServerPort:      getEnvWithDefault("SERVER_PORT", "8080"),
		ServerHost:      getEnvWithDefault("SERVER_HOST", "0.0.0.0"),
		UploadDir:       getEnvWithDefault("UPLOAD_DIR", "./uploads"),
		AllowedOrigins:  getEnvWithDefault("ALLOWED_ORIGINS", "http://localhost:5173"),
		LogLevel:        getEnvWithDefault("LOG_LEVEL", "info"),
	}

	// Parse integer values
	var err error
	cfg.JWTExpiryHours, err = strconv.Atoi(getEnvWithDefault("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_HOURS: %w", err)
	}

	cfg.MaxUploadSizeMB, err = strconv.Atoi(getEnvWithDefault("MAX_UPLOAD_SIZE_MB", "5"))
	if err != nil {
		return nil, fmt.Errorf("invalid MAX_UPLOAD_SIZE_MB: %w", err)
	}

	// Parse boolean values
	cfg.SecureCookies, _ = strconv.ParseBool(getEnvWithDefault("SECURE_COOKIES", "false"))

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks that all required configuration is present and valid
func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters for security")
	}

	if c.JWTExpiryHours <= 0 {
		return fmt.Errorf("JWT_EXPIRY_HOURS must be positive")
	}

	if c.MaxUploadSizeMB <= 0 {
		return fmt.Errorf("MAX_UPLOAD_SIZE_MB must be positive")
	}

	if c.ServerPort == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	return nil
}

// getEnvWithDefault returns the environment variable value or a default
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
