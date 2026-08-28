package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

<<<<<<< HEAD
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
=======
// RecoveryMiddleware recovers from panics that occur during HTTP request
// processing.
//
// It:
//   - prevents request panics from escaping the HTTP middleware chain;
//   - records the panic on the active OpenTelemetry span;
//   - logs the panic and stack trace;
//   - prevents panic details from being exposed to clients;
//   - writes the standard internal-server-error response when the response
//     has not already been committed.
//
// Recovery should be placed outside the middleware and handler components
// whose panics it is expected to recover.
func RecoveryMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Reuse the shared response recorder when a middleware farther
			// out in the chain has already created one.
			rw, ok := w.(*responseRecorder)
			if !ok {
				rw = &responseRecorder{
					ResponseWriter: w,
				}
			}

			defer func() {
				if rec := recover(); rec != nil {
					handleRecoveredPanic(
						logger,
						rw,
						r,
						rec,
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
					)
				}
			}()

<<<<<<< HEAD
				if rw.Header().Get("Content-Type") == "" {
					rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
				}
				rw.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintln(rw, "internal server error")
			}
		}()

		next.ServeHTTP(rw, r)
	})
=======
			next.ServeHTTP(rw, r)
		})
	}
}

// handleRecoveredPanic records, logs, and safely responds to a recovered
// panic.
//
// The original panic value is intentionally never returned to the client.
func handleRecoveredPanic(
	logger *slog.Logger,
	rw *responseRecorder,
	r *http.Request,
	rec any,
) {
	stack := debug.Stack()

	panicErr := fmt.Errorf(
		"panic recovered: %v",
		rec,
	)

	recordPanicSpan(
		r,
		panicErr,
		rec,
		stack,
	)

	logPanic(
		logger,
		r,
		rec,
		stack,
	)

	// Once HTTP response headers/body have been committed, replacing the
	// response is no longer safe or possible.
	if rw.Committed() {
		return
	}

	response.WriteError(
		rw,
		r,
		apperror.New(
			r.Context(),
			apperror.CodeInternalServerError,
			"internal server error",
			nil,
		),
	)
}

// recordPanicSpan records panic information on the active OpenTelemetry span.
func recordPanicSpan(
	r *http.Request,
	panicErr error,
	rec any,
	stack []byte,
) {
	span := trace.SpanFromContext(r.Context())

	if !span.SpanContext().IsValid() {
		return
	}

	span.RecordError(
		panicErr,
		trace.WithAttributes(
			attribute.String(
				"exception.type",
				fmt.Sprintf("%T", rec),
			),
			attribute.String(
				"exception.message",
				fmt.Sprint(rec),
			),
			attribute.String(
				"exception.stacktrace",
				string(stack),
			),
		),
	)

	span.SetStatus(
		codes.Error,
		"panic recovered",
	)

	span.SetAttributes(
		attribute.Bool(
			"http.panic.recovered",
			true,
		),
	)
}

// logPanic writes the full panic diagnostics to the structured logger.
//
// Sensitive implementation details are intentionally kept out of the HTTP
// response and are recorded only in internal observability systems.
func logPanic(
	logger *slog.Logger,
	r *http.Request,
	rec any,
	stack []byte,
) {
	logger.ErrorContext(
		r.Context(),
		"panic recovered from HTTP request",
		slog.Any("panic", rec),
		slog.String("stack_trace", string(stack)),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
}
