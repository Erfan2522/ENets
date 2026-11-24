package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCapabilitiesHandlerRejectsNonGet(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/capabilities", nil)
	rec := httptest.NewRecorder()
	capabilitiesHandler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected %d, got %d", http.StatusMethodNotAllowed, rec.Code)
	}
	if rec.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("expected Allow: GET, got %q", rec.Header().Get("Allow"))
	}
}

func TestCapabilitiesHandlerReturnsJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/capabilities", nil)
	rec := httptest.NewRecorder()
	capabilitiesHandler(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("unexpected content type: %q", got)
	}
	if rec.Body.Len() == 0 {
		t.Fatal("expected a JSON response body")
	}
}
