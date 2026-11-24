package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Safe development default: loopback only. Do not change this to 0.0.0.0
	// until authentication, CSRF protection, and hotspot-interface binding
	// have been implemented and reviewed.
	addr := os.Getenv("ENETS_LISTEN_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8080"
	}

	webDir := os.Getenv("ENETS_WEB_DIR")
	if webDir == "" {
		webDir = "web"
	}
	absWebDir, err := filepath.Abs(webDir)
	if err != nil {
		log.Fatalf("resolve web directory: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(absWebDir)))

	log.Printf("ENets development server listening on %s", addr)
	log.Printf("Serving dashboard assets from %s", absWebDir)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
