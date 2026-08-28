package health

import (
	"net/http"

	healthdto "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/dto/health"
	"github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http/response"
)

// LiveCheck godoc
//
// @Summary Health Live check
// @Description Returns service health live status
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} healthdto.HealthResponse
// @Failure 500 {object} map[string]any
// @Router /health/live [get]
func Live(w http.ResponseWriter, r *http.Request) {
	response.WriteResponse(
		w,
		http.StatusOK,
		healthdto.HealthResponse{
			Status:      "alive",
			Application: "saas-photo-listing-platform",
		},
	)
}
