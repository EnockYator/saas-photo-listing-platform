package cli

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/internal/infrastructure/database/postgres"

	"github.com/EnockYator/saas-photo-listing-platform/internal/cli"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Database migrations",
}

var migrateUpCmd = &cobra.Command{
	Use: "up",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrate("up", "")
	},
}

var migrateDownCmd = &cobra.Command{
	Use: "down",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrate("down", "")
	},
}

var migrateVersionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		runMigrate("version", "")
	},
}

var migrateForceCmd = &cobra.Command{
	Use:  "force [version]",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runMigrate("force", args[0])
	},
}

func runMigrate(command string, arg string) {
	cfg := config.Load()

	migrator, err := postgres.NewMigrator(cfg.Database.URL)
	if err != nil {
		slog.Error("Failed to connect to migrator", "error", err)
		os.Exit(1)
	}

	switch command {

	case "up":
		err := migrator.Up()
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("No new migrations to apply")
			return
		}
		if err != nil {
			slog.Error("Failed to apply migrations", "error", err)
		}
		slog.Info("✔ migrations completed", "status", "applied")

	case "down":
		if err := migrator.Down(); err != nil {
			slog.Error("Failed to drop migrations", "error", err)
		}
		slog.Info("✔ migrations completed", "status", "dropped")

	case "version":
		v, dirty, err := migrator.Version()
		if err != nil {
			slog.Error("Failed to get migration version", "error", err)
		}
		fmt.Printf("migration version:%d dirty=%v\n", v, dirty)

	case "force":
		version, err := strconv.Atoi(arg)
		if err != nil {
			slog.Error("Failed to convert value to integer", "error", err)
		}
		if err := migrator.Force(version); err != nil {
			slog.Error("Failed to force migration version", "error", err)
		}
		fmt.Printf("migration forced to version %d\n", version)
	}
}

func init() {
	migrateCmd.AddCommand(migrateUpCmd)
	migrateCmd.AddCommand(migrateDownCmd)
	migrateCmd.AddCommand(migrateVersionCmd)
	migrateCmd.AddCommand(migrateForceCmd)

	cli.RootCmd.AddCommand(migrateCmd)
}
