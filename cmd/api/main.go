package main

import (
	"log"
	"net/http"
	"fmt"
	"os"

	"github.com/EnockYator/saas-photo-listing-platform/internal/infrastructure/interfaces/http/handlers"
	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database using GORM

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()

	// Set up routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler.CheckHealth)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("SaaS Photo Listing Platform API - Version 1.0.0\n"))
	})

	// Configure server
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	
	// Start server
	log.Printf("Server running on port %v", cfg.Server.Port)
	log.Printf("Health check: http://localhost:%v/health", cfg.Server.Port)
	
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "Could not start server: %v", err)
		os.Exit(1)
	}
}