package middleware

import (
	"net/http"
	"strings"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/application"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// TokenValidator is an interface for validating JWT tokens and extracting claims.
type TokenValidator interface {
	Validate(token string) (*application.Claims, error)
}

// AuthMiddleware:
//   - extracts the Authorization header
//   - validates the bearer token
//   - extracts the authenticated identity
//   - stores request-scoped identity in context
//   - enriches the active OpenTelemetry span
func AuthMiddleware(validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 1. Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(
					w,
					r,
					apperror.New(
						r.Context(),
						apperror.CodeAuthTokenMissing,
						"authorization token required",
						nil,
					),
				)
				return
			}

			// 2. Validate Bearer format
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.WriteError(
					w,
					r,
					apperror.New(
						r.Context(),
						apperror.CodeAuthTokenInvalid,
						"invalid authorization format",
						nil,
					),
				)
				return
			}

			token := parts[1]

			// 3. Validate token
			claims, err := validator.Validate(token)
			if err != nil || claims == nil {
				span := trace.SpanFromContext(r.Context())
				if span.SpanContext().IsValid() {
					span.SetAttributes(
						attribute.String(
							"auth.status",
							"invalid_token",
						),
					)
				}

				response.WriteError(
					w,
					r,
					apperror.New(
						r.Context(),
						apperror.CodeAuthTokenInvalid,
						"invalid authorization token",
						err,
					),
				)

				return
			}

			// 4. Store claims in context
			ctx := requestcontext.WithUserID(r.Context(), claims.UserID)
			ctx = requestcontext.WithTenantID(ctx, claims.TenantID)
			ctx = requestcontext.WithRoles(ctx, claims.Roles)

			// 5. OpenTelemetry enrichment
			span := trace.SpanFromContext(ctx)
			if span.SpanContext().IsValid() {
				span.SetAttributes(
					attribute.String("auth.user_id", claims.UserID),
					attribute.String("auth.tenant_id", claims.TenantID),
					attribute.String("auth.roles", strings.Join(claims.Roles, ",")),
				)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
