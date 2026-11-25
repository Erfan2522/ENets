package main

import ("testing"; "strings")

func TestDiagnosticLogSecretRedaction(t *testing.T) {
	input := "hostapd PASSWORD=supersecret\npassphrase: anothersecret\nnormal startup"
	got := secretLogPattern.ReplaceAllString(input, "$1$2[REDACTED]")
	if contains := got == input; contains { t.Fatal("expected diagnostic secrets to be redacted") }
	for _, secret := range []string{"supersecret", "anothersecret"} {
		if strings.Contains(got, secret) { t.Fatalf("secret %q remained in diagnostic output: %q", secret, got) }
	}
}
