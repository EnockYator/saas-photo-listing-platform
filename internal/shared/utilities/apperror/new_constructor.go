package apperror

import (
	"context"

	"github.com/EnockYator/saas-photo-listing-platform/internal/domains/tenant/infrastructure/middleware"
)

// New is a Constructor that:
// 	- automatically pulls metadata from context
//	- constructs a new AppError instance
// Example usage:
	// return apperror.New(
		// 	ctx,
		// 	http.StatusNotFound,
		// 	apperror.CodeNotFound,
		// 	"album not found",
		// 	err,
	// )
func New(
	ctx context.Context,
	status int,
	code string,
	message string,
	err error,
	) *AppError {
	return &AppError{
		Status: status,
		Code: code,
		Message: message,
		Err: err,

		TraceID:   middleware.GetTraceID(ctx),
		RequestID: middleware.GetRequestID(ctx),
		UserID:    middleware.GetUserID(ctx),
		TenantID:  middleware.GetTenantID(ctx),
	}
}

// WithDetails allows builder-pattern style addition of validation failures
func (e *AppError) WithDetails(field, message string) *AppError {
	e.Details = append(e.Details, ErrorDetail{
		Field: field,
		Message: message,
	})
	return e
}