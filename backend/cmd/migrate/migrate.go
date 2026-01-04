package main

import (
	"log"

	"https://github.com/EnockYator/saas-photo-listing-platform/backend/internal/config"
	"https://github.com/EnockYator/saas-photo-listing-platform/backend/internal/infrastructure/database/postgres"
)

func main() {
	cfg := config.Load()

	if err := postgres.RunMigrations(
		cfg.DatabaseURL,
		"internal/infrastructure/database/postgres/migrations",
	); err != nil {
		log.Fatal(err)
	}
}
