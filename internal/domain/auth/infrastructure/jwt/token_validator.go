package jwt

import (
	"context"
	"fmt"
	"os"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/application"
	"github.com/coreos/go-oidc/v3/oidc"
)

// Validator validates OIDC ID tokens using the provider's JWKS.
type Validator struct {
	verifier *oidc.IDTokenVerifier
}

// NewValidator creates a Validator for the given OIDC issuer and audience.
// The audience is typically the OAuth2 client ID.
func NewValidator(ctx context.Context, issuer, audience string) (*Validator, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to create OIDC provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: audience,
	})

	return &Validator{verifier: verifier}, nil
}

// ValidateToken verifies the token and extracts application claims.
func (v *Validator) ValidateToken(token string) (*application.Claims, error) {
	ctx := context.Background()
	idToken, err := v.verifier.Verify(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	var claims oidcClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to extract claims: %w", err)
	}

	// Map OIDC claims to application claims
	appClaims := &application.Claims{
		UserID:   claims.Subject,
		TenantID: claims.TenantID,
		Roles:    claims.Roles,
	}

	return appClaims, nil
}

// oidcClaims defines the standard OIDC claims we care about, plus
// optional custom claims for multi‑tenancy and roles.
type oidcClaims struct {
	Subject  string   `json:"sub"`
	Email    string   `json:"email"`
	Name     string   `json:"name"`
	UserID   string   `json:"user_id"`
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`
}

// Default validator for package-level convenience.
var defaultValidator *Validator

func init() {
	issuer := os.Getenv("OIDC_ISSUER_URL")
	audience := os.Getenv("OIDC_AUDIENCE")
	if issuer != "" && audience != "" {
		v, err := NewValidator(context.Background(), issuer, audience)
		if err == nil {
			defaultValidator = v
		}
	}
}

// ValidateToken is a package-level convenience function that uses the
// default validator configured via OIDC_ISSUER and OIDC_AUDIENCE environment
// variables. It returns an error if no default validator is available.
func ValidateToken(token string) (*application.Claims, error) {
	if defaultValidator == nil {
		return nil, fmt.Errorf("no default validator configured; set OIDC_ISSUER and OIDC_AUDIENCE")
	}
	return defaultValidator.ValidateToken(token)
}