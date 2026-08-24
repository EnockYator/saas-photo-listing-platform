package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// RecoveryMiddleware safely recovers from panics in HTTP handlers.
//
// It:
//   - prevents server crash
//   - logs stack trace
//   - attaches panic info to OpenTelemetry spans
//   - returns HTTP 500 safely
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if rec := recover(); rec != nil {

				// Capture stack trace immediately (cheap and safe)
				stack := debug.Stack()

				// Get trace OpenTelemetry span from context
				span := trace.SpanFromContext(r.Context())

				if span.SpanContext().IsValid() {
					span.RecordError(fmt.Errorf("panic: %v", rec))
					span.SetAttributes(
						attribute.String("panic.value", fmt.Sprintf("%v", rec)),
						attribute.String("panic.stack", string(stack)),
					)
				}

				// Only write response if headers not already sent
				// (avoids "superfluous WriteHeader" issues)
				if w.Header().Get("Content-Type") == "" {
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				}

				w.WriteHeader(http.StatusInternalServerError)

				response.WriteError(
					w,
					r,
					apperror.New(
						r.Context(),
						apperror.CodeInternalServerError,
						"internal server error",
						nil,
					),
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
