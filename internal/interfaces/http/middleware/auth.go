package middleware

import (
	"context"
	"net/http"
	"strings"

	auth "github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/application"
	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/infrastructure/jwt"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type authContextKey struct{}

var claimsKey = authContextKey{}

// AuthMiddleware extracts and validates a bearer token, storing the
// resulting Claims in the request context and also populating the
// requestcontext values (user_id, tenant_id, roles).
//
// Should only be mounted on routes that require authentication.
func AuthMiddleware(validator jwt.TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			span := trace.SpanFromContext(ctx)

			// 1. Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.WriteError(
					w,
					r,
					apperror.New(
						ctx,
						apperror.CodeAuthTokenMissing,
						"missing authorization header",
						nil,
					),
				)
				span.SetAttributes(
					attribute.String(
						"auth.status",
						"missing_token",
					),
				)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
				response.WriteError(
					w,
					r,
					apperror.New(
						ctx,
						apperror.CodeAuthTokenInvalid,
						"invalid authorization format",
						nil,
					),
				)
				span.SetAttributes(
					attribute.String(
						"auth.status",
						"invalid_token_format",
					),
				)
				return
			}

			token := parts[1]

			// 2. Validate token using the provided validator
			claims, err := validator.Validate(token)
			if err != nil || claims == nil {
				if span.SpanContext().IsValid() {
					span.SetAttributes(attribute.String("auth.status", "invalid_token"))
				}

				code := apperror.CodeAuthTokenInvalid
				if apperror.IsCode(err, apperror.CodeAuthTokenInvalid) {
					code = apperror.CodeAuthTokenExpired
				}
				response.WriteError(
					w,
					r,
					apperror.New(
						ctx,
						code,
						"invalid or expired token",
						nil,
					),
				)
				span.SetAttributes(
					attribute.String(
						"auth.status",
						"invalid_or_expired_token",
					),
				)
				return
			}

			// 3. Store claims in context (value copy to prevent mutation)
			ctx = context.WithValue(ctx, claimsKey, *claims)

			// 4. Also set requestcontext values for easy access
			ctx = requestcontext.WithUserID(ctx, claims.UserID)
			ctx = requestcontext.WithTenantID(ctx, claims.TenantID)
			ctx = requestcontext.WithRoles(ctx, claims.Roles)

			// 5. Add trace attributes
			span.SetAttributes(
				attribute.String("auth.user_id", claims.UserID),
				attribute.String("auth.tenant_id", claims.TenantID),
				attribute.StringSlice("auth.roles", claims.Roles),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetClaims retrieves the authenticated Claims from context.
func GetClaims(ctx context.Context) (auth.Claims, bool) {
	c, ok := ctx.Value(claimsKey).(auth.Claims)
	return c, ok
}
