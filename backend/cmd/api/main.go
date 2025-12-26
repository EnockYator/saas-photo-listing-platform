package main

import (
	"log"
	"os"

	// "github.com/gin-gonic/gin"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/storage/postgres"
	"github.com/EnockYator/saas-photo-listing-platform/backend/internal/http"
)

func main() {
	// Load environment
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	// defer db.Close()

	// Setup Gin router
	r := http.NewRouter(cfg, db)


	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
