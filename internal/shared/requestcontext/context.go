// Package requestcontext is a helper package that provides a way to store and retrieve
// request-scoped values in a context.Context.
package requestcontext

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
	rolesKey
)
