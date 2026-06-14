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
