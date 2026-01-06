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

<<<<<<< HEAD
<<<<<<< HEAD
	// Build Postgres DSN
=======
	// Build the Postgres DSN
>>>>>>> 5e02025 (chore: configure sqlc and go-migrate)
=======
	// Build Postgres DSN
>>>>>>> 90020b5 (chore(migration): configure mirations)
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)

<<<<<<< HEAD
<<<<<<< HEAD
	// Path to  migrations folder
	migrationPath := "./internal/infrastructure/database/postgres/migrations"

	// Run migrations
	if err := postgres.RunMigrations(dbURL, migrationPath); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	log.Println("✅ All migrations applied successfully")
}
<<<<<<< HEAD
=======
=======
=======
	// Path to  migrations folder
	migrationPath := "./internal/infrastructure/database/postgres/migrations"

>>>>>>> 90020b5 (chore(migration): configure mirations)
	// Run migrations
	if err := postgres.RunMigrations(dbURL, migrationPath); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	log.Println("✅ All migrations applied successfully")
}
>>>>>>> 5e02025 (chore: configure sqlc and go-migrate)
>>>>>>> b475ca3 (chore: configure sqlc and go-migrate)
