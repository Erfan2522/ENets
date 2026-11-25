package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	virtualAPBinary = "/data/local/enets/bin/virtualap"
	virtualAPConfig = "/data/local/enets/ap.conf"
)

var hotspotOperationMu sync.Mutex
var interfaceNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]{1,15}$`)

type hotspotStartRequest struct {
	SSID     string `json:"ssid"`
	Password string `json:"password"`
	Security string `json:"security"`
	Band     string `json:"band"`
	Channel  string `json:"channel"`
	Width    string `json:"width"`
	Hidden   bool   `json:"hidden"`
}

type hotspotStatus struct {
	BackendAvailable bool              `json:"backend_available"`
	Running          bool              `json:"running"`
	Mode             string            `json:"mode,omitempty"`
	SSID             string            `json:"ssid,omitempty"`
	Band             string            `json:"band,omitempty"`
	Channel          string            `json:"channel,omitempty"`
	Security         string            `json:"security,omitempty"`
	Clients          int               `json:"clients"`
	Details          map[string]string `json:"details,omitempty"`
}

type hotspotClient struct {
	MAC  string `json:"mac"`
	IP   string `json:"ip"`
	Name string `json:"name"`
}

func virtualAPAvailable() bool {
	info, err := os.Stat(virtualAPBinary)
	return err == nil && !info.IsDir() && info.Mode()&0111 != 0
}

func runVirtualAP(timeout time.Duration, args ...string) (string, error) {
	if !virtualAPAvailable() {
		return "", errors.New("hotspot backend is not installed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, virtualAPBinary, args...)
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return "", ctx.Err()
	}
	if err != nil {
		return "", fmt.Errorf("virtualap command failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func parseKeyValues(output string) map[string]string {
	values := make(map[string]string)
	for _, line := range strings.Split(output, "\n") {
		key, value, ok := strings.Cut(strings.TrimSpace(line), "=")
		if ok && key != "" {
			values[key] = value
		}
	}
	return values
}

func readHotspotStatus() (hotspotStatus, error) {
	result := hotspotStatus{BackendAvailable: virtualAPAvailable()}
	if !result.BackendAvailable {
		return result, nil
	}
	output, err := runVirtualAP(5*time.Second, "status")
	if err != nil {
		return result, err
	}
	values := parseKeyValues(output)
	result.Running = values["running"] == "1"
	result.Mode = values["mode"]
	result.SSID = values["ssid"]
	result.Band = values["band"]
	result.Channel = values["channel"]
	result.Security = values["security"]
	result.Details = values
	if n, err := strconv.Atoi(values["clients"]); err == nil && n >= 0 {
		result.Clients = n
	}
	return result, nil
}

func validTextField(value string, maxBytes int) bool {
	if !utf8.ValidString(value) || len(value) > maxBytes {
		return false
	}
	for _, r := range value {
		if r == '\n' || r == '\r' || r == 0 || r < 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}

func validateHotspotRequest(input hotspotStartRequest) (map[string]string, error) {
	if !validTextField(input.SSID, 32) || len(input.SSID) == 0 {
		return nil, errors.New("SSID must contain 1 to 32 UTF-8 bytes and no control characters")
	}
	if !validTextField(input.Password, 63) {
		return nil, errors.New("password contains unsupported characters or is too long")
	}

	security := input.Security
	if security == "" {
		security = "wpa2"
	}
	switch security {
	case "open":
		if input.Password != "" {
			return nil, errors.New("open security does not use a password")
		}
	case "wpa2", "wpa3", "wpa2wpa3":
		if len(input.Password) < 8 || len(input.Password) > 63 {
			return nil, errors.New("WPA password must contain 8 to 63 bytes")
		}
	default:
		return nil, errors.New("unsupported security mode")
	}

	band := input.Band
	switch band {
	case "2.4", "2":
		band = "2"
	case "5":
		band = "5"
	default:
		return nil, errors.New("only 2.4 GHz and 5 GHz are available in this backend")
	}

	channel := strings.TrimSpace(input.Channel)
	if channel == "auto" {
		channel = ""
	}
	if channel != "" {
		n, err := strconv.Atoi(channel)
		if err != nil {
			return nil, errors.New("channel must be Auto or a valid channel number")
		}
		if band == "2" && (n < 1 || n > 13) {
			return nil, errors.New("2.4 GHz channel must be between 1 and 13")
		}
		if band == "5" && (n < 36 || n > 177) {
			return nil, errors.New("5 GHz channel must be between 36 and 177")
		}
	}

	width := input.Width
	if width == "" {
		width = "auto"
	}
	switch width {
	case "auto", "20":
	case "40":
		if band == "5" || band == "2" {
			// The engine and driver decide whether 40 MHz is actually usable.
		}
	case "80":
		if band != "5" {
			return nil, errors.New("80 MHz width is available only on 5 GHz")
		}
	default:
		return nil, errors.New("unsupported channel width")
	}

	hidden := "0"
	if input.Hidden {
		hidden = "1"
	}

	return map[string]string{
		"PHY_IFACE":   "wlan0",
		"UPSTREAM":    "auto",
		"SSID":        input.SSID,
		"PASSWORD":    input.Password,
		"BAND":        band,
		"CHANNEL":     channel,
		"WIDTH":       width,
		"IP_GW":       "192.168.42.1",
		"DNS_SERVERS": "",
		"HIDDEN":      hidden,
		"SECURITY":    security,
		"PMF":         "0",
		"TTL_FIX":     "0",
		"IPV6":        "1",
	}, nil
}

// VirtualAP's config parser accepts shell-style single-quoted values. This
// encoder preserves embedded apostrophes without invoking a shell.
func quoteVirtualAPConfig(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func writeVirtualAPConfig(values map[string]string) error {
	if os.Geteuid() != 0 {
		return errors.New("ENets service is not running as root")
	}
	if err := os.MkdirAll(filepath.Dir(virtualAPConfig), 0700); err != nil {
		return err
	}
	keys := []string{"PHY_IFACE", "UPSTREAM", "SSID", "PASSWORD", "BAND", "CHANNEL", "WIDTH", "IP_GW", "DNS_SERVERS", "HIDDEN", "SECURITY", "PMF", "TTL_FIX", "IPV6"}
	var builder strings.Builder
	for _, key := range keys {
		value, ok := values[key]
		if !ok {
			continue
		}
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(quoteVirtualAPConfig(value))
		builder.WriteByte('\n')
	}
	tmp, err := os.CreateTemp(filepath.Dir(virtualAPConfig), ".ap.conf-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	if err := tmp.Chmod(0600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.WriteString(builder.String()); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, virtualAPConfig)
}

func sameOriginRequest(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return false
	}
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Host == "" || parsed.Host != r.Host || parsed.Scheme != "http" {
		return false
	}
	if site := r.Header.Get("Sec-Fetch-Site"); site != "" && site != "same-origin" {
		return false
	}
	return net.ParseIP(strings.Split(r.Host, ":")[0]) != nil || strings.HasPrefix(r.Host, "localhost")
}

func requireControlRequest(w http.ResponseWriter, r *http.Request) bool {
	if !sameOriginRequest(r) {
		http.Error(w, "same-origin request required", http.StatusForbidden)
		return false
	}
	if os.Geteuid() != 0 {
		http.Error(w, "ENets service must run as root", http.StatusForbidden)
		return false
	}
	if !virtualAPAvailable() {
		http.Error(w, "hotspot backend is not installed", http.StatusServiceUnavailable)
		return false
	}
	return true
}

func hotspotStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	hotspotOperationMu.Lock()
	defer hotspotOperationMu.Unlock()
	status, err := readHotspotStatus()
	if err != nil {
		http.Error(w, "could not read hotspot status", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(status)
}

func hotspotStartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !requireControlRequest(w, r) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8192)
	var input hotspotStartRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		http.Error(w, "invalid hotspot settings", http.StatusBadRequest)
		return
	}
	values, err := validateHotspotRequest(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hotspotOperationMu.Lock()
	defer hotspotOperationMu.Unlock()
	if err := writeVirtualAPConfig(values); err != nil {
		http.Error(w, "could not save hotspot settings", http.StatusInternalServerError)
		return
	}
	if _, err := runVirtualAP(60*time.Second, "start"); err != nil {
		http.Error(w, "hotspot could not be started; check the ENets backend log", http.StatusBadGateway)
		return
	}
	status, err := readHotspotStatus()
	if err != nil {
		http.Error(w, "hotspot started but status verification failed", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(status)
}

func hotspotStopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !requireControlRequest(w, r) {
		return
	}
	hotspotOperationMu.Lock()
	defer hotspotOperationMu.Unlock()
	if _, err := runVirtualAP(30*time.Second, "stop"); err != nil {
		http.Error(w, "hotspot could not be stopped; check the ENets backend log", http.StatusBadGateway)
		return
	}
	status, err := readHotspotStatus()
	if err != nil {
		http.Error(w, "hotspot stopped but status verification failed", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(status)
}

func hotspotClientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	hotspotOperationMu.Lock()
	defer hotspotOperationMu.Unlock()
	status, err := readHotspotStatus()
	if err != nil {
		http.Error(w, "could not read hotspot status", http.StatusBadGateway)
		return
	}
	clients := make([]hotspotClient, 0)
	if status.Running {
		output, err := runVirtualAP(5*time.Second, "clients")
		if err != nil {
			http.Error(w, "could not read connected clients", http.StatusBadGateway)
			return
		}
		for _, line := range strings.Split(output, "\n") {
			fields := strings.Fields(line)
			if len(fields) == 3 {
				clients = append(clients, hotspotClient{MAC: fields[0], IP: fields[1], Name: fields[2]})
			}
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(clients)
}
