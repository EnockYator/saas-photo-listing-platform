package response

import (
	"encoding/json"
	"errors"
	"net/http"

	appErrors "github.com/EnockYator/saas-photo-listing-platform/internal/shared/utilities/apperror"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type APIResponse struct {
	Success bool       `json:"success"`
	Data    any        `json:"data,omitempty"`
	Error   *appErrors.AppError `json:"error,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    data,
	}); err != nil {
		http.Error(
			w,
			"failed to encode response",
			http.StatusInternalServerError,
		)
	}
}

func WriteAppError(
	w http.ResponseWriter,
	appErr *appErrors.AppError,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Status)

	_ = json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   appErr,
	})
}

func HandleError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	var appErr *appErrors.AppError

	if errors.As(err, &appErr) {
		WriteAppError(w, appErr)
		return
	}

	WriteAppError(w, &appErrors.AppError{
		Status:  http.StatusInternalServerError,
		Code:    appErrors.CodeInternal,
		Message: "internal server error",
	})
}