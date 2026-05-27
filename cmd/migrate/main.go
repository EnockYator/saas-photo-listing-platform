package main

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/internal/infrastructure/database/postgres"
)

func main() {
	cfg := config.Load()

	migrator, err := postgres.NewMigrator(cfg.Database.URL)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 2 {
		log.Fatal("command required")
	}

	command := os.Args[1]

	switch command {

	case "up":
		if err := migrator.Up(); err != nil {
			log.Fatal(err)
		}

		log.Println("migrations applied successfully")

	case "down":
		if err := migrator.Down(); err != nil {
			log.Fatal(err)
		}

		log.Println("migrations rolled back successfully")

	case "version":
		version, dirty, err := migrator.Version()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("version=%d dirty=%v\n", version, dirty)

	case "force":
		if len(os.Args) < 3 {
			log.Fatal("version required")
		}

		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		if err := migrator.Force(version); err != nil {
			log.Fatal(err)
		}

		log.Printf("forced migration version to %d\n", version)

	default:
		slog.Error("unknown command", "command:", command)
	}
}
