package main

import (
	"log"

	_ "github.com/lib/pq"

	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/internal/infrastructure/database/postgres"
	httpserver "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		log.Fatal("invalid config:", err)
	}

	// initialize database
	db, err := postgres.New(cfg.Database)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	server := httpserver.NewServer(cfg, db)

	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
