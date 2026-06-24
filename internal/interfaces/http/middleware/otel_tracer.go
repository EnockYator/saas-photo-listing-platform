package middleware

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// contextKey is a custom struct type to prevent context key collisions.
// No other package can replicate it
type traceContextKey struct{}

// traceIDKey defines a single instance of contextKey
var traceIDKey = traceContextKey{}

// TraceContextMiddleware exposes the current TraceID to:
//
// 1. Response headers (X-Trace-ID)
// 2. Context helper functions
//
// Span creation itself is handled by otelhttp.
func TraceContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		spanCtx := trace.SpanContextFromContext(r.Context())

		var traceID string
		if spanCtx.IsValid() {
			traceID = spanCtx.TraceID().String()
		}

		if traceID != "" {
			w.Header().Set("X-Trace-ID", traceID)
		}

		ctx := context.WithValue(
			r.Context(),
			traceIDKey,
			traceID,
		)

		next.ServeHTTP(
			w,
			r.WithContext(ctx),
		)
	})
}

// GetTraceID returns the current trace ID.
func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}

	spanCtx := trace.SpanContextFromContext(ctx)

	if !spanCtx.IsValid() {
		return ""
	}

	return spanCtx.TraceID().String()
}
