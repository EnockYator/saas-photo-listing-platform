package application

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTValidator struct {
	secret []byte
}

func parseRoles(raw any) []string {
	if raw == nil {
		return nil
	}

	switch v := raw.(type) {
	case []any:
		roles := make([]string, len(v))
		for i, r := range v {
			roles[i], _ = r.(string)
		}
		return roles
	case []string:
		return v
	default:
		return nil
	}
}

func (j *JWTValidator) Validate(token string) (*Claims, error) {

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return j.secret, nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	claimsMap, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	return &Claims{
		// sub is the standard JWT claim for “subject” (user identity) defined by RFC 7519.
		UserID:   safeString(claimsMap["sub"]),
		TenantID: safeString(claimsMap["tenant_id"]),
		Roles:    parseRoles(claimsMap["roles"]),
	}, nil
}

func safeString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
