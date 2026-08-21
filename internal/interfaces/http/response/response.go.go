package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/requestcontext"
)

func WriteJSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(APISuccessResponse{
		Success: true,
		Data:    data,
	}); err != nil {
		http.Error(
			w,
			"failed to encode response",
			statusFromCode(apperror.CodeInternalServerError),
		)
	}
}

func WriteError(w http.ResponseWriter, err error) {
	var appErr *apperror.AppError

	if errors.As(err, &appErr) {
		writeAppError(w, appErr)
		return
	}

	writeUnknownError(w)
}

func writeAppError(w http.ResponseWriter, err *apperror.AppError) {
	status := statusFromCode(err.HTTPCode)

	WriteJSONResponse(w, status, APIErrorResponse{
		HTTPCode:      err.HTTPCode,
		Message:   err.Message,
		Details:   err.Details,
		RequestID: err.RequestID,
	})
}

func writeUnknownError(w http.ResponseWriter) {
	ctx := context.TODO()

	WriteJSONResponse(w, http.StatusInternalServerError, APIErrorResponse{
		HTTPCode:    apperror.CodeInternalServerError,
		Message: "internal server error",
		RequestID: requestcontext.GetRequestID(ctx),
	})
}