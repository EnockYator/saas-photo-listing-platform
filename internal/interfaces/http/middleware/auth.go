package middleware

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type authContextKey struct{}

var claimsKey = authContextKey{}

// Claims is a single auth identity object.
type Claims struct {
	UserID   string
	TenantID string
	Roles    []string
}

// TokenValidator abstracts JWT/token validation so AuthMiddleware doesn't
// depend on a specific implementation (JWKS, symmetric key, etc.).
type TokenValidator interface {
	Validate(token string) (*Claims, error)
}

// AuthMiddleware extracts and validates a bearer token, storing the
// resulting Claims in the request context.
//
// Should only be mounted on routes that require authentication. Mounting it on a
// mux that also serves public routes (health checks, swagger, root) will
// 401 those routes too — router  should be splitted into a public chain and a
// protected chain.
func AuthMiddleware(validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// Validate the Authorization header format. It should be "Bearer <token>".
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "invalid authorization format", http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Validate the token using the provided TokenValidator.
			claims, err := validator.Validate(token)
			if err != nil || claims == nil {
				span := trace.SpanFromContext(r.Context())
				if span.SpanContext().IsValid() {
					span.SetAttributes(attribute.String("auth.status", "invalid_token"))
				}
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsKey, *claims)

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

// GetClaims retrieves the authenticated Claims from context.
func GetClaims(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsKey).(Claims)
	return c, ok
}

// GetUserID retrieves the authenticated user's ID from context.
func GetUserID(ctx context.Context) (string, bool) {
	c, ok := GetClaims(ctx)
	if !ok || c.UserID == "" {
		return "", false
	}
	return c.UserID, true
}

// GetRoles retrieves the authenticated user's roles from context.
func GetRoles(ctx context.Context) ([]string, bool) {
	c, ok := GetClaims(ctx)
	if !ok || len(c.Roles) == 0 {
		return nil, false
	}
	return c.Roles, true
}

// GetTenantID retrieves the authenticated user's tenant ID from context.
func GetTenantID(ctx context.Context) (string, bool) {
	c, ok := GetClaims(ctx)
	if !ok || c.TenantID == "" {
		return "", false
	}
	return c.TenantID, true
}
