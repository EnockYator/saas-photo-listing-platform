package application

type Claims struct {
	UserID   string
	TenantID string
	Roles    []string
}