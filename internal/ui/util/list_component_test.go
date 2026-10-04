package util

import (
	"testing"

	"github.com/rivo/tview"
)

func TestListComponent_ScrollingAndVisibility(t *testing.T) {
	app := tview.NewApplication()
	// Use a fixed MaxVisibleItems to avoid reliance on terminal height
	config := NewListComponentConfig().WithMaxVisibleItems(5)

	// Sample data
	entries := make([]*int, 20)
	for i := 0; i < 20; i++ {
		val := i
		entries[i] = &val
	}

	list := NewListComponent[int](
		app,
		config,
		func(entry *int) *tview.Flex { return tview.NewFlex() },
		func(entries []*int, inverted bool) []*int { return entries },
	)

	t.Run("InitialState", func(t *testing.T) {
		list.SetData(entries)
		min, max := list.GetVisibleRange()
		if min != 0 || max != 4 {
			t.Errorf("Expected initial visible range [0, 4], got [%d, %d]", min, max)
		}
	})

	t.Run("ScrollDown", func(t *testing.T) {
		list.scroll(2)
		min, max := list.GetVisibleRange()
		if min != 2 || max != 6 {
			t.Errorf("Expected visible range [2, 6] after scrolling down, got [%d, %d]", min, max)
		}
	})

	t.Run("ScrollToBottom", func(t *testing.T) {
		list.scroll(100) // Scroll way past bottom
		min, max := list.GetVisibleRange()
		// Should clamp to len(entries) - maxVisible = 20 - 5 = 15
		if min != 15 || max != 19 {
			t.Errorf("Expected clamped visible range [15, 19], got [%d, %d]", min, max)
		}
	})

	t.Run("ScrollToTop", func(t *testing.T) {
		list.scroll(-100) // Scroll way past top
		min, max := list.GetVisibleRange()
		if min != 0 || max != 4 {
			t.Errorf("Expected clamped visible range [0, 4], got [%d, %d]", min, max)
		}
	})

	t.Run("PageDown", func(t *testing.T) {
		list.startIndex = 0
		list.scrollByPage(1)
		min, max := list.GetVisibleRange()
		// Page size is 5, so startIndex should move from 0 to 5
		if min != 5 || max != 9 {
			t.Errorf("Expected range [5, 9] after PageDown, got [%d, %d]", min, max)
		}
	})

	t.Run("ScrollToSpecificEntry", func(t *testing.T) {
		target := entries[12]
		list.scrollTo(target)
		min, max := list.GetVisibleRange()
		// Entry 12 should be visible. scrollTo usually puts it at the edge of the window.
		if 12 < min || 12 > max {
			t.Errorf("Target entry 12 not in visible range [%d, %d]", min, max)
		}
	})
}

func TestListComponent_SelectionPreservedAcrossSetData(t *testing.T) {
	app := tview.NewApplication()
	config := NewListComponentConfig().WithMaxVisibleItems(5)

	entries := make([]*int, 10)
	for i := 0; i < 10; i++ {
		val := i
		entries[i] = &val
	}

	list := NewListComponent[int](
		app,
		config,
		func(entry *int) *tview.Flex { return tview.NewFlex() },
		func(entries []*int, inverted bool) []*int { return entries },
	)

	list.SetData(entries)

	// Initially select index 3
	list.SelectEntry(entries[3])
	if list.GetSelectedIndex() != 3 {
		t.Fatalf("expected selected index 3, got %d", list.GetSelectedIndex())
	}
	if list.GetSelectedItem() != entries[3] {
		t.Fatalf("expected selected item %v, got %v", entries[3], list.GetSelectedItem())
	}

	// Refresh with same data
	list.SetData(entries)
	if list.GetSelectedIndex() != 3 {
		t.Errorf("expected selected index to remain 3 after SetData, got %d", list.GetSelectedIndex())
	}
	if list.GetSelectedItem() != entries[3] {
		t.Errorf("expected selected item to remain %v after SetData, got %v", entries[3], list.GetSelectedItem())
	}

	// Next entry via shiftSelection
	list.shiftSelection(1)
	if list.GetSelectedIndex() != 4 {
		t.Errorf("expected selected index 4 after shiftSelection, got %d", list.GetSelectedIndex())
	}

	// Refresh again
	list.SetData(entries)
	if list.GetSelectedIndex() != 4 {
		t.Errorf("expected selected index to remain 4 after SetData, got %d", list.GetSelectedIndex())
	}

	// Triggering entriesLayout focus callback should not reset to 0
	list.entriesLayout.Focus(func(p tview.Primitive) {})
	if list.GetSelectedIndex() != 4 {
		t.Errorf("expected selected index to remain 4 after focus, got %d", list.GetSelectedIndex())
	}
}

