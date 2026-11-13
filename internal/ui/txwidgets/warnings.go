package txwidgets

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
	"strings"

	"github.com/rivo/tview"
)

// WarningMarker marks a fan, sensor or curve with warnings, and each of its warning messages.
const WarningMarker = "⚠"

// SetWindowWarnings sets the title of the window of the fan, sensor or curve with the given ID: with WarningMarker
// in the warning color if it has warnings, the plain ID (like util.SetupWindow) otherwise.
func SetWindowWarnings(window *tview.Box, id string, warnings []client.Warning) {
	if len(warnings) == 0 {
		window.SetTitle(theme.CreateTitleText(id))
		window.SetTitleColor(theme.Colors.Layout.Title)
		return
	}
	window.SetTitle(theme.CreateTitleText(WarningMarker + " " + id))
	window.SetTitleColor(theme.Colors.Warning.Title)
}

// WarningsText returns the warning messages, one paragraph per warning, each starting with WarningMarker.
func WarningsText(warnings []client.Warning) string {
	paragraphs := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		paragraphs = append(paragraphs, Span(theme.Colors.Warning.Text, "%s %s", WarningMarker, warning.Message))
	}
	return strings.Join(paragraphs, "\n")
}

// WarningsComponent shows the warning messages of a fan, sensor or curve, wrapped to its width. The info components
// hide it (height 0) while there are no warnings, see HasWarnings.
type WarningsComponent struct {
	layout *tview.TextView

	text string
}

func NewWarningsComponent() *WarningsComponent {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(true).
		SetWordWrap(true).
		SetScrollable(true)
	return &WarningsComponent{layout: textView}
}

func (c *WarningsComponent) GetPrimitive() tview.Primitive {
	return c.layout
}

// SetWarnings shows the given warnings. The text is only replaced if it changed, so the scroll position is kept
// across refreshes.
func (c *WarningsComponent) SetWarnings(warnings []client.Warning) {
	text := WarningsText(warnings)
	if text == c.text {
		return
	}
	c.text = text
	c.layout.SetText(text)
	c.layout.ScrollToBeginning()
}

// HasWarnings returns whether there are warnings to show.
func (c *WarningsComponent) HasWarnings() bool {
	return c.text != ""
}

// GetText returns the shown text without color tags (for tests).
func (c *WarningsComponent) GetText() string {
	return c.layout.GetText(true)
}
