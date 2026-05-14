package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	httpserver "github.com/EnockYator/saas-photo-listing-platform/internal/interfaces/http"
)

func TestHealthEndpoint(t *testing.T) {
	router := httpserver.NewRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d", rec.Code)
	}
}