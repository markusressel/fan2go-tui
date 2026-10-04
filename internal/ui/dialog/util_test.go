package dialog

import (
	"fmt"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Regression: tview hands mouse events to all pages until one consumes them (https://github.com/rivo/tview/issues/926),
// and the empty space around a dialog did not: a click next to it reached the page behind it.
func TestCreateModal_MouseEventsOutsideTheDialogDoNotReachThePageBehind(t *testing.T) {
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	app.SetScreen(screen)
	screen.SetSize(100, 40)
	app.EnableMouse(true)
	background := tview.NewTable().SetSelectable(true, false)
	for i := 0; i < 100; i++ {
		background.SetCell(i, 0, tview.NewTableCell(fmt.Sprintf("row %d", i)))
	}
	content := tview.NewTextView().SetText("a dialog")
	modal := createModal("Dialog", content, 40, 10)
	pages := tview.NewPages().
		AddPage("background", background, true, true).
		AddPage("dialog", modal, true, true)
	app.SetRoot(pages, true).SetFocus(background)
	go func() { _ = app.Run() }()
	defer app.Stop()

	onUiThread := func(f func()) {
		done := make(chan struct{})
		go func() {
			app.QueueUpdateDraw(f)
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for the UI thread")
		}
	}
	state := func() (row int, offset int) {
		onUiThread(func() {
			row, _ = background.GetSelection()
			offset, _ = background.GetOffset()
		})
		return row, offset
	}
	state()

	// row 2 of the table behind, at the top left, outside the centered dialog
	screen.InjectMouse(2, 2, tcell.Button1, tcell.ModNone)
	screen.InjectMouse(2, 2, tcell.ButtonNone, tcell.ModNone)
	for i := 0; i < 5; i++ {
		screen.InjectMouse(2, 2, tcell.WheelDown, tcell.ModNone)
		screen.InjectMouse(2, 2, tcell.ButtonNone, tcell.ModNone)
	}
	time.Sleep(200 * time.Millisecond)
	if row, offset := state(); row != 0 || offset != 0 {
		t.Errorf("the page behind the dialog changed: selected row %d, scrolled to %d", row, offset)
	}

	// on the border of the dialog, which its content does not cover
	var frameX, frameY, frameHeight int
	onUiThread(func() { frameX, frameY, _, frameHeight = content.GetRect() })
	time.Sleep(tview.DoubleClickInterval)
	screen.InjectMouse(frameX-1, frameY+frameHeight/2, tcell.Button1, tcell.ModNone)
	screen.InjectMouse(frameX-1, frameY+frameHeight/2, tcell.ButtonNone, tcell.ModNone)
	time.Sleep(200 * time.Millisecond)
	if row, _ := state(); row != 0 {
		t.Errorf("a click on the border of the dialog selected row %d of the page behind it", row)
	}

	// on the dialog, its content gets the click
	time.Sleep(tview.DoubleClickInterval)
	var x, y int
	onUiThread(func() { x, y, _, _ = content.GetRect() })
	screen.InjectMouse(x+1, y, tcell.Button1, tcell.ModNone)
	screen.InjectMouse(x+1, y, tcell.ButtonNone, tcell.ModNone)
	time.Sleep(200 * time.Millisecond)
	var focused bool
	onUiThread(func() { focused = content.HasFocus() })
	if !focused {
		t.Error("a click on the dialog does not reach its content")
	}
}
