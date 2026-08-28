package http

import (
	"context"
<<<<<<< HEAD
	"log"
=======
	"database/sql"
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/sql"

	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
<<<<<<< HEAD
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"
)

// Server owns everything needed to build and run the HTTP layer: config,
// the DB handle, and the constructed middleware dependencies (rate limiter,
// validator, CORS policy, structured logger).
type Server struct {
	cfg         *config.Config
	db          *sql.DB
	validator   middleware.TokenValidator
	rateLimiter *middleware.RateLimiter
	corsConfig  middleware.CorsConfig
	timeout     time.Duration
	logger      *slog.Logger
}

// NewServer constructs a Server
func NewServer(cfg *config.Config, db *sql.DB, validator middleware.TokenValidator) *Server {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// 10 req/s sustained, burst of 20 per client key.
	rl := middleware.NewRateLimiter(10, 20,
		// Only set true if this service sits exclusively behind a proxy/LB
		// you control that overwrites (not appends to) X-Forwarded-For.
		middleware.WithTrustProxy(false),
	)
	rl.StartCleanup(5 * time.Minute)

	corsCfg := middleware.CorsConfig{
		// Should be replaced with your real frontend origin(s), ideally sourced from config/env.
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           600,
	}

	return &Server{
		cfg:         cfg,
		db:          db,
		validator:   validator,
		rateLimiter: rl,
		corsConfig:  corsCfg,
		timeout:     10 * time.Second,
		logger:      logger,
=======
	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/infrastructure/jwt"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"

	"go.opentelemetry.io/otel/trace"
)

// Server owns the HTTP server and its dependencies.
type Server struct {
	cfg *config.Config
	db  *sql.DB

	logger         *slog.Logger
	jwtValidator   jwt.TokenValidator
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
	JWTValidator   jwt.TokenValidator
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

	return &Server{
		cfg: cfg,
		db:  db,

		logger:         logger,
		jwtValidator:   opts.JWTValidator,
		tracerProvider: opts.TracerProvider,
		cors:           opts.CORS,
		rateLimiter:    opts.RateLimiter,
		requestTimeout: opts.RequestTimeout,
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
	}
}

// Start starts the HTTP server and waits for either a server error or an
// operating-system shutdown signal.
func (s *Server) Start() error {
<<<<<<< HEAD
	router := NewRouter(RouterConfig{
		DB:             s.db,
		Validator:      s.validator,
		RateLimiter:    s.rateLimiter,
		CORS:           s.corsConfig,
		RequestTimeout: s.timeout,
		Logger:         s.logger,
	})
=======
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
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312

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

<<<<<<< HEAD
	// Graceful shutdown.
=======
	// ------------------------------------------------------------
	// Graceful shutdown
	// ------------------------------------------------------------

>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
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
<<<<<<< HEAD
	case <-stop:
		log.Println("shutting down server...")
=======

	case sig := <-stop:
		s.logger.Info(
			"shutdown signal received",
			slog.String("signal", sig.String()),
		)
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
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
