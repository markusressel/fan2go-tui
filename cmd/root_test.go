package cmd

import (
	"bytes"
	"fan2go-tui/cmd/global"
	"testing"
)

func TestRootCmd_FlagsAndArgs(t *testing.T) {
	// Flags
	flags := []string{"config", "no-color", "no-style", "verbose"}
	for _, f := range flags {
		if rootCmd.PersistentFlags().Lookup(f) == nil {
			t.Errorf("missing persistent flag %q on rootCmd", f)
		}
	}

	// Args validation
	if err := rootCmd.Args(rootCmd, []string{}); err != nil {
		t.Errorf("expected 0 args to be valid, got: %v", err)
	}
	if err := rootCmd.Args(rootCmd, []string{"one"}); err != nil {
		t.Errorf("expected 1 arg to be valid, got: %v", err)
	}
	if err := rootCmd.Args(rootCmd, []string{"one", "two"}); err == nil {
		t.Errorf("expected >1 args to fail MaximumNArgs(1)")
	}
}

func TestSetupUi(t *testing.T) {
	// Save state
	oldVerbose := global.Verbose
	oldNoColor := global.NoColor
	oldNoStyle := global.NoStyle
	defer func() {
		global.Verbose = oldVerbose
		global.NoColor = oldNoColor
		global.NoStyle = oldNoStyle
	}()

	global.Verbose = true
	global.NoColor = true
	global.NoStyle = true
	setupUi()

	global.Verbose = false
	global.NoColor = false
	global.NoStyle = false
	setupUi()
}

func TestVersionCmd_Output(t *testing.T) {
	oldVerbose := global.Verbose
	oldLong := long
	oldVersion := global.Version
	oldCommit := global.Commit
	oldDate := global.Date
	defer func() {
		global.Verbose = oldVerbose
		long = oldLong
		global.Version = oldVersion
		global.Commit = oldCommit
		global.Date = oldDate
	}()

	global.Version = "v1.2.3"
	global.Commit = "abcdef"
	global.Date = "2026-10-04"

	// Standard version
	global.Verbose = false
	long = false
	buf := new(bytes.Buffer)
	versionCmd.SetOut(buf)
	versionCmd.Run(versionCmd, nil)

	// Long version
	long = true
	versionCmd.Run(versionCmd, nil)

	// Verbose version
	global.Verbose = true
	versionCmd.Run(versionCmd, nil)
}
