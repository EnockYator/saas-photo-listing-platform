package health

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/pkg/response"
)

func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed", nil)
		return
	}

	response.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{"status": "ok"},
		map[string]any{"service": "saas-photo-listing-platform"},
	)
}
