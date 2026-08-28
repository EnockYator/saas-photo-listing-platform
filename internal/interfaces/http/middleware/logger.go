package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"

	"go.opentelemetry.io/otel/trace"
)

// LoggerMiddleware records one structured access log entry for each HTTP
// request.
//
// It records:
//   - request ID
//   - trace ID
//   - HTTP method
//   - request path
//   - remote address
//   - user agent
//   - response status
//   - response size
//   - request duration
//
// The middleware does not log sensitive request data such as authorization
// headers, cookies, request bodies, or JWTs.
func LoggerMiddleware(
	logger *slog.Logger,
) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

<<<<<<< HEAD
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// LoggerMiddleware provides structured request logging with OpenTelemetry
// and request-ID correlation.
//
// This middleware should be placed outside (before) Auth/Tenant/RateLimit in the chain
// so that every outcome — 401s, 429s, 500s, 200s — gets logged, not just
// requests that reach the final handler.
//
// Note the double invocation: LoggerMiddleware(base) returns a
// func(http.Handler) http.Handler, so should be called as
// middleware.LoggerMiddleware(logger)(next).
func LoggerMiddleware(base *slog.Logger) func(http.Handler) http.Handler {
=======
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Reuse an existing response recorder when an outer middleware
			// already created one. Otherwise create one here.
			rw, ok := w.(*responseRecorder)
			if !ok {
				rw = &responseRecorder{
					ResponseWriter: w,
				}
			}

			next.ServeHTTP(rw, r)

			traceID := traceIDFromRequest(r)

			requestID := requestcontext.GetRequestID(
				r.Context(),
			)

			logger.InfoContext(
				r.Context(),
				"http request completed",
				slog.String("request_id", requestID),
				slog.String("trace_id", traceID),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.Int("status", rw.Status()),
				slog.Int64("bytes", rw.BytesWritten()),
				slog.Duration(
					"duration",
					time.Since(start),
				),
			)
		})
	}
}

// traceIDFromRequest returns the trace ID associated with the request's
// active OpenTelemetry span.
func traceIDFromRequest(r *http.Request) string {
	spanContext := trace.SpanContextFromContext(
		r.Context(),
	)

	if !spanContext.IsValid() {
		return ""
	}

	return spanContext.TraceID().String()
}
