package health

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/pkg/response"
)

func Live(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(
		w,
		http.StatusOK,
		map[string]string{"status": "alive"},
		nil)
}
