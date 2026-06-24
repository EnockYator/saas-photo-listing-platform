package apperror

import (
	"fmt"
)

type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Message string `json:"message"`
}

type AppError struct {
	//HTTP metadata
	Status int `json:"-"`

	// Public API error fields
	Code    string `json:"code"`
	Message string `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
	
	// Internal error (this field should never be exposed)
	Err error `json:"-"` // Hide raw system/DB errors from clients

	// Observability metadata (injected from context)
	TraceID string `json:"trace_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	UserID string `json:"user_id,omitempty"`
	TenantID string `json:"tenant_id,omitempty"`
}

// Standard error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap so that errrors.Is() and errors.As() works
func (e *AppError) Unwrap() error {
	return e.Err
}
