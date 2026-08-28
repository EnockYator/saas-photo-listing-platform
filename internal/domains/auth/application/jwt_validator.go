// Package auth provides token validation implementations used by the HTTP
// middleware layer.
package application

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/middleware"
)

type jwtClaims struct {
	UserID   string   `json:"sub"`
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// JWTValidator validates HMAC-signed JWTs (HS256/384/512) against a shared
// secret and maps them into middleware.Claims.
type JWTValidator struct {
	secret []byte
}

// NewJWTValidator constructs a JWTValidator from a shared HMAC secret.
func NewJWTValidator(secret string) *JWTValidator {
	return &JWTValidator{secret: []byte(secret)}
}

// Validate implements middleware.TokenValidator.
func (v *JWTValidator) Validate(tokenString string) (*middleware.Claims, error) {
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	claims := &jwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Reject tokens signed with an unexpected algorithm — this guards
		// against algorithm-confusion attacks (e.g. "alg: none" or an
		// attacker-chosen RS256->HS256 downgrade).
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return v.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.UserID == "" {
		return nil, errors.New("token missing subject")
	}
	if claims.TenantID == "" {
		return nil, errors.New("token missing tenant_id")
	}

	return &middleware.Claims{
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
		Roles:    claims.Roles,
	}, nil
}
