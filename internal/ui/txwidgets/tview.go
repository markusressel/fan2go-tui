package txwidgets

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func Span(color tcell.Color, format string, args ...any) string {
	return ColorTag(color) + tview.Escape(fmt.Sprintf(format, args...)) + "[-]"
}

func ColorTag(color tcell.Color) string {
	r, g, b := color.RGB()
	return fmt.Sprintf("[#%02x%02x%02x]", uint8(r), uint8(g), uint8(b))
}
