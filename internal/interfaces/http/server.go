package http

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/infrastructure/jwt"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"

	"go.opentelemetry.io/otel/trace"
)

// Server owns the HTTP server and its dependencies.
type Server struct {
	cfg *config.Config
	db  *sql.DB

	logger         *slog.Logger
	jwtValidator   jwt.Validator
	tracerProvider trace.TracerProvider
	cors           middleware.CORSConfig
	rateLimiter    middleware.RateLimiterConfig
	requestTimeout time.Duration
}

// ServerOptions contains dependencies required by the HTTP server.
//
// Keeping infrastructure dependencies here avoids making the HTTP package
// responsible for constructing things such as JWT validators or tracer
// providers.
type ServerOptions struct {
	Logger         *slog.Logger
	TracerProvider trace.TracerProvider
	CORS           middleware.CORSConfig
	RateLimiter    middleware.RateLimiterConfig
	RequestTimeout time.Duration
}

// NewServer creates an HTTP server.
func NewServer(
	cfg *config.Config,
	db *sql.DB,
	opts ServerOptions,
) *Server {
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	// Create a validator for Google (or any OIDC provider)
	tokenValidator, err := jwt.NewValidator(
		context.Background(),
		"https://accounts.google.com", // OIDC issuer
		"YOUR_GOOGLE_CLIENT_ID.apps.googleusercontent.com", // audience / client ID
	)
	if err != nil {
		logger.Error("failed to validate token", err)
	}

	return &Server{
		cfg: cfg,
		db:  db,

		logger:         logger,
		jwtValidator:   tokenValidator,
		tracerProvider: opts.TracerProvider,
		cors:           opts.CORS,
		rateLimiter:    opts.RateLimiter,
		requestTimeout: opts.RequestTimeout,
	}
}

// Start starts the HTTP server and waits for either a server error or an
// operating-system shutdown signal.
func (s *Server) Start() error {
	router, err := NewRouter(RouterConfig{
		DB:             s.db,
		Logger:         s.logger,
		JWTValidator:   s.jwtValidator,
		TracerProvider: s.tracerProvider,
		CORS:           s.cors,
		RateLimiter:    s.rateLimiter,
		RequestTimeout: s.requestTimeout,
	})
	if err != nil {
		return err
	}

	defer router.Close()

	server := &http.Server{
		Addr: ":" + s.cfg.Server.Port,

		Handler: router.Handler(),

		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		s.logger.Info(
			"HTTP server starting",
			slog.String(
				"address",
				server.Addr,
			),
			slog.String(
				"environment",
				s.cfg.Env,
			),
			slog.Duration(
				"read_timeout",
				s.cfg.Server.ReadTimeout,
			),
			slog.Duration(
				"write_timeout",
				s.cfg.Server.WriteTimeout,
			),
			slog.Duration(
				"idle_timeout",
				s.cfg.Server.IdleTimeout,
			),
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// ------------------------------------------------------------
	// Graceful shutdown
	// ------------------------------------------------------------

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer signal.Stop(stop)

	select {
	case err := <-serverErr:
		return err

	case sig := <-stop:
		s.logger.Info(
			"shutdown signal received",
			slog.String("signal", sig.String()),
		)
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	s.logger.Info("HTTP server shutting down")

	if err := server.Shutdown(shutdownCtx); err != nil {
		s.logger.Error(
			"HTTP server shutdown failed",
			slog.Any("error", err),
		)

		return err
	}

	s.logger.Info("HTTP server stopped")

	return nil
}
