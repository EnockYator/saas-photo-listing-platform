package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	URL      string
	DBDriver string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type Config struct {
	Env      string
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// func getIntEnv(key string, defaultValue int) int {
// 	if value := os.Getenv(key); value != "" {
// 		if intVal, err := strconv.Atoi(value); err == nil {
// 			return intVal
// 		}
// 	}
// 	return defaultValue
// }

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func Load() *Config {
	// ONLY load dotenv in local development
	if os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "development" {
		if err := godotenv.Load(".env.development"); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}

	return &Config{
		Env: getEnv("APP_ENV", "production"),
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 20*time.Second),
		},
		Database: DatabaseConfig{
			URL: getEnv("DB_URL", ""),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "secret_key"),
			Expiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		},
	}
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("jwt secret is required")
	}

	return nil
}
