package response

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/apperror"
)

func statusFromCode(code apperror.ErrorCode) int {
	status, ok := statusByCode[code]
	if !ok {
		return http.StatusInternalServerError
	}

	return status
}
