package main

import (
	"log"

	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
	httpserver "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http"
)

func main() {
	cfg := config.Load()

	if err := cfg.Validate(); err != nil {
		log.Fatal("invalid config:", err)
	}

	server := httpserver.NewServer(cfg)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}