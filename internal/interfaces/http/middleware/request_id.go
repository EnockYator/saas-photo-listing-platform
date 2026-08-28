package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// contextKey is a custom struct type to prevent context key collisions.
// No other package can replicate it.
type requestIDContextKey struct{}

// requestIDKey defines a single instance of contextKey.
var requestIDKey = requestIDContextKey{}

// RequestIDMiddleware generates or extracts an X-Request-ID from the request
// header and attaches it to the context.
//
// A client-supplied X-Request-ID is only trusted if it parses as a valid
// UUID. This prevents arbitrary client-controlled strings from flowing into
// logs and traces (log injection, correlation-ID spoofing) if this service
// is reachable directly from untrusted clients rather than exclusively via a
// trusted reverse proxy that manages this header itself.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")

		if requestID == "" {
			requestID = uuid.NewString()
		} else if _, err := uuid.Parse(requestID); err != nil {
			requestID = uuid.NewString()
		}

		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

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
