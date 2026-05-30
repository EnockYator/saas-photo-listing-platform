package cli

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"

	"github.com/EnockYator/saas-photo-listing-platform/internal/cli"
	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/internal/infrastructure/database/postgres"
	httpserver "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start HTTP API server",
	Run: func(cmd *cobra.Command, args []string) {

		cfg := config.Load()

		if err := cfg.Validate(); err != nil {
			log.Fatal("invalid config:", err)
		}

		// initialize database
		db, err := postgres.New(cfg.Database)
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}
		defer func() {
			if err := db.Close(); err != nil {
				log.Printf("error closing database connection: %v", err)
			}
		}()

		server := httpserver.NewServer(cfg, db)

		if err := server.Start(); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	},
}

func init() {
	cli.RootCmd.AddCommand(apiCmd)
}