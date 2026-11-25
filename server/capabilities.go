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
	Status string `json:"status"`
	Details string `json:"details"`
	Evidence []string `json:"evidence,omitempty"`
}

type capabilityReport struct {
	GeneratedAt string `json:"generated_at"`
	Platform map[string]string `json:"platform"`
	Interfaces []string `json:"wifi_interfaces"`
	Capabilities map[string]capability `json:"capabilities"`
	Warnings []string `json:"warnings"`
}

type androidShellEvidence struct {
	WiFiHelpAvailable bool
	SoftAPStartListed bool
	SoftAPStopListed bool
	TetheringHelpAvailable bool
	ConnectivityHelpPresent bool
}

func commandOutput(timeout time.Duration, name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout); defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil { return "" }
	return strings.TrimSpace(string(out))
}

func commandOutputCombined(timeout time.Duration, name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout); defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil && len(out) == 0 { return "" }
	return strings.TrimSpace(string(out))
}

func property(name string) string { return commandOutput(2*time.Second, "getprop", name) }

func containsHelpCommand(help, command string) bool {
	pattern := regexp.MustCompile(`(?m)(^|\s)` + regexp.QuoteMeta(command) + `(\s|$)`)
	return pattern.MatchString(strings.ToLower(help))
}

func detectAndroidShellEvidence() androidShellEvidence {
	var result androidShellEvidence
	wifiHelp := commandOutputCombined(3*time.Second, "cmd", "wifi", "help")
	if wifiHelp != "" { result.WiFiHelpAvailable = true; result.SoftAPStartListed = containsHelpCommand(wifiHelp, "start-softap"); result.SoftAPStopListed = containsHelpCommand(wifiHelp, "stop-softap") }
	result.TetheringHelpAvailable = commandOutputCombined(3*time.Second, "cmd", "tethering", "help") != ""
	result.ConnectivityHelpPresent = commandOutputCombined(3*time.Second, "cmd", "connectivity", "help") != ""
	return result
}

func detectWiFiInterfaces() []string {
	entries, err := filepath.Glob("/sys/class/net/*"); if err != nil { return nil }
	var found []string
	for _, entry := range entries {
		if _, err := os.Stat(filepath.Join(entry, "wireless")); err == nil { found = append(found, filepath.Base(entry)); continue }
		if _, err := os.Stat(filepath.Join(entry, "phy80211")); err == nil { found = append(found, filepath.Base(entry)) }
	}
	return found
}

var frequencyPattern = regexp.MustCompile(`(?m)\*\s*([0-9]{4})\s*MHz`)

func detectHardwareBands() (map[string]capability, []string) {
	result := map[string]capability{
		"2.4_ghz": {Status: "unknown", Details: "No reliable radio capability data was found."},
		"5_ghz": {Status: "unknown", Details: "No reliable radio capability data was found."},
		"6_ghz": {Status: "unknown", Details: "No reliable radio capability data was found."},
	}
	iw := commandOutput(3*time.Second, "iw", "list"); if iw == "" { return result, nil }
	seen := map[string]bool{}
	for _, match := range frequencyPattern.FindAllStringSubmatch(iw, -1) {
		if len(match) < 2 { continue }
		var mhz int; for _, r := range match[1] { mhz = mhz*10 + int(r-'0') }
		band := ""
		switch { case mhz >= 2400 && mhz < 2500: band = "2.4_ghz"; case mhz >= 4900 && mhz < 5925: band = "5_ghz"; case mhz >= 5925 && mhz <= 7125: band = "6_ghz" }
		if band != "" { seen[band] = true }
	}
	var evidence []string
	for band := range seen {
		result[band] = capability{Status: "unknown", Details: "The Wi-Fi driver reports frequencies in this band, but hotspot AP support has not been validated.", Evidence: []string{"iw list frequency data"}}
		evidence = append(evidence, band)
	}
	return result, evidence
}

// Generic Android radio data cannot prove WPA3 SAE support in SoftAP mode.
func wpa3SoftAPValidated() bool { return false }

func buildCapabilityReport() capabilityReport {
	interfaces := detectWiFiInterfaces()
	bands, _ := detectHardwareBands()
	shell := detectAndroidShellEvidence()
	api, android, abi, device := property("ro.build.version.sdk"), property("ro.build.version.release"), property("ro.product.cpu.abi"), property("ro.product.device")
	if api == "" { api = "unknown" }; if android == "" { android = "unknown" }; if abi == "" { abi = "unknown" }; if device == "" { device = "unknown" }
	wifiShell, tetherShell, connectivityShell := "not_detected", "not_detected", "not_detected"
	if shell.WiFiHelpAvailable { wifiShell = "available" }; if shell.TetheringHelpAvailable { tetherShell = "available" }; if shell.ConnectivityHelpPresent { connectivityShell = "available" }
	softAPDetails := "The Android Wi-Fi shell interface was not detected."
	softAPEvidence := []string{}
	if shell.SoftAPStartListed { softAPDetails = "The device lists cmd wifi start-softap. This proves command discovery only, not permission or successful AP startup."; softAPEvidence = append(softAPEvidence, "cmd wifi help: start-softap") }
	if shell.SoftAPStopListed { softAPEvidence = append(softAPEvidence, "cmd wifi help: stop-softap") }
	backendCapability := capability{Status: "unsupported", Details: "The VirtualAP engine payload is not installed."}
	if virtualAPAvailable() { backendCapability = capability{Status: "unknown", Details: "The engine is installed; device-specific startup and internet sharing still require validation."} }
	report := capabilityReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Platform: map[string]string{"android_release": android, "android_api": api, "primary_abi": abi, "device": device, "manufacturer": property("ro.product.manufacturer"), "model": property("ro.product.model"), "wifi_shell": wifiShell, "tethering_shell": tetherShell, "connectivity_shell": connectivityShell},
		Interfaces: interfaces,
		Capabilities: map[string]capability{
			"virtualap_backend": backendCapability,
			"ssid_change": {Status: "unknown", Details: "Device-specific configuration validation is still required."},
			"password_change": {Status: "unknown", Details: "Device-specific configuration validation is still required."},
			"security_mode": {Status: "unknown", Details: "Generic root access does not prove hotspot security support."},
			"wpa3_softap": {Status: "unknown", Details: "WPA3 SoftAP SAE support has not been positively validated. WPA3 controls remain disabled for safety."},
			"hotspot_enable_disable": {Status: "unknown", Details: "A successful AP and client connectivity must be verified on this device."},
			"softap_shell_command": {Status: "unknown", Details: softAPDetails, Evidence: softAPEvidence},
			"internet_tethering": {Status: "unknown", Details: "No validated Android tethering adapter is integrated."},
			"connected_clients": {Status: "unknown", Details: "Client metadata is supplied by the bundled VirtualAP engine."},
			"band_2_4_ghz": bands["2.4_ghz"], "band_5_ghz": bands["5_ghz"], "band_6_ghz": bands["6_ghz"],
		},
		Warnings: []string{
			"This endpoint performs read-only detection.",
			"Backend presence does not prove device compatibility.",
			"Shell help output does not prove permission or successful operation.",
			"Reported radio frequencies are not proof of hotspot AP support.",
			"Soft AP startup and internet tethering must be validated independently.",
			"WPA3 is disabled until SoftAP SAE support is positively validated for this device and firmware.",
		},
	}
	if len(interfaces) == 0 { report.Warnings = append(report.Warnings, "No Wi-Fi interface was identified through standard sysfs markers.") }
	return report
}

func capabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { w.Header().Set("Allow", http.MethodGet); http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
	w.Header().Set("Content-Type", "application/json; charset=utf-8"); w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(buildCapabilityReport())
}
