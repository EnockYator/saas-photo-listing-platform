package middleware

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
	"github.com/google/uuid"
)

// RequestIDMiddleware generates or extracts an X-Request-ID from request header and attaches it to the context.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request already has an ID (e.g., passed from a reverse proxy)
		requestID := r.Header.Get("X-Request-ID")

		// If not, generate a new UUID
		if requestID == "" {
			requestID = uuid.NewString()
		}

		// Add the request ID to the response headers
		w.Header().Set("X-Request-ID", requestID)

		// Inject the request ID into the request context
		ctx := requestcontext.WithRequestID(r.Context(), requestID)

		// Call the next handler in the chain with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
