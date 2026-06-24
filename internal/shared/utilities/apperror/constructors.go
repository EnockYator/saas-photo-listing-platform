package apperror

import (
	"context"
	"net/http"
)

func BadRequest(
	ctx context.Context,
	message string,
	err error,
) *AppError {
	return New(
		ctx,
		http.StatusBadRequest,
		CodeBadRequest,
		message,
		err,
	)
}

func Validation(
	ctx context.Context,
	message string,
	err error,
) *AppError {
	return New(
		ctx,
		http.StatusBadRequest,
		CodeValidation,
		message,
		err,
	)
}

func Unauthorized(
	ctx context.Context,
	message string,
	err error,
) *AppError {
	return New(
		ctx,
		http.StatusUnauthorized,
		CodeUnauthorized,
		message,
		err,
	)
}

func Forbidden(
	ctx context.Context,
	message string,
	err error,
) *AppError {
	return New(
		ctx,
		http.StatusForbidden,
		CodeForbidden,
		message,
		err,
	)
}

func NotFound(
	ctx context.Context,
	message string,
	err error,
) *AppError {
	return New(
		ctx,
		http.StatusNotFound,
		CodeNotFound,
		message,
		err,
	)
}

func Internal(
	ctx context.Context,
	message string,
	err error,
) *AppError {
	return New(
		ctx,
		http.StatusInternalServerError,
		CodeInternal,
		message,
		err,
	)
}