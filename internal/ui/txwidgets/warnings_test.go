package txwidgets

import (
	"fan2go-tui/internal/client"
	"testing"

	"github.com/rivo/tview"
)

func TestSetWindowWarnings(t *testing.T) {
	window := tview.NewBox()

	SetWindowWarnings(window, "cpu", []client.Warning{{Code: "flatRpmCurve", Message: "flat"}})
	if title := window.GetTitle(); title != " ⚠ cpu " {
		t.Errorf("expected the marked title, got %q", title)
	}

	// the warning is gone
	SetWindowWarnings(window, "cpu", nil)
	if title := window.GetTitle(); title != " cpu " {
		t.Errorf("expected the plain title, got %q", title)
	}
}

func TestWarningsComponent(t *testing.T) {
	c := NewWarningsComponent()
	if c.HasWarnings() || c.GetText() != "" {
		t.Fatalf("expected no warnings initially, got %q", c.GetText())
	}

	c.SetWarnings([]client.Warning{
		{Code: "flatRpmCurve", Message: "RPM stays at ~750 [not a color tag]"},
		{Code: "other", Message: "second"},
	})
	if !c.HasWarnings() {
		t.Fatalf("expected warnings")
	}
	// the messages from the daemon are escaped
	if expected := "⚠ RPM stays at ~750 [not a color tag]\n⚠ second"; c.GetText() != expected {
		t.Errorf("expected %q, got %q", expected, c.GetText())
	}

	c.SetWarnings([]client.Warning{})
	if c.HasWarnings() || c.GetText() != "" {
		t.Errorf("expected no warnings after they are gone, got %q", c.GetText())
	}
}
