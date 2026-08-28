package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

// CorsConfig configures allowed origins, methods, headers, and credentials
// behavior for CorsMiddleware.
type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// CorsMiddleware returns a middleware that enforces the given CORS policy.
// Note the double invocation: CorsMiddleware(cfg) returns a
// func(http.Handler) http.Handler, so call it as
// middleware.CorsMiddleware(cfg)(next), not middleware.CorsMiddleware(next).
func CorsMiddleware(cfg CorsConfig) func(http.Handler) http.Handler {
	allowedOrigins := make(map[string]struct{}, len(cfg.AllowedOrigins))
	for _, o := range cfg.AllowedOrigins {
		allowedOrigins[o] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowedOrigins[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Add("Vary", "Origin")

					if cfg.AllowCredentials {
						w.Header().Set("Access-Control-Allow-Credentials", "true")
					}
					if len(cfg.AllowedMethods) > 0 {
						w.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
					}
					if len(cfg.AllowedHeaders) > 0 {
						w.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
					}
					if cfg.MaxAge > 0 {
						w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cfg.MaxAge))
					}
				}
			}

			// Preflight requests short-circuit here, before Auth ever runs.
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
