package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/state"
	"strings"
	"testing"

	"github.com/rivo/tview"
)

// A fan with warnings has a marked title and shows the messages in its info panel, until the warnings are gone.
func TestFanListItemComponentShowsWarnings(t *testing.T) {
	app := tview.NewApplication()
	warnings := []client.Warning{{Code: "flatRpmCurve", Message: "RPM stays at ~750 across the whole PWM range"}}
	c := NewFanListItemComponent(app, &state.FanState{
		Fan: &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Warnings: warnings},
	}, nil)

	if title := c.window.GetTitle(); title != " ⚠ fan-1 " {
		t.Errorf("expected the marked title, got %q", title)
	}
	if text := c.fanInfoComponent.warningsComponent.GetText(); !strings.Contains(text, "⚠ RPM stays at ~750") {
		t.Errorf("expected the warning in the info panel, got %q", text)
	}

	c.SetFan(&state.FanState{Fan: &client.Fan{Config: client.FanConfig{ID: "fan-1"}}})

	if title := c.window.GetTitle(); title != " fan-1 " {
		t.Errorf("expected the plain title, got %q", title)
	}
	if c.fanInfoComponent.warningsComponent.HasWarnings() {
		t.Errorf("expected no warnings in the info panel, got %q", c.fanInfoComponent.warningsComponent.GetText())
	}
}
