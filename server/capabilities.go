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
	Status   string   `json:"status"`
	Details  string   `json:"details"`
	Evidence []string `json:"evidence,omitempty"`
}

type capabilityReport struct {
	GeneratedAt  string                `json:"generated_at"`
	Platform     map[string]string     `json:"platform"`
	Interfaces   []string              `json:"wifi_interfaces"`
	Capabilities map[string]capability `json:"capabilities"`
	Warnings     []string              `json:"warnings"`
}

type androidShellEvidence struct {
	WiFiHelpAvailable       bool
	SoftAPStartListed       bool
	SoftAPStopListed        bool
	TetheringHelpAvailable  bool
	ConnectivityHelpPresent bool
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

func commandOutputCombined(timeout time.Duration, name string, args ...string) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil && len(out) == 0 {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func property(name string) string {
	return commandOutput(2*time.Second, "getprop", name)
}

func containsHelpCommand(help, command string) bool {
	pattern := regexp.MustCompile(`(?m)(^|\s)` + regexp.QuoteMeta(command) + `(\s|$)`)
	return pattern.MatchString(strings.ToLower(help))
}

func detectAndroidShellEvidence() androidShellEvidence {
	var result androidShellEvidence
	wifiHelp := commandOutputCombined(3*time.Second, "cmd", "wifi", "help")
	if wifiHelp != "" {
		result.WiFiHelpAvailable = true
		result.SoftAPStartListed = containsHelpCommand(wifiHelp, "start-softap")
		result.SoftAPStopListed = containsHelpCommand(wifiHelp, "stop-softap")
	}
	result.TetheringHelpAvailable = commandOutputCombined(3*time.Second, "cmd", "tethering", "help") != ""
	result.ConnectivityHelpPresent = commandOutputCombined(3*time.Second, "cmd", "connectivity", "help") != ""
	return result
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
	shell := detectAndroidShellEvidence()
	api := property("ro.build.version.sdk")
	android := property("ro.build.version.release")
	abi := property("ro.product.cpu.abi")
	device := property("ro.product.device")
	if api == "" { api = "unknown" }
	if android == "" { android = "unknown" }
	if abi == "" { abi = "unknown" }
	if device == "" { device = "unknown" }

	wifiShell := "not_detected"
	if shell.WiFiHelpAvailable { wifiShell = "available" }
	tetherShell := "not_detected"
	if shell.TetheringHelpAvailable { tetherShell = "available" }
	connectivityShell := "not_detected"
	if shell.ConnectivityHelpPresent { connectivityShell = "available" }

	softAPDetails := "The Android Wi-Fi shell interface was not detected."
	softAPEvidence := []string{}
	if shell.SoftAPStartListed {
		softAPDetails = "The device lists cmd wifi start-softap. This proves command discovery only, not permission, successful AP startup, or internet tethering."
		softAPEvidence = append(softAPEvidence, "cmd wifi help: start-softap")
	}
	if shell.SoftAPStopListed {
		softAPEvidence = append(softAPEvidence, "cmd wifi help: stop-softap")
	}

	backendCapability := capability{Status: "unsupported", Details: "The VirtualAP engine payload is not installed."}
	if virtualAPAvailable() {
		backendCapability = capability{Status: "unknown", Details: "The pinned VirtualAP engine is installed. Device-specific hotspot startup and internet sharing have not yet been verified."}
	}

	report := capabilityReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Platform: map[string]string{
			"android_release": android,
			"android_api": api,
			"primary_abi": abi,
			"device": device,
			"manufacturer": property("ro.product.manufacturer"),
			"model": property("ro.product.model"),
			"wifi_shell": wifiShell,
			"tethering_shell": tetherShell,
			"connectivity_shell": connectivityShell,
		},
		Interfaces: interfaces,
		Capabilities: map[string]capability{
			"virtualap_backend": backendCapability,
			"ssid_change": {Status: "unknown", Details: "No device-specific adapter has safely read and applied an SSID yet."},
			"password_change": {Status: "unknown", Details: "No device-specific adapter has safely read and applied a hotspot passphrase yet."},
			"security_mode": {Status: "unknown", Details: "Hotspot security mode cannot be safely inferred from generic root shell access."},
			"hotspot_enable_disable": {Status: "unknown", Details: "The backend can attempt start/stop, but a successful AP and client connectivity have not yet been validated on this device."},
			"softap_shell_command": {Status: "unknown", Details: softAPDetails, Evidence: softAPEvidence},
			"internet_tethering": {Status: "unknown", Details: "No safe, validated Android tethering adapter is integrated. Soft AP startup is not equivalent to internet sharing."},
			"connected_clients": {Status: "unknown", Details: "Read-only client discovery is not implemented yet."},
			"band_2_4_ghz": bands["2.4_ghz"],
			"band_5_ghz": bands["5_ghz"],
			"band_6_ghz": bands["6_ghz"],
		},
		Warnings: []string{
			"This endpoint performs read-only detection. It does not change hotspot settings.",
			"The bundled VirtualAP engine is an experimental backend; its presence is not proof of device compatibility.",
			"Shell help output confirms only that a command is listed; it does not prove permissions or successful operation.",
			"Reported radio frequencies are hardware/driver clues, not proof of hotspot AP support.",
			"Soft AP startup and internet tethering are separate capabilities and must be validated independently.",
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
