package postgres

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs all Postgres migrations in order
func RunMigrations(dsn string, migrationPath string) error {
	// Ensure path has file:// prefix
	migrationURL := fmt.Sprintf("file://%s", migrationPath)

	m, err := migrate.New(migrationURL, dsn)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	// Apply all up migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	log.Printf("✅ Migrations applied successfully from %s", migrationPath)
	return nil
}
