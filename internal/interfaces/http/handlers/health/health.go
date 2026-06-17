package health

import (
	"net/http"

	healthdto "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/dto/health"
	"github.com/EnockYator/saas-photo-listing-platform/pkg/response"
)

// HealthCheck godoc
//
// @Summary Health check
// @Description Returns service health status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} healthdto.HealthResponse
// @Failure 500 {object} map[string]any
// @Router /health [get]
func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed", nil)
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
