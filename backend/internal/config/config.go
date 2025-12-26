package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret   string
	JWTDuration time.Duration
}

// LoadEnvVar loads an environment variable by name, and returns an error if it is missing.
func LoadEnvVar(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("Missing required environment variable: %s", key)
	}
	return value, nil
}

// LoadConfig loads the configuration from environment variables (with fallback to .env file).
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	} 

	// Define the environment variables we need to load
	envVars := map[string]*string{
		"DB_HOST":     nil,
		"DB_PORT":     nil,
		"DB_USER":     nil,
		"DB_PASSWORD": nil,
		"DB_NAME":     nil,
		"DB_SSLMODE":  nil,
		"JWT_SECRET":  nil,
		"JWT_EXPIRES_IN": nil,
	}

	// Load required environment variables into the map
	for key := range envVars {
		val, err := LoadEnvVar(key)
		if err != nil {
			return nil, err
		}
		envVars[key] = &val
	}

	// Parse JWT expiration duration
	duration, err := time.ParseDuration(*envVars["JWT_EXPIRES_IN"])
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRES_IN value: %v", err)
	}

	// Populate the Config struct
	cfg := &Config{
		DBHost:     *envVars["DB_HOST"],
		DBPort:     *envVars["DB_PORT"],
		DBUser:     *envVars["DB_USER"],
		DBPassword: *envVars["DB_PASSWORD"],
		DBName:     *envVars["DB_NAME"],
		DBSSLMode:  *envVars["DB_SSLMODE"],
		JWTSecret:  *envVars["JWT_SECRET"],
		JWTDuration: duration,
	}

	return cfg, nil
}
