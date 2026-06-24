package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// responseWriter wraps http.ResponseWriter to capture status code and response size.
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
	wrote  bool
}

// WriteHeader captures HTTP status code.
func (rw *responseWriter) WriteHeader(code int) {
	if rw.wrote {
		return
	}
	rw.status = code
	rw.wrote = true
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures response size and ensures status defaults to 200.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wrote {
		rw.WriteHeader(http.StatusOK)
	}

	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// LoggerMiddleware provides structured request logging with OpenTelemetry and request correlation.
func LoggerMiddleware(base *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()

			// Wrap response writer
			rw := &responseWriter{
				ResponseWriter: w,
				status:         0,
			}

			// Extract trace ID from OpenTelemetry span (if available)
			span := trace.SpanFromContext(r.Context())

			traceID := ""
			if span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
			}

			// Extract request ID from context
			requestID := GetRequestID(r.Context())

			// Build request-scoped logger
			logger := base.With(
				slog.String("trace_id", traceID),
				slog.String("request_id", requestID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)

			// Execute next handler in chain
			next.ServeHTTP(rw, r)

			// Compute latency
			duration := time.Since(start)

			// Final structured log
			logger.Info("http request completed",
				slog.Int("status", rw.status),
				slog.Int("size_bytes", rw.size),
				slog.Duration("duration", duration),
			)
		})
	}
}