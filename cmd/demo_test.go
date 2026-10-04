package cmd

import (
	"testing"
)

func TestDemoCommandRegistration(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"demo"})
	if err != nil {
		t.Fatalf("expected to find 'demo' command: %v", err)
	}

	if cmd == nil || cmd.Name() != "demo" {
		t.Fatalf("expected demo command, got %v", cmd)
	}

	hostFlag := cmd.Flag("host")
	if hostFlag == nil {
		t.Fatalf("missing --host flag")
	}
	if hostFlag.DefValue != "127.0.0.1" {
		t.Fatalf("expected default host 127.0.0.1, got %s", hostFlag.DefValue)
	}

	portFlag := cmd.Flag("port")
	if portFlag == nil {
		t.Fatalf("missing --port flag")
	}
	if portFlag.DefValue != "9001" {
		t.Fatalf("expected default port 9001, got %s", portFlag.DefValue)
	}
	if portFlag.Shorthand != "p" {
		t.Fatalf("expected shorthand -p for port, got %s", portFlag.Shorthand)
	}

	uiFlag := cmd.Flag("ui")
	if uiFlag == nil {
		t.Fatalf("missing --ui flag")
	}
	if uiFlag.DefValue != "false" {
		t.Fatalf("expected default ui false, got %s", uiFlag.DefValue)
	}
}
