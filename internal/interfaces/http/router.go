package http

import (
	"database/sql"
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
	"log/slog"
	"net/http"
	"time"

<<<<<<< HEAD
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
=======
	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/infrastructure/jwt"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handler/health"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/handler/root"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/otel/trace"
)

// RouterConfig contains everything required to construct the HTTP router.
type RouterConfig struct {
	DB             *sql.DB
	Logger         *slog.Logger
	JWTValidator   jwt.TokenValidator
	TracerProvider trace.TracerProvider
	CORS           middleware.CORSConfig
	RateLimiter    middleware.RateLimiterConfig
	RequestTimeout time.Duration
}

// Router owns the HTTP handler and resources that require lifecycle
// management.
type Router struct {
	handler     http.Handler
	rateLimiter *middleware.RateLimiter
}

// NewRouter constructs and validates the complete HTTP middleware stack.
//
// Configuration errors are returned immediately so that the application
// fails during startup rather than discovering invalid configuration while
// serving requests.
func NewRouter(cfg RouterConfig) (*Router, error) {
	if cfg.DB == nil {
		return nil, fmt.Errorf("http router: nil database")
	}

	if cfg.JWTValidator == nil {
		return nil, fmt.Errorf("http router: nil JWT validator")
	}

	logger := cfg.Logger
	if logger == nil {
		logger = slog.Default()
	}

	corsMiddleware, err := middleware.NewCORS(cfg.CORS)
	if err != nil {
		return nil, fmt.Errorf("http router: initialize CORS: %w", err)
	}

	rateLimiter, err := middleware.NewRateLimiter(
		cfg.RateLimiter,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"http router: initialize rate limiter: %w",
			err,
		)
	}

	timeoutMiddleware, err := middleware.NewTimeout(
		cfg.RequestTimeout,
	)
	if err != nil {
		rateLimiter.Close()

		return nil, fmt.Errorf(
			"http router: initialize request timeout: %w",
			err,
		)
	}

	// ------------------------------------------------------------
	// Public routes
	// ------------------------------------------------------------

	publicMux := http.NewServeMux()

	publicMux.HandleFunc("/", root.RootHandler)
	publicMux.HandleFunc("/health", health.Health)
	publicMux.HandleFunc("/health/live", health.Live)
	publicMux.HandleFunc("/health/ready", health.Ready(cfg.DB))
	publicMux.Handle(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)

	// ------------------------------------------------------------
	// Protected routes
	// ------------------------------------------------------------

	protectedMux := http.NewServeMux()

	// Register protected API endpoints here.
	//
	// Example:
	//
	// protectedMux.HandleFunc(
	//     "/api/v1/photos",
	//     photoHandler.List,
	// )

	protectedHandler := middleware.AuthMiddleware(
		cfg.JWTValidator,
	)(
		middleware.TenantMiddleware(
			protectedMux,
		),
	)

	// ------------------------------------------------------------
	// Root router
	// ------------------------------------------------------------

	rootMux := http.NewServeMux()

	// Everything under /api/ is protected.
	rootMux.Handle("/api/", protectedHandler)

	// Everything else is public.
	rootMux.Handle("/", publicMux)

	// ------------------------------------------------------------
	// Common middleware
	// ------------------------------------------------------------

	// Middleware executes from outside → inside.
	//
	// Recovery
	//   Request ID
	//     Tracing
	//       CORS
	//         Rate limiting
	//           Timeout
	//             Logging
	//               Router
	handler := middleware.RecoveryMiddleware(
		logger,
	)(
		middleware.RequestIDMiddleware(
			middleware.TracingMiddleware(
				cfg.TracerProvider,
			)(
				corsMiddleware(
					rateLimiter.RateLimitMiddleware(
						timeoutMiddleware(
							middleware.LoggerMiddleware(
								logger,
							)(
								rootMux,
							),
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
						),
					),
				),
			),
		),
	)
<<<<<<< HEAD
=======

	return &Router{
		handler:     handler,
		rateLimiter: rateLimiter,
	}, nil
}

// Handler returns the HTTP handler used by http.Server.
func (r *Router) Handler() http.Handler {
	return r.handler
}

// Close releases resources owned by the router.
//
// Currently this stops the rate limiter's cleanup goroutine.
func (r *Router) Close() {
	if r == nil {
		return
	}

	if r.rateLimiter != nil {
		r.rateLimiter.Close()
	}
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
}
