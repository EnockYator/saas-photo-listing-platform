// Package response provides consistent HTTP JSON responses for the application.
package response

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
)

// WriteJSON encodes data as JSON and writes it to the HTTP response.
//
// The response is not committed until JSON encoding succeeds. This prevents
// an encoding failure from leaving the client with a partially initialized
// HTTP response.
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(body)
	if err != nil {
		return err
	}

	return nil
}

// WriteSuccess writes a successful JSON response.
func WriteResponse(w http.ResponseWriter, status int, data any) {
	err := WriteJSON(w, status, APISuccessResponse{
		Success: true,
		Data:    data,
	})
	if err != nil {
		slog.Error(
			"failed to write success response",
			"error", err,
			"status", status,
		)
	}
}

// WriteError converts an application error into a standardized HTTP error
// response.
//
// Known application errors are translated according to their application
// error code. Unknown errors are treated as internal server errors and their
// underlying details are never exposed to the client.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var appErr *apperror.AppError

	if errors.As(err, &appErr) {
		writeAppError(w, r, appErr)
		return
	}

	writeUnknownError(w, r, err)
}

// writeAppError writes a known application error as an HTTP response.
func writeAppError(
	w http.ResponseWriter,
	r *http.Request,
	err *apperror.AppError,
) {
	response := APIErrorResponse{
		Code:      err.Code,
		Message:   err.Message,
		Details:   err.Details,
		RequestID: requestcontext.GetRequestID(r.Context()),
	}

	if writeErr := WriteJSON(w, statusFromCode(err.Code), response); writeErr != nil {
		slog.ErrorContext(
			r.Context(),
			"failed to write application error response",
			"error", writeErr,
			"app_error_code", err.Code,
		)
	}
}

// writeUnknownError logs the original error and writes a generic internal
// server error response to the client.
func writeUnknownError(
	w http.ResponseWriter,
	r *http.Request,
	originalErr error,
) {
	slog.ErrorContext(
		r.Context(),
		"unhandled internal server error",
		"error", originalErr,
	)

	response := APIErrorResponse{
		Code:      apperror.CodeInternalServerError,
		Message:   "internal server error",
		RequestID: requestcontext.GetRequestID(r.Context()),
	}

	if writeErr := WriteJSON(
		w,
		http.StatusInternalServerError,
		response,
	); writeErr != nil {
		slog.ErrorContext(
			r.Context(),
			"failed to write internal server error response",
			"error", writeErr,
		)
	}
}