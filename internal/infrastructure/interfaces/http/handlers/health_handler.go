package handlers

import (
	"time"
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/internal/pkg/utils"
	"github.com/EnockYator/saas-photo-listing-platform/internal/domain/health"
)

type HealthHandler struct {}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	status := health.HealthStatus{
		Status: "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
		Version: "1.0.0",
	}
	utils.SuccessResponse(w, http.StatusOK, status)
}