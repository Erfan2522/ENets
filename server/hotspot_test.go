package main

import (
	"net/http/httptest"
	"testing"
)

func TestValidateHotspotRequest(t *testing.T) {
	input := hotspotStartRequest{SSID: "ENets Test", Password: "password123", Security: "wpa2", Band: "2.4", Channel: "auto", Width: "auto"}
	values, err := validateHotspotRequest(input)
	if err != nil { t.Fatalf("valid request rejected: %v", err) }
	if values["BAND"] != "2" || values["SSID"] != input.SSID || values["SECURITY"] != "wpa2" { t.Fatalf("unexpected normalized config: %#v", values) }
}

func TestValidateHotspotRequestRejectsUnsupportedSettings(t *testing.T) {
	cases := []struct { name string; edit func(*hotspotStartRequest) }{
		{"6 GHz", func(v *hotspotStartRequest) { v.Band = "6" }},
		{"short password", func(v *hotspotStartRequest) { v.Password = "short" }},
		{"newline SSID", func(v *hotspotStartRequest) { v.SSID = "bad\nssid" }},
		{"invalid 2.4 GHz channel", func(v *hotspotStartRequest) { v.Channel = "14" }},
		{"80 MHz on 2.4 GHz", func(v *hotspotStartRequest) { v.Width = "80" }},
	}
	for _, tc := range cases { t.Run(tc.name, func(t *testing.T) {
		input := hotspotStartRequest{SSID: "ENets", Password: "password123", Security: "wpa2", Band: "2.4", Channel: "auto", Width: "auto"}
		tc.edit(&input)
		if _, err := validateHotspotRequest(input); err == nil { t.Fatal("expected validation error") }
	}) }
}

func TestValidateOpenSecurity(t *testing.T) {
	input := hotspotStartRequest{SSID: "OpenNet", Security: "open", Band: "5", Width: "80"}
	values, err := validateHotspotRequest(input)
	if err != nil { t.Fatalf("valid open network rejected: %v", err) }
	if values["PASSWORD"] != "" || values["SECURITY"] != "open" { t.Fatalf("unexpected open config: %#v", values) }
}

func TestWPA3RejectedUntilDeviceCapabilityIsValidated(t *testing.T) {
	for _, mode := range []string{"wpa3", "wpa2wpa3"} {
		input := hotspotStartRequest{SSID: "ENets", Password: "password123", Security: mode, Band: "2.4", Width: "auto"}
		if _, err := validateHotspotRequest(input); err == nil { t.Fatalf("%s must remain unavailable without validated SoftAP SAE support", mode) }
	}
}

func TestQuoteVirtualAPConfig(t *testing.T) {
	got := quoteVirtualAPConfig("owner's network")
	want := "'owner'\\''s network'"
	if got != want { t.Fatalf("quoteVirtualAPConfig() = %q, want %q", got, want) }
}

func TestSameOriginRequest(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1:8080/api/v1/hotspot/start", nil)
	req.Host = "127.0.0.1:8080"
	req.Header.Set("Origin", "http://127.0.0.1:8080")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	if !sameOriginRequest(req) { t.Fatal("same-origin request was rejected") }
	req.Header.Set("Origin", "http://attacker.example")
	if sameOriginRequest(req) { t.Fatal("cross-origin request was accepted") }
}
