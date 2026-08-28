package middleware

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

<<<<<<< HEAD
// TenantMiddleware enforces multi-tenant data boundaries by requiring valid
// Claims (set by AuthMiddleware) with a non-empty TenantID.
//
// Must run after AuthMiddleware in the chain, and should only be mounted
// on protected routes.
func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetClaims(r.Context())
		if !ok {
			http.Error(w, "missing auth context", http.StatusUnauthorized)
			return
		}

		tenantID := claims.TenantID
		if tenantID == "" {
			http.Error(w, "missing tenant", http.StatusUnauthorized)
			return
		}

		// Attach the tenant to the span for observability.
		span := trace.SpanFromContext(r.Context())
=======
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

>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
		if span.SpanContext().IsValid() {
			span.SetAttributes(attribute.String("tenant.id", tenantID))
		}

		next.ServeHTTP(w, r)
	})
}
<<<<<<< HEAD
=======

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
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
