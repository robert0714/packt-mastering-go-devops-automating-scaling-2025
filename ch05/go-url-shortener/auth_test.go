package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthenticationMiddleware(t *testing.T) {
	handler := AuthenticationMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	)
	req, _ := http.NewRequest("GET", "/", nil)
	// Test without API key
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Errorf("expected status 403, got %d", rec.Code)
	}
	// Test with API key
	req.Header.Set("X-API-Key", "key-admin-123")
	rec = httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
}
