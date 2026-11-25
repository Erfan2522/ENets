package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	settingsPath = "/data/local/enets/settings.json"
	autoStartPath = "/data/local/enets/autostart"
	maxProfiles = 12
)

type savedSettings struct {
	Hotspot hotspotStartRequest `json:"hotspot"`
	AutoStart bool `json:"auto_start"`
	UpdatedAt string `json:"updated_at"`
}

type hotspotProfile struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Hotspot hotspotStartRequest `json:"hotspot"`
}

type profilesDocument struct {
	Profiles []hotspotProfile `json:"profiles"`
}

func atomicPrivateWrite(path string, data []byte) error {
	if os.Geteuid() != 0 { return errors.New("ENets service must run as root") }
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil { return err }
	tmp, err := os.CreateTemp(filepath.Dir(path), ".enets-*")
	if err != nil { return err }
	name := tmp.Name()
	defer os.Remove(name)
	if err := tmp.Chmod(0600); err != nil { tmp.Close(); return err }
	if _, err := tmp.Write(data); err != nil { tmp.Close(); return err }
	if err := tmp.Sync(); err != nil { tmp.Close(); return err }
	if err := tmp.Close(); err != nil { return err }
	return os.Rename(name, path)
}

func loadSavedSettings() (savedSettings, error) {
	var saved savedSettings
	data, err := os.ReadFile(settingsPath)
	if errors.Is(err, os.ErrNotExist) { return saved, nil }
	if err != nil { return saved, err }
	if err := json.Unmarshal(data, &saved); err != nil { return saved, err }
	return saved, nil
}

func saveAutoStart(enabled bool) error {
	value := "0\n"
	if enabled { value = "1\n" }
	return atomicPrivateWrite(autoStartPath, []byte(value))
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		saved, err := loadSavedSettings()
		if err != nil { http.Error(w, "saved settings could not be read", http.StatusInternalServerError); return }
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(saved)
	case http.MethodPut:
		if !requireControlRequest(w, r) { return }
		r.Body = http.MaxBytesReader(w, r.Body, 8192)
		var saved savedSettings
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&saved); err != nil { http.Error(w, "invalid settings payload", http.StatusBadRequest); return }
		values, err := validateHotspotRequest(saved.Hotspot)
		if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		phy, err := defaultPHYInterface()
		if err != nil { http.Error(w, err.Error(), http.StatusServiceUnavailable); return }
		values["PHY_IFACE"] = phy
		hotspotOperationMu.Lock()
		defer hotspotOperationMu.Unlock()
		if err := writeVirtualAPConfig(values); err != nil { http.Error(w, "hotspot settings could not be saved", http.StatusInternalServerError); return }
		saved.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		data, err := json.Marshal(saved)
		if err != nil { http.Error(w, "settings could not be encoded", http.StatusInternalServerError); return }
		if err := atomicPrivateWrite(settingsPath, data); err != nil { http.Error(w, "settings could not be saved", http.StatusInternalServerError); return }
		if err := saveAutoStart(saved.AutoStart); err != nil { http.Error(w, "auto-start preference could not be saved", http.StatusInternalServerError); return }
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(saved)
	default:
		w.Header().Set("Allow", "GET, PUT")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func loadProfiles() (profilesDocument, error) {
	var doc profilesDocument
	path := "/data/local/enets/profiles.json"
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) { return doc, nil }
	if err != nil { return doc, err }
	err = json.Unmarshal(data, &doc)
	return doc, err
}

func profilesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		doc, err := loadProfiles()
		if err != nil { http.Error(w, "profiles could not be read", http.StatusInternalServerError); return }
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		_ = json.NewEncoder(w).Encode(doc.Profiles)
	case http.MethodPost:
		if !requireControlRequest(w, r) { return }
		r.Body = http.MaxBytesReader(w, r.Body, 8192)
		var profile hotspotProfile
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&profile); err != nil { http.Error(w, "invalid profile", http.StatusBadRequest); return }
		profile.Name = strings.TrimSpace(profile.Name)
		if profile.Name == "" || len(profile.Name) > 32 { http.Error(w, "profile name must contain 1 to 32 characters", http.StatusBadRequest); return }
		if _, err := validateHotspotRequest(profile.Hotspot); err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
		doc, err := loadProfiles()
		if err != nil { http.Error(w, "profiles could not be read", http.StatusInternalServerError); return }
		if len(doc.Profiles) >= maxProfiles { http.Error(w, "profile limit reached", http.StatusConflict); return }
		id := make([]byte, 8)
		if _, err := rand.Read(id); err != nil { http.Error(w, "profile ID could not be generated", http.StatusInternalServerError); return }
		profile.ID = hex.EncodeToString(id)
		doc.Profiles = append(doc.Profiles, profile)
		data, err := json.Marshal(doc)
		if err == nil { err = atomicPrivateWrite("/data/local/enets/profiles.json", data) }
		if err != nil { http.Error(w, "profile could not be saved", http.StatusInternalServerError); return }
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(profile)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func diagnosticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	report := map[string]any{
		"root": os.Geteuid() == 0,
		"backend_available": virtualAPAvailable(),
		"settings_file": fileExists(settingsPath),
		"profiles_file": fileExists("/data/local/enets/profiles.json"),
		"auto_start": readAutoStart(),
		"generated_at": time.Now().UTC().Format(time.RFC3339),
		"log_path": "/data/adb/modules/enets/enets.log",
		"warnings": []string{"Diagnostic output excludes hotspot passwords and saved profile secrets.", "WPA3 support is not declared until a device-specific SoftAP capability check exists."},
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(report)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func readAutoStart() bool {
	data, err := os.ReadFile(autoStartPath)
	return err == nil && strings.TrimSpace(string(data)) == "1"
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(map[string]string{
		"version": "v0.3.0-alpha",
		"release_url": "https://github.com/Erfan2522/ENets/releases",
		"update_mode": "manual_confirmation_required",
	})
}
