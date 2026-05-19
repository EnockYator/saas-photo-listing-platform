package http

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EnockYator/saas-photo-listing-platform/internal/config"
)

type Server struct {
	cfg *config.Config
	db *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	return &Server{
		cfg: cfg,
		db: db,
	}
}

func (s *Server) Start() error {
	router := NewRouter(s.db)

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

	// graceful shutdown
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
