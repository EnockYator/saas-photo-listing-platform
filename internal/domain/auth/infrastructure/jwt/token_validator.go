package jwt

import (
	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/application"
)

// TokenValidator validates an access token and returns the trusted
// application identity contained within it.
//
// It abstracts JWT/token validation so AuthMiddleware doesn't
// depend on a specific implementation (JWKS, symmetric key, etc.).
type TokenValidator interface {
	Validate(token string) (*application.Claims, error)
}
