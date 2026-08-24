package middleware

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TenantMiddleware is responsible for identity extraction thus enforcing data boundaries
//   - extracts TenantID from request context (populated by AuthMiddleware)
//   - attaches TenantID to request context for downstream handlers
//   - enriches OpenTelemetry spans with tenant information
func TenantMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Get tenant ID from request context (populated by AuthMiddleware)
		tenantID := requestcontext.GetTenantID(r.Context())

		if tenantID == "" {
			response.WriteError(
				w,
				r,
				apperror.New(
					r.Context(),
					apperror.CodeTenantIDMissing,
					"tenant ID missing in request context",
					nil,
				),
			)
			return
		}

		// 2. Attach tenant to context (optional convenience)
		ctx := requestcontext.WithTenantID(r.Context(), tenantID)

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
