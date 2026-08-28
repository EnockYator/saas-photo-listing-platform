package middleware

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
	"github.com/google/uuid"
)

<<<<<<< HEAD
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
=======
const requestIDHeader = "X-Request-ID"

// RequestIDMiddleware assigns a unique server-generated identifier to every
// HTTP request.
//
// The identifier is:
//   - returned in the X-Request-ID response header;
//   - stored in the request context;
//   - available to logging, error handling, and downstream application code.
//
// Client-provided request IDs are intentionally ignored so that external
// callers cannot control the application's correlation identifiers.
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.NewString()

		w.Header().Set(requestIDHeader, requestID)

		ctx := requestcontext.WithRequestID(
			r.Context(),
			requestID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
	})
}
