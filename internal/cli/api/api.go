package cli

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"

	"github.com/EnockYator/saas-photo-listing-platform/internal/cli"
	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/infrastructure/jwt"
	"github.com/EnockYator/saas-photo-listing-platform/internal/infrastructure/database/postgres"
	httpserver "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start HTTP API server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		if err := cfg.Validate(); err != nil {
			log.Fatal("invalid config:", err)
		}

		logger := slog.Default()

		// Read the public key file from configuration path
		pubKeyBytes, err := os.ReadFile(cfg.JWT.PublicKeyPath)
		if err != nil {
			logger.Error("failed to read public key file", slog.Any("error", err))
			os.Exit(1)
		}

		// Parse the PEM encoded block
		block, _ := pem.Decode(pubKeyBytes)
		if block == nil || block.Type != "PUBLIC KEY" {
			logger.Error("failed to decode PEM block containing public key")
			os.Exit(1)
		}

		// Parse the public key to PKIX format
		pubKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			logger.Error("failed to parse PKIX public key", slog.Any("error", err))
			os.Exit(1)
		}

		// Assert type to *rsa.PublicKey
		publicKey, ok := pubKeyInterface.(*rsa.PublicKey)
		if !ok {
			logger.Error("not an RSA public key")
			os.Exit(1)
		}

		// Pass it into your existing validator setup
		tokenValidator, err := jwt.NewValidator(jwt.Config{
			PublicKey: publicKey,
			Issuer:    cfg.JWT.Issuer,
			Audience:  cfg.JWT.Audience,
		})

		if err != nil {
			logger.Error(
				"failed to initialize JWT validator",
				slog.Any("error", err),
			)
			os.Exit(1)
		}


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

		server := httpserver.NewServer(
			cfg,
			db,
			httpserver.ServerOptions{
				Logger:         logger,
				JWTValidator:   tokenValidator,
				TracerProvider: otel.GetTracerProvider(),

				CORS: middleware.CORSConfig{
					AllowedOrigins: []string{
						"http://localhost:3000",
					},
					AllowedMethods: []string{
						http.MethodGet,
						http.MethodPost,
						http.MethodPut,
						http.MethodPatch,
						http.MethodDelete,
						http.MethodOptions,
					},
					AllowedHeaders: []string{
						"Authorization",
						"Content-Type",
						"X-Request-ID",
					},
					AllowCredentials: false,
					MaxAge:            3600,
				},

				RateLimiter: middleware.RateLimiterConfig{
					RequestsPerSecond: 10,
					Burst:              20,
					CleanupInterval:   10 * time.Minute,
					TrustProxy:        false,
				},

				RequestTimeout: 30 * time.Second,
			},
		)

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