package requestcontext

import (
	"context"
)

// contextKey is an unexported type used for context keys.
//
// Using a private/custom type prevents collisions with keys defined by other
// packages, even when those packages use the same underlying type.
type contextKey int

const (
	traceIDKey contextKey = iota
	requestIDKey
	userIDKey
	tenantIDKey
)

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