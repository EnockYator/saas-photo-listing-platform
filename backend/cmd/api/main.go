package main

import (
	"log"
	"os"

	// "github.com/gin-gonic/gin"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/infrastructure/database/postgres"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/interfaces/http"
)

func main() {
	// Load environment
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database using GORM
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Initialize the HTTP Gin router and pass the GORM DB
	r := http.NewRouter(cfg, db)

	// Determine the port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
