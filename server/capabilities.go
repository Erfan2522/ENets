package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type capability struct {
	Status  string   `json:"status"`
	Details string   `json:"details"`
	Evidence []string `json:"evidence,omitempty"`
}

type capabilityReport struct {
	GeneratedAt string                `json:"generated_at"`
	Platform    map[string]string     `json:"platform"`
	Interfaces  []string              `json:"wifi_interfaces"`
	Capabilities map[string]capability `json:"capabilities"`
	Warnings    []string               `json:"warnings"`
}

func commandOutput(timeout time.Duration, name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func property(name string) string {
	return commandOutput(2*time.Second, "getprop", name)
}

func detectWiFiInterfaces() []string {
	entries, err := filepath.Glob("/sys/class/net/*")
	if err != nil {
		return nil
	}
	var found []string
	for _, entry := range entries {
		if _, err := os.Stat(filepath.Join(entry, "wireless")); err == nil {
			found = append(found, filepath.Base(entry))
			continue
		}
		if _, err := os.Stat(filepath.Join(entry, "phy80211")); err == nil {
			found = append(found, filepath.Base(entry))
		}
	}
	return found
}

var frequencyPattern = regexp.MustCompile(`(?m)\*\s*([0-9]{4})\s*MHz`)

func detectHardwareBands() (map[string]capability, []string) {
	result := map[string]capability{
		"2.4_ghz": {Status: "unknown", Details: "No reliable radio capability data was found."},
		"5_ghz":   {Status: "unknown", Details: "No reliable radio capability data was found."},
		"6_ghz":   {Status: "unknown", Details: "No reliable radio capability data was found."},
	}
	iw := commandOutput(3*time.Second, "iw", "list")
	if iw == "" {
		return result, nil
	}
	var evidence []string
	seen := map[string]bool{}
	for _, match := range frequencyPattern.FindAllStringSubmatch(iw, -1) {
		if len(match) < 2 {
			continue
		}
		var mhz int
		for _, r := range match[1] {
			mhz = mhz*10 + int(r-'0')
		}
		band := ""
		switch {
		case mhz >= 2400 && mhz < 2500:
			band = "2.4_ghz"
		case mhz >= 4900 && mhz < 5925:
			band = "5_ghz"
		case mhz >= 5925 && mhz <= 7125:
			band = "6_ghz"
		}
		if band != "" {
			seen[band] = true
		}
	}
	for band := range seen {
		result[band] = capability{
			Status: "unknown",
			Details: "The Wi-Fi driver reports frequencies in this band, but hotspot AP support has not been validated.",
			Evidence: []string{"iw list frequency data"},
		}
		evidence = append(evidence, band)
	}
	return result, evidence
}

func buildCapabilityReport() capabilityReport {
	interfaces := detectWiFiInterfaces()
	bands, _ := detectHardwareBands()
	api := property("ro.build.version.sdk")
	android := property("ro.build.version.release")
	abi := property("ro.product.cpu.abi")
	device := property("ro.product.device")
	if api == "" { api = "unknown" }
	if android == "" { android = "unknown" }
	if abi == "" { abi = "unknown" }
	if device == "" { device = "unknown" }

	report := capabilityReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Platform: map[string]string{
			"android_release": android,
			"android_api": api,
			"primary_abi": abi,
			"device": device,
			"manufacturer": property("ro.product.manufacturer"),
			"model": property("ro.product.model"),
		},
		Interfaces: interfaces,
		Capabilities: map[string]capability{
			"ssid_change": {Status: "unknown", Details: "No tested Android hotspot adapter is integrated. This feature is not implemented."},
			"password_change": {Status: "unknown", Details: "No tested Android hotspot adapter is integrated. This feature is not implemented."},
			"security_mode": {Status: "unknown", Details: "Hotspot security mode cannot be safely inferred from generic root shell access."},
			"hotspot_enable_disable": {Status: "unknown", Details: "Hotspot control is not implemented."},
			"connected_clients": {Status: "unknown", Details: "Read-only client discovery is not implemented yet."},
			"band_2_4_ghz": bands["2.4_ghz"],
			"band_5_ghz": bands["5_ghz"],
			"band_6_ghz": bands["6_ghz"],
		},
		Warnings: []string{
			"This endpoint performs read-only detection. It does not change hotspot settings.",
			"Reported radio frequencies are hardware/driver clues, not proof of hotspot AP support.",
			"Root access alone does not guarantee that Android or vendor hotspot settings can be changed.",
		},
	}
	if len(interfaces) == 0 {
		report.Warnings = append(report.Warnings, "No Wi-Fi interface was identified through standard sysfs markers.")
	}
	return report
}

func capabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(buildCapabilityReport())
}
