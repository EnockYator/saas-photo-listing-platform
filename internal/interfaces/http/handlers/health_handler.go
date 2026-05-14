package handlers

import (
	"net/http"

	"github.com/EnockYator/saas-photo-listing-platform/pkg/response"
)

func Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "method not allowed", nil)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	}, map[string]any{
		"service": "saas-photo-listing-platform",
	})
}

func Live(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, map[string]string{
		"status": "alive",
	}, nil)
}

// func Ready(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		if err := db.Ping(); err != nil {
// 			response.WriteError(
// 				w,
// 				http.StatusServiceUnavailable,
// 				"DB_NOT_READY",
// 				"database connection failed",
// 				nil,
// 			)
// 			return
// 		}

// 		response.WriteJSON(w, http.StatusOK, map[string]string{
// 			"status": "ready",
// 		}, map[string]any{
// 			"db": "ok",
// 		})
// 	}
// }
