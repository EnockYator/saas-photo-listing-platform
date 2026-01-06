package main

import (
	"fmt"
	"log"

	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/infrastructure/database/postgres"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Build Postgres DSN
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

	// Path to  migrations folder
	migrationPath := "./internal/infrastructure/database/postgres/migrations"

	// Run migrations
	if err := postgres.RunMigrations(dbURL, migrationPath); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	log.Println("✅ All migrations applied successfully")
}
