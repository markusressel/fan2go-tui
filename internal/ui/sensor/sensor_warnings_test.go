package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/state"
	"strings"
	"testing"

	"github.com/rivo/tview"
)

// A sensor with warnings has a marked title and shows the messages in its info panel, until the warnings are gone.
func TestSensorListItemComponentShowsWarnings(t *testing.T) {
	app := tview.NewApplication()
	warnings := []client.Warning{{Code: "test", Message: "sensor warning"}}
	c := NewSensorListItemComponent(app, &state.SensorState{
		Sensor: &client.Sensor{Config: client.SensorConfig{ID: "sensor-1"}, Warnings: warnings},
	})

	if title := c.window.GetTitle(); title != " ⚠ sensor-1 " {
		t.Errorf("expected the marked title, got %q", title)
	}
	if text := c.sensorInfoComponent.warningsComponent.GetText(); !strings.Contains(text, "⚠ sensor warning") {
		t.Errorf("expected the warning in the info panel, got %q", text)
	}

	c.SetSensor(&state.SensorState{Sensor: &client.Sensor{Config: client.SensorConfig{ID: "sensor-1"}}})

	if title := c.window.GetTitle(); title != " sensor-1 " {
		t.Errorf("expected the plain title, got %q", title)
	}
	if c.sensorInfoComponent.warningsComponent.HasWarnings() {
		t.Errorf("expected no warnings in the info panel, got %q", c.sensorInfoComponent.warningsComponent.GetText())
	}
}
