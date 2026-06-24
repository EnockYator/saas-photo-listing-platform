package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// contextKey is a custom struct type to prevent context key collisions.
// No other package can replicate it
type authContextKey struct{}

// requestIDKey defines a single instance of contextKey
var claimsKey = authContextKey{}

// Claims is a single auth identity object
type Claims struct {
	UserID    string
	TenantID  string
	Roles     []string
}

// Token validator abstraction
type TokenValidator interface {
	Validate(token string) (*Claims, error)
}

// AuthMiddleware
//   - extracts and validates jwt token
//   - receive Claims struct (UserID, TenantID, Roles)
//   - store Claims in context
func AuthMiddleware(validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// 1. Extract Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(
					w,
					"missing authorization header",
					http.StatusUnauthorized,
				)
				return
			}

			// 2. Validate Bearer format
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(
					w,
					"invalid authorization format",
					http.StatusUnauthorized,
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

				http.Error(
					w,
					"invalid token",
					http.StatusUnauthorized,
				)
				return
			}

			// 4. Store claims in context (single source of truth)
			ctx := context.WithValue(r.Context(), claimsKey, *claims)

			// 5. OpenTelemetry enrichment
			span := trace.SpanFromContext(ctx)
			if span.SpanContext().IsValid() {
				span.SetAttributes(
					attribute.String("auth.user_id", claims.UserID),
					attribute.String("auth.tenant_id", claims.TenantID),
				)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetClaims(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsKey).(Claims)
	return c, ok
}

func GetUserID(ctx context.Context) (string, bool) {
	c, ok := GetClaims(ctx)
	if !ok {
		return "", false
	}
	return c.UserID, true
}

func GetRoles(ctx context.Context) ([]string, bool) {
	c, ok := GetClaims(ctx)
	if !ok {
		return nil, false
	}
	return c.Roles, true
}
