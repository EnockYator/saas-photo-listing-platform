package requestcontext

import "context"

// GetTraceID returns the trace ID stored in ctx.
//
// It returns an empty string when no trace ID is present.
func GetTraceID(ctx context.Context) string {
	return getString(ctx, traceIDKey)
}

// GetRequestID returns the request ID stored in ctx.
//
// It returns an empty string when no request ID is present.
func GetRequestID(ctx context.Context) string {
	return getString(ctx, requestIDKey)
}

// GetUserID returns the authenticated user ID stored in ctx.
//
// It returns an empty string when no user ID is present.
func GetUserID(ctx context.Context) string {
	return getString(ctx, userIDKey)
}

// GetTenantID returns the tenant ID stored in ctx.
//
// It returns an empty string when no tenant ID is present.
func GetTenantID(ctx context.Context) string {
	return getString(ctx, tenantIDKey)
}

// GetRoles returns the roles of the authenticated user stored in ctx.
//
// It returns nil when no roles are present.
func GetRoles(ctx context.Context) []string {
	if ctx == nil {
		return nil
	}

	value, ok := ctx.Value(rolesKey).([]string)
	if !ok {
		return nil
	}

	return value
}

// getString retrieves a string value from ctx.
//
// A missing value or an unexpected value type results in an empty string.
func getString(ctx context.Context, k contextKey) string {
	if ctx == nil {
		return ""
	}

	value, ok := ctx.Value(k).(string)
	if !ok {
		return ""
	}

	return value
}
