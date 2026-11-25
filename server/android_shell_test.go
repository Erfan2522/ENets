package main

import "testing"

func TestSoftAPHelpDetectionDoesNotAssumeInternetTethering(t *testing.T) {
	help := "Wi-Fi commands:\n start-softap <ssid> wpa2 <passphrase>\n stop-softap"
	if !containsHelpCommand(help, "start-softap") {
		t.Fatal("expected start-softap command to be detected")
	}
	if !containsHelpCommand(help, "stop-softap") {
		t.Fatal("expected stop-softap command to be detected")
	}
	if containsHelpCommand(help, "start-tethering") {
		t.Fatal("must not infer tethering support from Soft AP commands")
	}
}
