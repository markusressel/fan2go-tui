package util

import (
	"fan2go-tui/internal/testutil"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Regression: the size was polled from another goroutine, while the UI thread drew the box (a data race, found by
// "go test -race"). It is checked while drawing now.
func TestSetupReactiveResize(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	app := tview.NewApplication().SetScreen(screen)
	box := tview.NewFlex()
	box.SetBorder(true)

	resizes := 0
	var innerWidth, innerHeight int
	SetupReactiveResize(app, box.Box, func() {
		resizes++
		_, _, innerWidth, innerHeight = box.GetInnerRect()
	})
	app.SetRoot(box, true)
	testutil.RunApplication(t, app)

	testutil.WaitFor(t, app, "resize after the first draw", func() bool { return resizes == 1 })
	width, height := screen.Size()
	if innerWidth != width-2 || innerHeight != height-2 {
		t.Errorf("expected the inner size %dx%d (without the border), got %dx%d", width-2, height-2, innerWidth, innerHeight)
	}

	testutil.OnUiThread(t, app, func() {}) // another draw of the same size
	if resizes != 1 {
		t.Errorf("expected no resize while the size is the same, got %d", resizes)
	}

	screen.SetSize(width+10, height+5)
	testutil.WaitFor(t, app, "resize after the screen size changed", func() bool { return resizes == 2 })
	if innerWidth != width+8 || innerHeight != height+3 {
		t.Errorf("expected the inner size %dx%d, got %dx%d", width+8, height+3, innerWidth, innerHeight)
	}

	// draws and size changes for a while, for the race detector
	for i := range 20 {
		screen.SetSize(width+i%2, height)
		testutil.OnUiThread(t, app, func() {})
		time.Sleep(10 * time.Millisecond)
	}
}
