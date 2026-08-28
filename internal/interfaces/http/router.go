package http

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/EnockYator/saas-photo-listing-platform/docs"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handlers/health"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handlers/root"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"
)

// RouterConfig bundles everything NewRouter needs to construct the
// middleware chain. They're required constructor arguments rather than defaults baked into NewRouter.
type RouterConfig struct {
	DB             *sql.DB
	Validator      middleware.TokenValidator
	RateLimiter    *middleware.RateLimiter
	CORS           middleware.CorsConfig
	RequestTimeout time.Duration
	Logger         *slog.Logger
}

// NewRouter builds the full HTTP handler: a public mux (health, root,
// swagger) and a protected mux (everything requiring auth + tenant scoping),
// each wrapped in the same outer observability/security chain, so that
// health checks and swagger docs aren't accidentally gated behind auth.
//
// Execution order (outermost first) and why:
//
//  1. Recovery       - must see every panic, including ones in the
//                       middleware below it.
//  2. RequestID       - assigns a correlation ID before anything else logs
//                       or traces.
//  3. TraceContext    - exposes the trace ID for logging/response headers.
//  4. Logger          - wraps everything after IDs exist, so 401s/429s/500s
//                       get logged, not just requests that reach a handler.
//  5. CORS            - must short-circuit preflight OPTIONS before Auth,
//                       since preflight requests never carry auth headers.
//  6. Timeout         - deadline should cover RateLimit/Auth/Tenant and the
//                       handler, including any I/O they do (DB, JWKS, etc).
//  7. RateLimiter     - reject cheaply before spending CPU/DB calls on auth.
//  8. Auth / Tenant   - only applied to the protected mux, closest to the
//                       actual handlers since they're the most
//                       request-specific.
func NewRouter(cfg RouterConfig) http.Handler {
	publicMux := http.NewServeMux()

	// register public endpoints
	publicMux.HandleFunc("/", root.RootHandler)
	publicMux.HandleFunc("/health", health.Health)
	publicMux.HandleFunc("/health/live", health.Live)
	publicMux.HandleFunc("/health/ready", health.Ready(cfg.DB))
	publicMux.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // points to the generated spec
	))

	protectedMux := http.NewServeMux()

	// register protected endpoints
	// protectedMux.HandleFunc("/api/photos", root.PhotosHandler)

	protectedChain := middleware.AuthMiddleware(cfg.Validator)(
		middleware.TenantMiddleware(protectedMux),
	)
	publicMux.Handle("/api", protectedChain)

	// Shared outer chain applies to both public and protected routes.
	return middleware.RecoveryMiddleware(
		middleware.RequestIDMiddleware(
			middleware.TraceContextMiddleware(
				middleware.LoggerMiddleware(cfg.Logger)(
					middleware.CorsMiddleware(cfg.CORS)(
						middleware.TimeoutMiddleware(cfg.RequestTimeout)(
							cfg.RateLimiter.RateLimitMiddleware(publicMux),
						),
					),
				),
			),
		),
	)
}
