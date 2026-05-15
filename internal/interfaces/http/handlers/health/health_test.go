package health_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http"
)

func TestHealthEndpoint(t *testing.T) {
	var db *sql.DB = nil
	router := httpserver.NewRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}
