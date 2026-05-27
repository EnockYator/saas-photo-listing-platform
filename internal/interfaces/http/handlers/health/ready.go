package health

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/EnockYator/saas-photo-listing-platform/pkg/response"
)

func Ready(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			response.WriteError(
				w,
				http.StatusServiceUnavailable,
				"DB_NOT_READY",
				"database not reachable",
				nil,
			)
			return
		}

		response.WriteJSON(
			w,
			http.StatusOK,
			map[string]string{"status": "ready"},
			map[string]any{"db": "ok"},
		)
	}
}
