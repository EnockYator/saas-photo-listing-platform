package response

import (
	"encoding/json"
	"net/http"

	appErrors "github.com/EnockYator/saas-photo-listing-platform/pkg/errors"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type APIResponse struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
	Meta    any        `json:"meta,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any, meta any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
		Meta:    meta,
	})
}

func WriteError(w http.ResponseWriter, status int, code string, message string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error: &ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

func BadRequest(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusBadRequest, "BAD_REQUEST", message, nil)
}

func Unauthorized(w http.ResponseWriter) {
	WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "authentication required", nil)
}

func InternalServerError(w http.ResponseWriter) {
	WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "something went wrong", nil)
}

func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	// typed domain error
	if appErr, ok := err.(*appErrors.AppError); ok {
		const (
			ErrNotFound     = "NOT_FOUND"
			ErrUnauthorized = "UNAUTHORIZED"
			ErrValidation   = "VALIDATION_ERROR"
		)

		switch appErr.Code {
		case ErrNotFound:
			WriteError(w, http.StatusNotFound, appErr.Code, appErr.Message, nil)

		case ErrUnauthorized:
			WriteError(w, http.StatusUnauthorized, appErr.Code, appErr.Message, nil)

		case ErrValidation:
			WriteError(w, http.StatusBadRequest, appErr.Code, appErr.Message, nil)

		default:
			WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error", nil)
		}
		return
	}

	// fallback unknown error
	WriteError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected error", nil)
}
