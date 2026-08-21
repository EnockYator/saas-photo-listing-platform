package apperror

import (
	"fmt"
	"context"
	"errors"

	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
)

// ErrorDetail represents a specific detail associated with an application
// error, typically a validation failure.
//
// Example:
// {
//  "field": "email",
//  "message": "must be a valid email address"
// }
type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Message string `json:"message"`
	Code string `json:"code,omitempty"`
}


// AppError represents a known application-level error.
//
// AppError separates:
//   - application semantics (Code)
//   - safe client-facing information (Message)
//   - the underlying technical cause (Err)
//   - optional validation details (Details)
//   - request/observability metadata
type AppError struct {
	HTTPCode    ErrorCode `json:"code"`
	Message string `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
	
	// Internal error (this field should never be exposed)
	Err error `json:"-"`

	// Observability metadata
	TraceID string `json:"trace_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
}

// New creates a new application error.
//
// The context is used only to capture request/observability metadata.
// HTTP status codes must not be supplied here; those are determined
// by the HTTP transport layer.
func New(
	ctx context.Context,
	code ErrorCode,
	message string,
	err error,
) *AppError {
	return &AppError{
		HTTPCode:    code,
		Message: message,
		Err:     err,
		Details: nil,

		TraceID:   requestcontext.GetTraceID(ctx),
		RequestID: requestcontext.GetRequestID(ctx),
		UserID:    requestcontext.GetUserID(ctx),
		TenantID:  requestcontext.GetTenantID(ctx),
	}
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e == nil {
		return "<nil>"
	}

	if e.Err == nil {
		return fmt.Sprintf("%s: %s", e.HTTPCode, e.Message)
	}

	return fmt.Sprintf("%s: %s: %v", e.HTTPCode, e.Message, e.Err)
}

// Unwrap exposes the underlying cause to the standard errors package.
//
// This allows:
//   - errors.Is(err, target)
//   - errors.As(err, &target)
func (e *AppError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

// WithDetails adds a validation or field-level error detail.
func (e *AppError) WithDetails(
	field string,
	message string,
	code string,
) *AppError {
	e.Details = append(e.Details, ErrorDetail{
		Field:   field,
		Message: message,
		Code:    code,
	})

	return e
}

// Is allows errors.Is to compare AppError codes.
//
// Example:
//
//	errors.Is(err, apperror.CodeNotFound)
//
// This provides a convenient application-level comparison.
func (e *AppError) Is(target error) bool {
	targetErr, ok := target.(*AppError)
	if !ok {
		return false
	}

	return e.HTTPCode == targetErr.HTTPCode
}

// NewCode creates a sentinel application error for use with errors.Is.
func NewCode(code ErrorCode) *AppError {
	return &AppError{
		HTTPCode: code,
	}
}

// IsCode reports whether err contains an AppError with the given code.
func IsCode(err error, code ErrorCode) bool {
	var appErr *AppError

	if !errors.As(err, &appErr) {
		return false
	}

	return appErr.HTTPCode == code
}