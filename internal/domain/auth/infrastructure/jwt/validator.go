package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/application"
	golangjwt "github.com/golang-jwt/jwt/v5"
)

const signingAlgorithm = "RS256"

type Config struct {
	PublicKey *rsa.PublicKey
	Issuer    string
	Audience  string
}

type Validator struct {
	publicKey *rsa.PublicKey
	issuer    string
	audience  string
}

type claims struct {
	TenantID string   `json:"tenant_id"`
	Roles    []string `json:"roles"`

	golangjwt.RegisteredClaims
}

func NewValidator(cfg Config) (*Validator, error) {
	if cfg.PublicKey == nil {
		return nil, errors.New("JWT public key is required")
	}

	if cfg.Issuer == "" {
		return nil, errors.New("JWT issuer is required")
	}

	if cfg.Audience == "" {
		return nil, errors.New("JWT audience is required")
	}

	return &Validator{
		publicKey: cfg.PublicKey,
		issuer:    cfg.Issuer,
		audience:  cfg.Audience,
	}, nil
}

func (v *Validator) Validate(
	tokenString string,
) (*application.Claims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := golangjwt.ParseWithClaims(
		tokenString,
		&claims{},
		func(token *golangjwt.Token) (any, error) {
			if token.Method != golangjwt.SigningMethodRS256 {
				return nil, fmt.Errorf(
					"%w: %s",
					ErrUnsupportedMethod,
					token.Method.Alg(),
				)
			}

			return v.publicKey, nil
		},
		golangjwt.WithValidMethods([]string{
			signingAlgorithm,
		}),
		golangjwt.WithIssuer(v.issuer),
		golangjwt.WithAudience(v.audience),
		golangjwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"%w: %w",
			ErrInvalidToken,
			err,
		)
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	jwtClaims, ok := token.Claims.(*claims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	if jwtClaims.Subject == "" {
		return nil, ErrMissingUserID
	}

	if jwtClaims.TenantID == "" {
		return nil, ErrMissingTenantID
	}

	if jwtClaims.Roles == nil {
		return nil, ErrInvalidRoles
	}

	return &application.Claims{
		UserID:   jwtClaims.Subject,
		TenantID: jwtClaims.TenantID,
		Roles:    append([]string(nil), jwtClaims.Roles...),
	}, nil
}
