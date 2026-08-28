package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// recoveryWriter wraps http.ResponseWriter to make it safe to call
// WriteHeader from the recovery handler even if the wrapped handler already
// started writing a response before it panicked. Recovery sits outermost in
// the chain, so it can't see whether an inner middleware (e.g. Logger) has
// already tracked a write; this local guard makes that irrelevant.
type recoveryWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (rw *recoveryWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *recoveryWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// RecoveryMiddleware safely recovers from panics in HTTP handlers.
//
// It:
//   - prevents the server from crashing
//   - logs the stack trace onto the active span
//   - attaches panic info to OpenTelemetry spans
//   - returns HTTP 500 safely, without a double WriteHeader panic
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &recoveryWriter{ResponseWriter: w}

		defer func() {
			if rec := recover(); rec != nil {
				stack := debug.Stack()

				span := trace.SpanFromContext(r.Context())
				if span.SpanContext().IsValid() {
					span.RecordError(fmt.Errorf("panic: %v", rec))
					span.SetAttributes(
						attribute.String("panic.value", fmt.Sprintf("%v", rec)),
						attribute.String("panic.stack", string(stack)),
					)
				}

				if rw.Header().Get("Content-Type") == "" {
					rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
				}
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintln(rw, "internal server error")
			}
		}()

		next.ServeHTTP(rw, r)
	})
}
