package fan

import (
	"fan2go-tui/internal/demo"
	"strconv"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func screenContent(screen tcell.SimulationScreen) string {
	cells, width, height := screen.GetContents()
	var content strings.Builder
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if runes := cells[y*width+x].Runes; len(runes) > 0 {
				content.WriteRune(runes[0])
			} else {
				content.WriteRune(' ')
			}
		}
		content.WriteRune('\n')
	}
	return content.String()
}

// Regression test: the y-axis was drawn with and without its labels in turns on every update at some widths,
// because the points of the curve were only shown if a column of the plot matched them exactly, which changed
// the scale of the y-axis, the width of its labels and so the columns of the plot.
func TestFanRpmCurveComponent_LayoutIsStableAtAllWidths(t *testing.T) {
	fan := demo.NewSimulator().GetFans()["gpu-fan"]

	for width := 40; width <= 220; width++ {
		screen := tcell.NewSimulationScreen("UTF-8")
		if err := screen.Init(); err != nil {
			t.Fatal(err)
		}
		app := tview.NewApplication()
		app.SetScreen(screen)
		screen.SetSize(width, 22)
		component := NewFanRpmCurveComponent(app, fan)
		app.SetRoot(component.GetLayout(), true)

		// several draws per update, like the application does
		component.SetFan(fan)
		app.ForceDraw()
		first := screenContent(screen)
		for i := 0; i < 3; i++ {
			component.SetFan(fan)
			app.ForceDraw()
			app.ForceDraw()
			if content := screenContent(screen); content != first {
				t.Fatalf("width %d: the graph changed without a change of the data:\n%s\n%s", width, first, content)
			}
		}
		// the scale of the y-axis goes up to the highest point of the curve (2600 rpm)
		if highest := highestYAxisLabel(first); highest < 2000 {
			t.Fatalf("width %d: expected the y-axis to go up to the curve, got %d:\n%s", width, highest, first)
		}
		screen.Fini()
	}
}

// highestYAxisLabel returns the highest value labeled on the y-axis.
func highestYAxisLabel(content string) int {
	highest := -1
	for _, line := range strings.Split(content, "\n") {
		label, _, found := strings.Cut(line, "┤")
		if !found {
			continue
		}
		if value, err := strconv.Atoi(strings.TrimSpace(label)); err == nil {
			highest = max(highest, value)
		}
	}
	return highest
}
