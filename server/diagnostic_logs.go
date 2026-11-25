package main

import (
	"net/http"
	"os"
	"regexp"
	"strings"
)

var secretLogPattern = regexp.MustCompile(`(?i)(password|passphrase|psk)(\s*[:=]\s*)(?:'[^']*'|"[^"]*"|[^\s&]+)`)

func diagnosticLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, err := os.ReadFile("/data/adb/modules/enets/enets.log")
	if err != nil {
		if os.IsNotExist(err) { data = []byte("No ENets log is available yet.") } else {
			http.Error(w, "diagnostic log could not be read", http.StatusInternalServerError)
			return
		}
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) > 120 { lines = lines[len(lines)-120:] }
	for i, line := range lines { lines[i] = secretLogPattern.ReplaceAllString(line, "$1$2[REDACTED]") }
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(strings.Join(lines, "\n")))
}
