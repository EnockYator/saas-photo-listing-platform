package cli

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/EnockYator/saas-photo-listing-platform/internal/cli"
	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/internal/infrastructure/database/postgres"
	"github.com/EnockYator/saas-photo-listing-platform/internal/observability/tracing"

	httpserver "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/dto/auth"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start HTTP API server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		if err := cfg.Validate(); err != nil {
			log.Fatal("invalid config:", err)
		}

		logger := slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelInfo,
				},
			),
		)

		ctx, cancel := context.WithCancel(context.Background())

		tp, err := tracing.Init(
			ctx,
			tracing.Config{
				ServiceName:    "cloud-gallery-api",
				ServiceVersion: "1.0.0",
				Environment:    "development",
				SamplingRatio:  0.10,
			},
			logger,
		)
		if err != nil {
			logger.Error(
				"failed to initialize tracing",
				"error", err,
			)
			os.Exit(1)
		}

		defer func() {
			if err := tracing.Shutdown(tp, 5*time.Second); err != nil {
				logger.Error(
					"failed to shutdown tracing",
					"error", err,
				)
			}
		}()

		db, err := postgres.New(cfg.Database)
		if err != nil {
			log.Fatal("failed to connect to database:", err)
		}

		defer func() {
			if err := db.Close(); err != nil {
				logger.Error(
					"failed to close database",
					slog.Any("error", err),
				)
			}
		}()

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Fatal("JWT_SECRET environment variable is required")
		}

		// server := httpserver.NewServer(
		// 	cfg,
		// 	db,
		// 	httpserver.ServerOptions{
		// 		Logger:         logger,
		// 		JWTValidator:   tokenValidator,
		// 		TracerProvider: otel.GetTracerProvider(),

		// 		CORS: middleware.CORSConfig{
		// 			AllowedOrigins: []string{
		// 				"http://localhost:3000",
		// 			},
		// 			AllowedMethods: []string{
		// 				http.MethodGet,
		// 				http.MethodPost,
		// 				http.MethodPut,
		// 				http.MethodPatch,
		// 				http.MethodDelete,
		// 				http.MethodOptions,
		// 			},
		// 			AllowedHeaders: []string{
		// 				"Authorization",
		// 				"Content-Type",
		// 				"X-Request-ID",
		// 			},
		// 			AllowCredentials: false,
		// 			MaxAge:           3600,
		// 		},

		// 		RateLimiter: middleware.RateLimiterConfig{
		// 			RequestsPerSecond: 10,
		// 			Burst:             20,
		// 			CleanupInterval:   10 * time.Minute,
		// 			TrustProxy:        false,
		// 		},

		// 		RequestTimeout: 30 * time.Second,
		// 	},
		// )

		server := httpserver.NewServer(cfg, db, nil)
		if err := server.Start(); err != nil {
			logger.Error(
				"server stopped with error",
				slog.Any("error", err),
			)
			os.Exit(1)
		}
	},
}

func init() {
	cli.RootCmd.AddCommand(apiCmd)
}
