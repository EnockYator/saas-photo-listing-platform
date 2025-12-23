package config

import (
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

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	duration, err := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "enock"),
		DBPassword: getEnv("DB_PASSWORD", "Enock02postgresql"),
		DBName:     getEnv("DB_NAME", "saas_photo_listing"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		JWTDuration: duration,
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
