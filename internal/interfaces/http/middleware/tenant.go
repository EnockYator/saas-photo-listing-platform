package middleware

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// contextKey is a custom struct type to prevent context key collisions.
// No other package can replicate it
type tenantContextKey struct{}

// tenantIDKey defines a single instance of contextKey
var tenantIDKey = tenantContextKey{}

// TenantMiddleware is responsible for identity extraction thus enforcing data boundaries
//   - reads Claims which are set by AuthMiddleware
//   - extracts TenantID
//   - attaches convenience value
func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Get claims from context (source of truth)
		claims, ok := GetClaims(r.Context())
		if !ok {
			http.Error(
				w,
				"missing auth context",
				http.StatusUnauthorized,
			)
			return
		}

		tenantID := claims.TenantID

		if tenantID == "" {
			http.Error(
				w,
				"missing tenant",
				http.StatusUnauthorized,
			)
			return
		}

		// 2. Attach tenant to context (optional convenience)
		ctx := context.WithValue(r.Context(), tenantIDKey, tenantID)

		// 3. OpenTelemetry enrichment
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().IsValid() {
			span.SetAttributes(
				attribute.String("tenant.id", tenantID),
			)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetTenantID(ctx context.Context) string {
	if id, ok := ctx.Value(tenantIDKey).(string); ok {
		return id
	}
	return ""
}
