package middleware

import (
	"net/http"
	"strings"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/auth/infrastructure/jwt"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
)

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

func bearerToken(header string) (string, bool) {
	parts := strings.Fields(header)

	if len(parts) != 2 {
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
}
