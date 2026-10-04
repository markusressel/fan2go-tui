package util

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Regression: mouse captures get all mouse events, not only those on their primitive
// (https://github.com/rivo/tview/issues/926), so double clicks anywhere called the callback of the table.
func TestRowSelectionTable_DoubleClickOnlyOnTheTable(t *testing.T) {
	table := NewTableContainer[string](tview.NewApplication(),
		func(row int, columns []*Column, entry *string) []*tview.TableCell { return nil },
		func(entries []*string, column *Column, inverted bool) []*string { return entries })
	doubleClicks := 0
	table.SetDoubleClickCallback(func() { doubleClicks++ })
	table.GetLayout().SetRect(10, 10, 20, 10)
	capture := table.GetLayout().GetMouseCapture()

	capture(tview.MouseLeftDoubleClick, tcell.NewEventMouse(50, 5, tcell.Button1, tcell.ModNone))
	if doubleClicks != 0 {
		t.Errorf("a double click outside of the table called its callback")
	}
	_, event := capture(tview.MouseLeftDoubleClick, tcell.NewEventMouse(15, 12, tcell.Button1, tcell.ModNone))
	if doubleClicks != 1 || event != nil {
		t.Errorf("a double click on the table: %d calls, event consumed: %v", doubleClicks, event == nil)
	}
}
