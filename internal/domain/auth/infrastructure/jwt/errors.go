package jwt

import (
	"errors"
)

var (
	ErrInvalidToken      = errors.New("invalid token")
	ErrInvalidClaims     = errors.New("invalid claims")
	ErrInvalidIssuer     = errors.New("invalid issuer")
	ErrInvalidAudience   = errors.New("invalid audience")
	ErrMissingUserID     = errors.New("missing user ID")
	ErrMissingTenantID   = errors.New("missing tenant ID")
	ErrInvalidRoles      = errors.New("invalid roles")
	ErrUnsupportedMethod = errors.New("unsupported signing method")
)
