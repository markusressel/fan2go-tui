package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/state"
	"strings"
	"testing"

	"github.com/rivo/tview"
)

// A curve with warnings has a marked title and shows the messages in its info panel, until the warnings are gone.
func TestCurveListItemComponentShowsWarnings(t *testing.T) {
	app := tview.NewApplication()
	warnings := []client.Warning{{Code: "test", Message: "curve warning"}}
	c := NewCurveListItemComponent(app, &state.CurveState{
		Curve: &client.Curve{Config: client.CurveConfig{ID: "curve-1"}, Warnings: warnings},
	}, nil, nil)

	if title := c.window.GetTitle(); title != " ⚠ curve-1 " {
		t.Errorf("expected the marked title, got %q", title)
	}
	if text := c.curveInfoComponent.warningsComponent.GetText(); !strings.Contains(text, "⚠ curve warning") {
		t.Errorf("expected the warning in the info panel, got %q", text)
	}

	c.SetCurve(&state.CurveState{Curve: &client.Curve{Config: client.CurveConfig{ID: "curve-1"}}})

	if title := c.window.GetTitle(); title != " curve-1 " {
		t.Errorf("expected the plain title, got %q", title)
	}
	if c.curveInfoComponent.warningsComponent.HasWarnings() {
		t.Errorf("expected no warnings in the info panel, got %q", c.curveInfoComponent.warningsComponent.GetText())
	}
}
