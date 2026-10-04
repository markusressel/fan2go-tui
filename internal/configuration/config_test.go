package configuration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestConfiguration_DefaultsAndLoad(t *testing.T) {
	viper.Reset()
	InitConfig("")
	LoadConfig()

	if CurrentConfig.Api.Host != "127.0.0.1" {
		t.Errorf("expected default API host 127.0.0.1, got %s", CurrentConfig.Api.Host)
	}
	if CurrentConfig.Api.Port != 9001 {
		t.Errorf("expected default API port 9001, got %d", CurrentConfig.Api.Port)
	}
	if CurrentConfig.Ui.UpdateInterval != 500*time.Millisecond {
		t.Errorf("expected default UpdateInterval 500ms, got %v", CurrentConfig.Ui.UpdateInterval)
	}
	if CurrentConfig.Profiling.Enabled {
		t.Errorf("expected profiling disabled by default")
	}
}

func TestConfiguration_ConfigFileLoading(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.yaml")
	cfgContent := `
api:
  host: "192.168.1.100"
  port: 8080
ui:
  updateInterval: "1s"
`
	if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
		t.Fatalf("failed to write tmp config: %v", err)
	}

	viper.Reset()
	InitConfig(cfgPath)
	detected := DetectAndReadConfigFile()
	if detected != cfgPath {
		t.Errorf("expected detected config path %s, got %s", cfgPath, detected)
	}

	LoadConfig()
	if CurrentConfig.Api.Host != "192.168.1.100" {
		t.Errorf("expected host 192.168.1.100, got %s", CurrentConfig.Api.Host)
	}
	if CurrentConfig.Api.Port != 8080 {
		t.Errorf("expected port 8080, got %d", CurrentConfig.Api.Port)
	}
	if CurrentConfig.Ui.UpdateInterval != 1*time.Second {
		t.Errorf("expected UpdateInterval 1s, got %v", CurrentConfig.Ui.UpdateInterval)
	}
}
