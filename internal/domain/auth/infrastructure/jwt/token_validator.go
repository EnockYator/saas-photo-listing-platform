package jwt

import (
	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/application"
)

// TokenValidator validates an access token and returns the trusted
// application identity contained within it.
type TokenValidator interface {
	Validate(token string) (*application.Claims, error)
}
