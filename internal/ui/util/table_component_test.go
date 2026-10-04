package util

import (
	"sort"
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

func TestRowSelectionTable_SortSelectionAndData(t *testing.T) {
	app := tview.NewApplication()

	col1 := &Column{Id: 1, Title: "Name", Alignment: tview.AlignLeft}
	col2 := &Column{Id: 2, Title: "Value", Alignment: tview.AlignRight}
	columns := []*Column{col1, col2}

	sortTable := func(entries []*string, column *Column, inverted bool) []*string {
		res := append([]*string{}, entries...)
		sort.Slice(res, func(i, j int) bool {
			if inverted {
				return *res[i] > *res[j]
			}
			return *res[i] < *res[j]
		})
		return res
	}

	toCells := func(row int, cols []*Column, entry *string) []*tview.TableCell {
		return []*tview.TableCell{
			tview.NewTableCell(*entry),
			tview.NewTableCell("val"),
		}
	}

	table := NewTableContainer[string](app, toCells, sortTable)
	table.SetTitle("Test Table")
	table.SetColumnSpec(columns, col1, false)

	if !table.IsEmpty() {
		t.Errorf("expected table to be empty")
	}

	a, b, c := "apple", "banana", "cherry"
	table.SetData([]*string{&c, &a, &b})

	if table.IsEmpty() {
		t.Errorf("expected table not empty")
	}
	entries := table.GetEntries()
	if len(entries) != 3 || *entries[0] != "apple" {
		t.Fatalf("expected sorted entries starting with apple, got %v", *entries[0])
	}

	// Selection
	var selected *string
	table.SetSelectionChangedCallback(func(entry *string) {
		selected = entry
	})

	table.SelectFirstIfExists()
	if selected == nil || *selected != "apple" {
		t.Errorf("expected apple selected, got %v", selected)
	}
	if table.GetSelectedEntry() != entries[0] {
		t.Errorf("expected GetSelectedEntry to return apple")
	}

	// Select cherry
	table.Select(&c)
	if table.GetSelectedEntry() != &c {
		t.Errorf("expected cherry selected")
	}

	// Select unknown entry
	unknown := "dragonfruit"
	table.Select(&unknown)
	if table.GetSelectedEntry() != &c {
		t.Errorf("expected selection unchanged on unknown entry")
	}

	// Select header (row 0)
	table.SelectHeader()
	if table.GetSelectedEntry() != nil {
		t.Errorf("expected nil selected entry when header selected, got %v", table.GetSelectedEntry())
	}

	// Test sort navigation via key inputs when header selected
	capture := table.GetLayout().GetInputCapture()
	if capture != nil {
		// KeyRight triggers nextSortOrder
		capture(tcell.NewEventKey(tcell.KeyRight, ' ', tcell.ModNone))
		if table.sortByColumn != col2 {
			t.Errorf("expected sortByColumn to be col2 after KeyRight")
		}

		// KeyLeft triggers previousSortOrder
		capture(tcell.NewEventKey(tcell.KeyLeft, ' ', tcell.ModNone))
		if table.sortByColumn != col1 {
			t.Errorf("expected sortByColumn to be col1 after KeyLeft")
		}

		// KeyEnter toggles sort direction
		capture(tcell.NewEventKey(tcell.KeyEnter, ' ', tcell.ModNone))
		if !table.sortInverted {
			t.Errorf("expected sort to be inverted after Enter")
		}
		entries = table.GetEntries()
		if *entries[0] != "cherry" {
			t.Errorf("expected cherry first after invert sort, got %v", *entries[0])
		}
	}

	// Focus check
	_ = table.HasFocus()
}
