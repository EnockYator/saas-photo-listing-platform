package middleware

import (
	"net/http"
	"strings"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/infrastructure/jwt"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
)

<<<<<<< HEAD
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
=======
const authorizationHeader = "Authorization"

func AuthMiddleware(
	validator jwt.TokenValidator,
) func(http.Handler) http.Handler {
	if validator == nil {
		panic("middleware: nil token validator")
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			token, ok := bearerToken(
				r.Header.Get(authorizationHeader),
			)

			if !ok {
				writeAuthenticationError(w, r)
				return
			}

			claims, err := validator.Validate(token)
			if err != nil {
				writeAuthenticationError(w, r)
				return
			}

			ctx := requestcontext.WithUserID(
				r.Context(),
				claims.UserID,
			)

			ctx = requestcontext.WithTenantID(
				ctx,
				claims.TenantID,
			)
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312

			ctx = requestcontext.WithRoles(
				ctx,
				claims.Roles,
			)

			next.ServeHTTP(
				w,
				r.WithContext(ctx),
			)
		})
	}
}

<<<<<<< HEAD
// GetClaims retrieves the authenticated Claims from context.
func GetClaims(ctx context.Context) (Claims, bool) {
	c, ok := ctx.Value(claimsKey).(Claims)
	return c, ok
}

// GetUserID retrieves the authenticated user's ID from context.
func GetUserID(ctx context.Context) (string, bool) {
	c, ok := GetClaims(ctx)
	if !ok || c.UserID == "" {
=======
func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)

	if len(parts) != 2 {
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
		return "", false
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}

	if parts[1] == "" {
		return "", false
	}

	return parts[1], true
}

<<<<<<< HEAD
// GetRoles retrieves the authenticated user's roles from context.
func GetRoles(ctx context.Context) ([]string, bool) {
	c, ok := GetClaims(ctx)
	if !ok || len(c.Roles) == 0 {
		return nil, false
	}
	return c.Roles, true
=======
func writeAuthenticationError(
	w http.ResponseWriter,
	r *http.Request,
) {
	response.WriteError(
		w,
		r,
		apperror.New(
			r.Context(),
			apperror.CodeUnauthorized,
			"authentication required",
			nil,
		),
	)
>>>>>>> 9f919071e0e60cf5b08cae69e04231b06af1f312
}

// GetTenantID retrieves the authenticated user's tenant ID from context.
func GetTenantID(ctx context.Context) (string, bool) {
	c, ok := GetClaims(ctx)
	if !ok || c.TenantID == "" {
		return "", false
	}
	return c.TenantID, true
}
