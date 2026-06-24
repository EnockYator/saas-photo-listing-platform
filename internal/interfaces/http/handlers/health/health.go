package health

import (
	"net/http"

	healthdto "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/dto/health"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/utilities/apperror"
	"github.com/EnockYator/saas-photo-listing-platform/internal/shared/utilities/response"
)

// HealthCheck godoc
//
// @Summary Health check
// @Description Returns service health status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} healthdto.HealthResponse
// @Failure 405 {object} response.APIResponse
// @Router /health [get]
func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.HandleError(
			w,
			apperror.New(
				r.Context(),
				http.StatusMethodNotAllowed,
				apperror.CodeMethodNotAllowed,
				"method not allowed",
				nil,
			),
		)
		return
	}

	response.WriteJSON(
		w,
		http.StatusOK,
		healthdto.HealthResponse{
			Status:  "ok",
			Service: "saas-photo-listing-platform",
		},
	)
}