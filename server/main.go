package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	// ENets deliberately binds to loopback only. Do not make this configurable
	// until authentication, CSRF protection, and interface restrictions exist.
	const addr = "127.0.0.1:8080"

	webDir := os.Getenv("ENETS_WEB_DIR")
	if webDir == "" {
		webDir = "web"
	}
	absWebDir, err := filepath.Abs(webDir)
	if err != nil {
		log.Fatalf("resolve web directory: %v", err)
	}
	if info, err := os.Stat(absWebDir); err != nil || !info.IsDir() {
		log.Fatalf("web directory is unavailable: %s", absWebDir)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/health", healthHandler)
	mux.HandleFunc("/api/v1/capabilities", capabilitiesHandler)
	mux.HandleFunc("/api/v1/hotspot/status", hotspotStatusHandler)
	mux.HandleFunc("/api/v1/hotspot/start", hotspotStartHandler)
	mux.HandleFunc("/api/v1/hotspot/stop", hotspotStopHandler)
	mux.HandleFunc("/api/v1/hotspot/clients", hotspotClientsHandler)
	mux.Handle("/", http.FileServer(http.Dir(absWebDir)))

	server := &http.Server{
		Addr:              addr,
		Handler:           securityHeaders(localHostOnly(mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	log.Printf("ENets local server listening on %s", addr)
	log.Printf("Serving dashboard assets from %s", absWebDir)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server stopped: %v", err)
	}
}

func localHostOnly(next http.Handler) http.Handler {
	allowedHosts := map[string]bool{
		"127.0.0.1:8080": true,
		"localhost:8080": true,
		"[::1]:8080":     true,
		"127.0.0.1":      true,
		"localhost":      true,
		"[::1]":          true,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host := strings.ToLower(r.Host)
		if !allowedHosts[host] {
			http.Error(w, "host not allowed", http.StatusForbidden)
			return
		}
		remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil || net.ParseIP(remoteHost) == nil || !net.ParseIP(remoteHost).IsLoopback() {
			http.Error(w, "local access only", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; object-src 'none'; base-uri 'none'; frame-ancestors 'none'")
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok", "mode": "read_only_preview"})
}
