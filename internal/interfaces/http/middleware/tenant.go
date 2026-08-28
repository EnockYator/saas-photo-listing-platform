package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

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
		if span.SpanContext().IsValid() {
			span.SetAttributes(attribute.String("tenant.id", tenantID))
		}

		next.ServeHTTP(w, r)
	})
}
