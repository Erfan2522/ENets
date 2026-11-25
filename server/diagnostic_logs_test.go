package main

import (
	"strings"
	"testing"
)

func TestDiagnosticLogSecretRedaction(t *testing.T) {
	input := "PASSWORD=supersecret\npassphrase: anothersecret\nPSK='quoted secret'\nnormal startup"
	got := secretLogPattern.ReplaceAllString(input, "$1$2[REDACTED]")
	for _, secret := range []string{"supersecret", "anothersecret", "quoted secret"} {
		if strings.Contains(got, secret) { t.Fatalf("secret %q remained in diagnostic output: %q", secret, got) }
	}
	if !strings.Contains(got, "[REDACTED]") { t.Fatalf("expected redaction marker: %q", got) }
}
