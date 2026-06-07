package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DBConfig
	Server ServerConfig
	JWT    JWTConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN returns the PostgreSQL connection string.
func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		db.Host, db.Port, db.User, db.Password, db.Name, db.SSLMode,
	)
}

type ServerConfig struct {
	Port int
}

// Addr returns the server address in ":port" format.
func (s ServerConfig) Addr() string {
	return fmt.Sprintf(":%d", s.Port)
}

type JWTConfig struct {
	Secret                    string
	ExpirationHours           int
	RefreshTokenExpirationDays int
}

// Load reads the .env file (if present) and populates a Config
// from environment variables. It returns an error if any required
// variable is missing or invalid.
func Load() (*Config, error) {
	_ = godotenv.Load()

	dbPort, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	serverPort, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	jwtExpHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "1"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRATION_HOURS: %w", err)
	}

	refreshDays, err := strconv.Atoi(getEnv("REFRESH_TOKEN_EXPIRATION_DAYS", "7"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRATION_DAYS: %w", err)
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	cfg := &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "auth_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: serverPort,
		},
		JWT: JWTConfig{
			Secret:                    jwtSecret,
			ExpirationHours:           jwtExpHours,
			RefreshTokenExpirationDays: refreshDays,
		},
	}

	return cfg, nil
}

// getEnv reads an environment variable or returns a fallback default.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
