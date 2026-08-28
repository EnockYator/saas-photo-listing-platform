package middleware

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TenantMiddleware enforces the presence of an authenticated tenant identity
// in the request context.
//
// The tenant ID must have been established by AuthMiddleware from validated
// authentication claims. This middleware never accepts tenant identity from
// client-controlled HTTP input and never replaces the authenticated tenant
// value.
func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := requestcontext.GetTenantID(r.Context())

		if tenantID == "" {
			writeTenantError(w, r)
			return
		}

		span := trace.SpanFromContext(r.Context())

		if span.SpanContext().IsValid() {
			span.SetAttributes(
				attribute.String("tenant.id", tenantID),
			)
		}

		next.ServeHTTP(w, r)
	})
}

// writeTenantError returns a generic error without exposing authentication
// implementation details.
func writeTenantError(
	w http.ResponseWriter,
	r *http.Request,
) {
	response.WriteError(
		w,
		r,
		apperror.New(
			r.Context(),
			apperror.CodeTenantIDMissing,
			"tenant context required",
			nil,
		),
	)
}
