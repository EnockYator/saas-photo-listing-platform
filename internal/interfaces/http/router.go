package http

import (
	"database/sql"
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handlers/health"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handlers/root"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// register endpoints
	mux.HandleFunc("/", root.RootHandler)
	mux.HandleFunc("/health", health.Health)
	mux.HandleFunc("/health/live", health.Live)
	mux.HandleFunc("/health/ready", health.Ready(db))

	// wrap the handler with middleware chain and return
	return middleware.RequestID(
		middleware.Logging(
			middleware.Recovery(
				middleware.Timeout(mux),
			),
		),
	)
}
