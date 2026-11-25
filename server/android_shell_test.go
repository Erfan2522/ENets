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

func TestCapabilityStatusesUseStableVocabulary(t *testing.T) {
	report := buildCapabilityReport()
	allowed := map[string]bool{
		"supported": true,
		"unsupported": true,
		"unknown": true,
		"error": true,
	}
	for name, item := range report.Capabilities {
		if !allowed[item.Status] {
			t.Errorf("capability %q has unsupported status %q", name, item.Status)
		}
	}
}

func TestSoftAPHelpDetectionDoesNotAssumeInternetTethering(t *testing.T) {
	help := "Wi-Fi commands:\n start-softap <ssid> wpa2 <passphrase>\n stop-softap"
	if !containsHelpCommand(help, "start-softap") {
		t.Fatal("expected start-softap command to be detected")
	}
	if !containsHelpCommand(help, "stop-softap") {
		t.Fatal("expected stop-softap command to be detected")
	}
	if containsHelpCommand(help, "start-tethering") {
		t.Fatal("must not infer tethering support from Soft AP commands")
	}
}
