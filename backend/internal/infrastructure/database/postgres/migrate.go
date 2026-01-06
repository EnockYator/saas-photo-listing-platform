package postgres

import (
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations runs migrations for all domains
func RunMigrations(dsn string) error {
	domains := []string{
		"internal/infrastructure/database/postgres/migrations/users",
		"internal/infrastructure/database/postgres/migrations/listings",
		"internal/infrastructure/database/postgres/migrations/listing_photos",
		"internal/infrastructure/database/postgres/migrations/files",
		"internal/infrastructure/database/postgres/migrations/payments",
		"internal/infrastructure/database/postgres/migrations/subscriptions",
		"internal/infrastructure/database/postgres/migrations/plans",
		"internal/infrastructure/database/postgres/migrations/plan_limits",
		"internal/infrastructure/database/postgres/migrations/audit_logs",
		"internal/infrastructure/database/postgres/migrations/notifications",
		"internal/infrastructure/database/postgres/migrations/share_links",
		"internal/infrastructure/database/postgres/migrations/shared",
		"internal/infrastructure/database/postgres/migrations/tenants",
		"internal/infrastructure/database/postgres/migrations/tenant_storage_usage",
		"internal/infrastructure/database/postgres/migrations/tenant_users",
		"internal/infrastructure/database/postgres/migrations/usage_stats",
	}

	for _, dir := range domains {
		m, err := migrate.New("file://"+dir, dsn)
		if err != nil {
			log.Fatalf("failed to initialize migrations for %s: %v", dir, err)
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration failed for %s: %v", dir, err)
		}

		log.Printf("Migrations applied for %s", dir)
	}

	return nil
}
