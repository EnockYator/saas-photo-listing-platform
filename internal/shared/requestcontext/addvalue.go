package requestcontext

import "context"

// WithTraceID returns a new context containing the trace ID.
//
// The original context is not modified.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

// WithRequestID returns a new context containing the request ID.
//
// The original context is not modified.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// WithUserID returns a new context containing the authenticated user ID.
//
// The original context is not modified.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// WithTenantID returns a new context containing the tenant ID.
//
// The original context is not modified.
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// WithRoles returns a new context containing the roles of the authenticated user.
//
// The original context is not modified.
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, rolesKey, roles)
}