package http

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/sql"

	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
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
	}
}

func (s *Server) Start() error {
	router := NewRouter(RouterConfig{
		DB:             s.db,
		Validator:      s.validator,
		RateLimiter:    s.rateLimiter,
		CORS:           s.corsConfig,
		RequestTimeout: s.timeout,
		Logger:         s.logger,
	})

	server := &http.Server{
		Addr:         ":" + s.cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  s.cfg.Server.ReadTimeout,
		WriteTimeout: s.cfg.Server.WriteTimeout,
		IdleTimeout:  s.cfg.Server.IdleTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Printf(
			"server running on %s\n\tAppEnv: %s\n\tReadTimeout: %v\n\tWriteTimeout: %v\n\tIdleTimeout: %v\n",
			s.cfg.Server.Port,
			s.cfg.Env,
			s.cfg.Server.ReadTimeout,
			s.cfg.Server.WriteTimeout,
			s.cfg.Server.IdleTimeout,
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	// Graceful shutdown.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case <-stop:
		log.Println("shutting down server...")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return server.Shutdown(ctx)
}
