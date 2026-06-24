package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a custom struct type to prevent context key collisions.
// No other package can replicate it
type contextKey struct{}

// requestIDKey defines a single instance of contextKey
var requestIDKey = contextKey{}

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
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		// Call the next handler in the chain with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID extracts the request ID from the context.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
