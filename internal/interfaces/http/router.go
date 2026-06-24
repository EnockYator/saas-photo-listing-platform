package http

import (
	"database/sql"
	"net/http"

	_ "github.com/EnockYator/saas-photo-listing-platform/docs"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handlers/health"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handlers/root"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// register endpoints
	mux.HandleFunc("/", root.RootHandler)
	mux.HandleFunc("/health", health.Health)
	mux.HandleFunc("/health/live", health.Live)
	mux.HandleFunc("/health/ready", health.Ready(db))
	mux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // points to the generated spec
	))

	// wrap the handler with middleware chain and return
	return middleware.RecoveryMiddleware(
		middleware.CorsMiddleware(
			middleware.RateLimiter(
				middleware.RequestIDMiddleware(
					middleware.AuthMiddleware(
						middleware.TenantMiddleware(
							middleware.TraceContextMiddleware(
								middleware.TimeoutMiddleware(
									middleware.LoggerMiddleware(
										mux
									),
								),
							),
						),
					),
				),
				
			),
		),
	)



	// return middleware.RequestID(
	// 	middleware.Logging(
	// 		middleware.Recovery(
	// 			middleware.Timeout(mux),
	// 		),
	// 	),
	// )
}
