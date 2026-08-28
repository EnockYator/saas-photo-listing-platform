package jwt_test

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"testing"
	"time"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/infrastructure/jwt"
	golangjwt "github.com/golang-jwt/jwt/v5"
)

const (
	testIssuer   = "https://auth.example.com"
	testAudience = "saas-photo-listing-api"
)

func TestNewValidator(t *testing.T) {
	privateKey := generateKey(t)

	tests := []struct {
		name    string
		config  jwt.Config
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: jwt.Config{
				PublicKey: &privateKey.PublicKey,
				Issuer:    testIssuer,
				Audience:  testAudience,
			},
		},
		{
			name: "missing public key",
			config: jwt.Config{
				Issuer:   testIssuer,
				Audience: testAudience,
			},
			wantErr: true,
		},
		{
			name: "missing issuer",
			config: jwt.Config{
				PublicKey: &privateKey.PublicKey,
				Audience:  testAudience,
			},
			wantErr: true,
		},
		{
			name: "missing audience",
			config: jwt.Config{
				PublicKey: &privateKey.PublicKey,
				Issuer:    testIssuer,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jwt.NewValidator(tt.config)

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"expected error=%v, got error=%v",
					tt.wantErr,
					err,
				)
			}
		})
	}
}

func TestValidator_Validate(t *testing.T) {
	privateKey := generateKey(t)

	validator, err := jwt.NewValidator(jwt.Config{
		PublicKey: &privateKey.PublicKey,
		Issuer:    testIssuer,
		Audience:  testAudience,
	})
	if err != nil {
		t.Fatalf("create validator: %v", err)
	}

	t.Run("accepts valid token", func(t *testing.T) {
		token := createToken(
			t,
			privateKey,
			golangjwt.SigningMethodRS256,
			"user-123",
			"tenant-123",
			[]string{"admin", "editor"},
			testIssuer,
			testAudience,
			time.Now().Add(time.Hour),
		)

		got, err := validator.Validate(token)
		if err != nil {
			t.Fatalf("expected valid token: %v", err)
		}

		if got.UserID != "user-123" {
			t.Errorf(
				"expected user ID %q, got %q",
				"user-123",
				got.UserID,
			)
		}

		if got.TenantID != "tenant-123" {
			t.Errorf(
				"expected tenant ID %q, got %q",
				"tenant-123",
				got.TenantID,
			)
		}

		if len(got.Roles) != 2 {
			t.Fatalf(
				"expected 2 roles, got %d",
				len(got.Roles),
			)
		}
	})

	t.Run("rejects empty token", func(t *testing.T) {
		_, err := validator.Validate("")

		if !errors.Is(err, jwt.ErrInvalidToken) {
			t.Fatalf(
				"expected ErrInvalidToken, got %v",
				err,
			)
		}
	})

	t.Run("rejects expired token", func(t *testing.T) {
		token := createToken(
			t,
			privateKey,
			golangjwt.SigningMethodRS256,
			"user-123",
			"tenant-123",
			[]string{"admin"},
			testIssuer,
			testAudience,
			time.Now().Add(-time.Hour),
		)

		_, err := validator.Validate(token)

		if err == nil {
			t.Fatal("expected expired token to be rejected")
		}
	})

	t.Run("rejects wrong issuer", func(t *testing.T) {
		token := createToken(
			t,
			privateKey,
			golangjwt.SigningMethodRS256,
			"user-123",
			"tenant-123",
			[]string{"admin"},
			"https://evil.example.com",
			testAudience,
			time.Now().Add(time.Hour),
		)

		_, err := validator.Validate(token)

		if err == nil {
			t.Fatal("expected wrong issuer to be rejected")
		}
	})

	t.Run("rejects wrong audience", func(t *testing.T) {
		token := createToken(
			t,
			privateKey,
			golangjwt.SigningMethodRS256,
			"user-123",
			"tenant-123",
			[]string{"admin"},
			testIssuer,
			"another-api",
			time.Now().Add(time.Hour),
		)

		_, err := validator.Validate(token)

		if err == nil {
			t.Fatal("expected wrong audience to be rejected")
		}
	})

	t.Run("rejects token without expiration", func(t *testing.T) {
		token := createTokenWithoutExpiration(
			t,
			privateKey,
			golangjwt.SigningMethodRS256,
			"user-123",
			"tenant-123",
			[]string{"admin"},
			testIssuer,
			testAudience,
		)

		_, err := validator.Validate(token)

		if err == nil {
			t.Fatal(
				"expected token without expiration to be rejected",
			)
		}
	})

	t.Run("rejects invalid signature", func(t *testing.T) {
		otherKey := generateKey(t)

		token := createToken(
			t,
			otherKey,
			golangjwt.SigningMethodRS256,
			"user-123",
			"tenant-123",
			[]string{"admin"},
			testIssuer,
			testAudience,
			time.Now().Add(time.Hour),
		)

		_, err := validator.Validate(token)

		if err == nil {
			t.Fatal(
				"expected invalid signature to be rejected",
			)
		}
	})

	t.Run("rejects missing user ID", func(t *testing.T) {
		token := createToken(
			t,
			privateKey,
			golangjwt.SigningMethodRS256,
			"",
			"tenant-123",
			[]string{"admin"},
			testIssuer,
			testAudience,
			time.Now().Add(time.Hour),
		)

		_, err := validator.Validate(token)

		if !errors.Is(err, jwt.ErrMissingUserID) {
			t.Fatalf(
				"expected ErrMissingUserID, got %v",
				err,
			)
		}
	})

	t.Run("rejects missing tenant ID", func(t *testing.T) {
		token := createToken(
			t,
			privateKey,
			golangjwt.SigningMethodRS256,
			"user-123",
			"",
			[]string{"admin"},
			testIssuer,
			testAudience,
			time.Now().Add(time.Hour),
		)

		_, err := validator.Validate(token)

		if !errors.Is(err, jwt.ErrMissingTenantID) {
			t.Fatalf(
				"expected ErrMissingTenantID, got %v",
				err,
			)
		}
	})

	t.Run("rejects missing roles", func(t *testing.T) {
		token := createToken(
			t,
			privateKey,
			golangjwt.SigningMethodRS256,
			"user-123",
			"tenant-123",
			nil,
			testIssuer,
			testAudience,
			time.Now().Add(time.Hour),
		)

		_, err := validator.Validate(token)

		if !errors.Is(err, jwt.ErrInvalidRoles) {
			t.Fatalf(
				"expected ErrInvalidRoles, got %v",
				err,
			)
		}
	})
}

func generateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate RSA key: %v", err)
	}

	return key
}

func createToken(
	t *testing.T,
	key *rsa.PrivateKey,
	method *golangjwt.SigningMethodRSA,
	userID string,
	tenantID string,
	roles []string,
	issuer string,
	audience string,
	expiration time.Time,
) string {
	t.Helper()

	token := golangjwt.NewWithClaims(
		method,
		golangjwt.MapClaims{
			"sub":       userID,
			"tenant_id": tenantID,
			"roles":     roles,
			"iss":       issuer,
			"aud":       audience,
			"exp":       expiration.Unix(),
		},
	)

	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return signed
}

func createTokenWithoutExpiration(
	t *testing.T,
	key *rsa.PrivateKey,
	method *golangjwt.SigningMethodRSA,
	userID string,
	tenantID string,
	roles []string,
	issuer string,
	audience string,
) string {
	t.Helper()

	token := golangjwt.NewWithClaims(
		method,
		golangjwt.MapClaims{
			"sub":       userID,
			"tenant_id": tenantID,
			"roles":     roles,
			"iss":       issuer,
			"aud":       audience,
		},
	)

	signed, err := token.SignedString(key)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	return signed
}
