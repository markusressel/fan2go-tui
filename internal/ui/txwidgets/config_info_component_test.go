package txwidgets

import (
	"fan2go-tui/internal/client"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

func intPtr(i int) *int {
	return &i
}

func floatPtr(f float64) *float64 {
	return &f
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}

func strPtr(s string) *string {
	return &s
}

func TestConfigInfoComponent_RenderingAndClicks(t *testing.T) {
	comp := NewConfigInfoComponent()
	if comp.GetPrimitive() == nil {
		t.Fatalf("expected non-nil primitive")
	}

	var clickedLabel, clickedValue string
	comp.SetFieldClickablePredicate(func(sectionTitle, label, value string) bool {
		return label == "Curve"
	})
	comp.SetFieldClickHandler(func(sectionTitle, label, value string) {
		clickedLabel = label
		clickedValue = value
	})

	sections := []ConfigInfoSection{
		{
			Title:  "General",
			Accent: ConfigInfoAccentGeneral,
			Fields: []ConfigInfoField{
				{Label: "ID", Value: "fan-1"},
				{Label: "Curve", Value: "cpu-curve"},
			},
		},
		{
			Title:  "Source",
			Accent: ConfigInfoAccentSource,
			Fields: []ConfigInfoField{
				{Label: "Sensor", Value: "cpu-temp"},
			},
		},
	}

	comp.SetSections(sections)
	text := comp.layout.GetText(true)
	if !strings.Contains(text, "fan-1") || !strings.Contains(text, "cpu-curve") {
		t.Fatalf("expected rendered text to contain fields, got: %s", text)
	}

	// Trigger click on region
	comp.layout.Highlight("field-0")
	if clickedLabel != "Curve" || clickedValue != "cpu-curve" {
		t.Fatalf("expected click handler invoked with Curve/cpu-curve, got %s/%s", clickedLabel, clickedValue)
	}

	// Click unknown region ID
	comp.layout.Highlight("non-existent-region")

	// Test key input capture
	handler := comp.layout.InputHandler()
	if handler != nil {
		handler(tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone), nil)
		handler(tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModNone), nil)
		handler(tcell.NewEventKey(tcell.KeyRune, 'l', tcell.ModNone), nil)
		handler(tcell.NewEventKey(tcell.KeyRune, 'h', tcell.ModNone), nil)
	}

	// Scroll horizontal
	comp.ScrollHorizontal(4)
	comp.ScrollHorizontal(-4)

	// Set nil sections clears text
	comp.SetSections(nil)
	if comp.layout.GetText(true) != "" {
		t.Fatalf("expected empty text after setting nil sections")
	}
}

func TestFanConfigSections(t *testing.T) {
	activeMode := client.ControlModeValue("curve")
	cfg := client.FanConfig{
		ID:        "case-fan",
		Curve:     "my-curve",
		NeverStop: true,
		MinPwm:    intPtr(40),
		StartPwm:  intPtr(60),
		MaxPwm:    intPtr(255),
		HwMon: &client.HwMonFanConfig{
			Platform: "nct6798",
			Index:    1,
		},
		ControlMode: &client.ControlModeConfig{
			Active: &activeMode,
		},
		ControlAlgorithm: &client.ControlAlgorithmConfig{
			Direct: &client.DirectControlAlgorithmConfig{
				MaxPwmChangePerCycle: intPtr(10),
			},
		},
		SanityCheck: &client.SanityCheckConfig{
			PwmValueChangedByThirdParty: client.PwmValueChangedByThirdPartyConfig{
				Enabled: true,
			},
		},
	}

	sections := FanConfigSections(cfg)
	if len(sections) < 3 {
		t.Fatalf("expected at least 3 sections for fan config, got %d", len(sections))
	}
}

func TestCurveConfigSections(t *testing.T) {
	// Linear curve
	linearCfg := client.CurveConfig{
		ID: "linear-curve",
		Linear: &client.LinearCurveConfig{
			Sensor: "cpu-temp",
			Min:    30,
			Max:    80,
			Steps: map[int]float64{
				40: 100,
				60: 200,
			},
		},
	}
	s1 := CurveConfigSections(linearCfg)
	if len(s1) == 0 {
		t.Fatalf("expected sections for linear curve")
	}

	// Function curve
	funcCfg := client.CurveConfig{
		ID: "func-curve",
		Function: &client.FunctionCurveConfig{
			Type:   client.FunctionMaximum,
			Curves: []string{"curve-1", "curve-2"},
		},
	}
	s2 := CurveConfigSections(funcCfg)
	if len(s2) == 0 {
		t.Fatalf("expected sections for function curve")
	}

	// PID curve
	pidCfg := client.CurveConfig{
		ID: "pid-curve",
		PID: &client.PidCurveConfig{
			Sensor:   "cpu-temp",
			SetPoint: 60.0,
			P:        1.0,
			I:        0.1,
			D:        0.5,
		},
	}
	s3 := CurveConfigSections(pidCfg)
	if len(s3) == 0 {
		t.Fatalf("expected sections for pid curve")
	}
}

func TestSensorConfigSections(t *testing.T) {
	// HwMon sensor
	hwCfg := client.SensorConfig{
		ID: "hw-sensor",
		HwMon: &client.HwMonSensorConfig{
			Platform: "k10temp",
			Index:    0,
			Channel:  1,
		},
	}
	s1 := SensorConfigSections(hwCfg)
	if len(s1) == 0 {
		t.Fatalf("expected sections for hwmon sensor")
	}

	// Nvidia sensor
	nvCfg := client.SensorConfig{
		ID: "nv-sensor",
		Nvidia: &client.NvidiaSensorConfig{
			Device: "RTX 4090",
			Index:  0,
		},
	}
	s2 := SensorConfigSections(nvCfg)
	if len(s2) == 0 {
		t.Fatalf("expected sections for nvidia sensor")
	}

	// Disk sensor
	diskCfg := client.SensorConfig{
		ID: "disk-sensor",
		Disk: &client.DiskSensorConfig{
			Device: "/dev/nvme0n1",
		},
	}
	s3 := SensorConfigSections(diskCfg)
	if len(s3) == 0 {
		t.Fatalf("expected sections for disk sensor")
	}
}
