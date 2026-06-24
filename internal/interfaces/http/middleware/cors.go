package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

type CorsConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

func CorsMiddleware(cfg CorsConfig) func(http.Handler) http.Handler {

	allowedOrigins := make(map[string]struct{})
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
						w.Header().Set(
							"Access-Control-Allow-Methods",
							strings.Join(cfg.AllowedMethods, ", "),
						)
					}

					if len(cfg.AllowedHeaders) > 0 {
						w.Header().Set(
							"Access-Control-Allow-Headers",
							strings.Join(cfg.AllowedHeaders, ", "),
						)
					}

					if cfg.MaxAge > 0 {
						w.Header().Set(
							"Access-Control-Max-Age",
							strconv.Itoa(cfg.MaxAge),
						)
					}
				}
			}

			// Preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
