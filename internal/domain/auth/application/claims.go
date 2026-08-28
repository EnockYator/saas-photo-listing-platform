package application

// Claims represents the trusted identity extracted from an authenticated
// access token.
//
// The application layer deliberately exposes only identity information
// required by downstream application code. JWT-specific implementation
// details remain inside the infrastructure layer.
type Claims struct {
	UserID   string
	TenantID string
	Roles    []string
}
