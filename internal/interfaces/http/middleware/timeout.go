package middleware

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type timeoutContextKey string

const timeoutKey timeoutContextKey = "requestTimeout"

// TimeoutMiddleware applies a request deadline using context cancellation.
//
// It does NOT terminate handlers.
// It relies on downstream operations respecting context cancellation:
//
//   - db.QueryContext()
//   - req.WithContext()
//   - gRPC calls using context
//   - any select on ctx.Done()
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Preserve shorter upstream deadlines.
			if deadline, ok := r.Context().Deadline(); ok {
				if time.Until(deadline) <= timeout {
					next.ServeHTTP(w, r)
					return
				}
			}

			ctx, cancel := context.WithTimeout(
				r.Context(),
				timeout,
			)
			defer cancel()

			// Store timeout value in context for
			// logging, diagnostics, or error handling.
			ctx = context.WithValue(
				ctx,
				timeoutKey,
				timeout,
			)

			span := trace.SpanFromContext(ctx)

			if span.SpanContext().IsValid() {
				span.SetAttributes(
					attribute.String(
						"http.timeout",
						timeout.String(),
					),
				)
			}

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

// GetTimeout returns the configured request timeout.
func GetTimeout(ctx context.Context) (time.Duration, bool) {
	timeout, ok := ctx.Value(timeoutKey).(time.Duration)
	return timeout, ok
}
