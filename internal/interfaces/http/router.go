package http

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

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
			middleware.NewTraceMiddleware(
				middleware.WithServiceName("saas-photo-listing-platform"),
			)(
				corsMiddleware(
					rateLimiter.RateLimitMiddleware(
						timeoutMiddleware(
							middleware.LoggerMiddleware(
								logger,
							)(
								rootMux,
							),
						),
					),
				),
			),
		),
	)

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
}



// // Configuring NewTraceMiddleware
// handler = middleware.NewTraceMiddleware(
//     middleware.WithServiceName("saas-photo-listing-platform"),
//     middleware.WithFilter(func(r *http.Request) bool {
// 		// Exclude health, metrics, and other infrastructure endpoints
// 		return r.URL.Path == "/health" ||
// 			r.URL.Path == "/ready" ||
// 			strings.HasPrefix(r.URL.Path, "/metrics") ||
// 			strings.HasPrefix(r.URL.Path, "/debug/pprof")
// 	}),
//     middleware.WithPublicEndpointFn(), // treat all incoming requests as root spans
// )(handler)