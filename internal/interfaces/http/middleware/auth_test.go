package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/application"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
)

type fakeTokenValidator struct {
	claims *application.Claims
	err    error
	token  string
}

func (f *fakeTokenValidator) Validate(
	token string,
) (*application.Claims, error) {
	f.token = token

	if f.err != nil {
		return nil, f.err
	}

	return f.claims, nil
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("rejects missing authorization header", func(t *testing.T) {
		validator := &fakeTokenValidator{}

		nextCalled := false

		next := http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			nextCalled = true
		})

		req := httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

		rec := httptest.NewRecorder()

		AuthMiddleware(validator)(next).ServeHTTP(rec, req)

		if nextCalled {
			t.Fatal("expected next handler not to be called")
		}

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusUnauthorized,
				rec.Code,
			)
		}
	})

	t.Run("rejects malformed authorization header", func(t *testing.T) {
		tests := []string{
			"Basic token",
			"Bearer",
			"Bearer ",
			"token",
			"Bearer token extra",
		}

		for _, header := range tests {
			t.Run(header, func(t *testing.T) {
				validator := &fakeTokenValidator{}

				nextCalled := false

				next := http.HandlerFunc(func(
					w http.ResponseWriter,
					r *http.Request,
				) {
					nextCalled = true
				})

				req := httptest.NewRequest(
					http.MethodGet,
					"/protected",
					nil,
				)

				req.Header.Set(
					authorizationHeader,
					header,
				)

				rec := httptest.NewRecorder()

				AuthMiddleware(validator)(next).ServeHTTP(
					rec,
					req,
				)

				if nextCalled {
					t.Fatal(
						"expected next handler not to be called",
					)
				}

				if rec.Code != http.StatusUnauthorized {
					t.Fatalf(
						"expected status %d, got %d",
						http.StatusUnauthorized,
						rec.Code,
					)
				}
			})
		}
	})

	t.Run("rejects invalid token", func(t *testing.T) {
		validator := &fakeTokenValidator{
			err: errors.New("invalid token"),
		}

		nextCalled := false

		next := http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			nextCalled = true
		})

		req := httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

		req.Header.Set(
			authorizationHeader,
			"Bearer invalid-token",
		)

		rec := httptest.NewRecorder()

		AuthMiddleware(validator)(next).ServeHTTP(rec, req)

		if nextCalled {
			t.Fatal("expected next handler not to be called")
		}

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusUnauthorized,
				rec.Code,
			)
		}

		if validator.token != "invalid-token" {
			t.Fatalf(
				"expected validator to receive token %q, got %q",
				"invalid-token",
				validator.token,
			)
		}
	})

	t.Run("accepts valid token and stores claims in context", func(t *testing.T) {
		validator := &fakeTokenValidator{
			claims: &application.Claims{
				UserID:   "user-123",
				TenantID: "tenant-456",
				Roles:    []string{"admin", "editor"},
			},
		}

		nextCalled := false

		next := http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			nextCalled = true

			if got := requestcontext.GetUserID(r.Context()); got != "user-123" {
				t.Errorf(
					"expected user ID %q, got %q",
					"user-123",
					got,
				)
			}

			if got := requestcontext.GetTenantID(r.Context()); got != "tenant-456" {
				t.Errorf(
					"expected tenant ID %q, got %q",
					"tenant-456",
					got,
				)
			}

			roles := requestcontext.GetRoles(r.Context())

			if len(roles) != 2 {
				t.Fatalf(
					"expected 2 roles, got %d",
					len(roles),
				)
			}

			if roles[0] != "admin" || roles[1] != "editor" {
				t.Fatalf(
					"unexpected roles: %#v",
					roles,
				)
			}
		})

		req := httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

		req.Header.Set(
			authorizationHeader,
			"Bearer valid-token",
		)

		rec := httptest.NewRecorder()

		AuthMiddleware(validator)(next).ServeHTTP(rec, req)

		if !nextCalled {
			t.Fatal("expected next handler to be called")
		}
	})

	t.Run("accepts case insensitive bearer scheme", func(t *testing.T) {
		validator := &fakeTokenValidator{
			claims: &application.Claims{
				UserID:   "user-123",
				TenantID: "tenant-456",
				Roles:    []string{},
			},
		}

		next := http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
		})

		req := httptest.NewRequest(
			http.MethodGet,
			"/protected",
			nil,
		)

		req.Header.Set(
			authorizationHeader,
			"bearer valid-token",
		)

		rec := httptest.NewRecorder()

		AuthMiddleware(validator)(next).ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				rec.Code,
			)
		}
	})
}

func TestBearerToken(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
		ok     bool
	}{
		{
			name:   "valid bearer token",
			header: "Bearer abc123",
			want:   "abc123",
			ok:     true,
		},
		{
			name:   "lowercase bearer",
			header: "bearer abc123",
			want:   "abc123",
			ok:     true,
		},
		{
			name:   "mixed case bearer",
			header: "BeArEr abc123",
			want:   "abc123",
			ok:     true,
		},
		{
			name:   "missing header",
			header: "",
			ok:     false,
		},
		{
			name:   "missing token",
			header: "Bearer",
			ok:     false,
		},
		{
			name:   "basic authentication",
			header: "Basic abc123",
			ok:     false,
		},
		{
			name:   "too many fields",
			header: "Bearer abc123 extra",
			ok:     false,
		},
		{
			name:   "token without scheme",
			header: "abc123",
			ok:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := bearerToken(tt.header)

			if ok != tt.ok {
				t.Fatalf(
					"expected ok=%v, got %v",
					tt.ok,
					ok,
				)
			}

			if got != tt.want {
				t.Fatalf(
					"expected token %q, got %q",
					tt.want,
					got,
				)
			}
		})
	}
}
