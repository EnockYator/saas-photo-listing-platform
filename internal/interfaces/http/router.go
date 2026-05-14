package http

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handlers"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// endpoints
	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/health/live", handlers.Live)
	// mux.HandleFunc("/health/ready", handler.Ready(db))

	// middleware chain
	handler := middleware.RequestID(
		middleware.Logging(
			middleware.Recovery(
				middleware.Timeout(mux),
			),
		),
	)

	return handler
}
