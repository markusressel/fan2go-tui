package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"sync"
)

const (
	MaxVisibleItems = 3
)

type ListComponent[T comparable] struct {
	application *tview.Application

	layout *tview.Flex

	entries      []*T
	entriesMutex sync.Mutex

	entryVisibilityMap map[*T]bool

	toLayout                 func(row int, entry *T) (layout tview.Primitive)
	inputCapture             func(event *tcell.EventKey) *tcell.EventKey
	selectionChangedCallback func(selectedEntry *T)

	sortInverted bool
}

func NewListComponent[T comparable](
	application *tview.Application,
	toLayout func(row int, entry *T) (layout tview.Primitive),
) *ListComponent[T] {
	listComponent := &ListComponent[T]{
		application:        application,
		entries:            []*T{},
		entriesMutex:       sync.Mutex{},
		entryVisibilityMap: map[*T]bool{},
		toLayout:           toLayout,
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		selectionChangedCallback: func(selectedEntry *T) {},
	}
	listComponent.createLayout()
	listComponent.SetDirection(tview.FlexRow)
	return listComponent
}

func (c *ListComponent[T]) createLayout() {
	layout := tview.NewFlex()

	SetupWindow(layout, "")

	layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = c.inputCapture(event)
		if event == nil {
			return event
		}
		key := event.Key()
		if key == tcell.KeyUp {
			return nil
		} else if key == tcell.KeyDown {
			return nil
		} else if key == tcell.KeyLeft {
			return nil
		} else if key == tcell.KeyRight {
			return nil
		} else if key == tcell.KeyEnter {
			return nil
		}

		if key == tcell.KeyUp {
			c.ShiftFocusUp()
			//FocusPreviousItem()
			return nil
		} else if key == tcell.KeyDown {
			c.ShiftFocusDown()
			return nil
		} else if key == tcell.KeyEnter {
			//selectedEntry := c.GetSelectedEntry()
			//c.selectionChangedCallback(selectedEntry)
			//return nil
		}

		return event
	})

	c.layout = layout
}

func (c *ListComponent[T]) SetDirection(direction int) {
	c.layout.SetDirection(direction)
}

func (c *ListComponent[T]) updateLayout() {
	c.updateVisibleEntries()
	c.application.ForceDraw()
}

func (c *ListComponent[T]) GetLayout() *tview.Flex {
	return c.layout
}

func (c *ListComponent[T]) SetTitle(title string) {
	SetupWindow(c.layout, title)
}

func (c *ListComponent[T]) GetData() []*T {
	return c.entries
}

func (c *ListComponent[T]) SetData(entries []*T) {
	c.entriesMutex.Lock()
	c.entries = entries
	c.entriesMutex.Unlock()
	c.updateLayout()
}

func (c *ListComponent[comparable]) SortBy(inverted bool) {
	c.entriesMutex.Lock()
	c.sortInverted = inverted
	// c.entries = c.sortTableEntries(c.entries, c.sortByColumn, c.sortInverted)
	c.entriesMutex.Unlock()
}

func (c *ListComponent[abc]) HasFocus() bool {
	return c.layout.HasFocus()
}

func (c *ListComponent[T]) GetEntries() []*T {
	return c.entries
}

func (c *ListComponent[T]) IsEmpty() bool {
	return len(c.entries) <= 0
}

func (c *ListComponent[T]) SetInputCapture(inputCapture func(event *tcell.EventKey) *tcell.EventKey) {
	c.inputCapture = inputCapture
}

func (c *ListComponent[T]) SetSelectionChangedCallback(f func(selectedEntry *T)) {
	c.selectionChangedCallback = f
}

func (c *ListComponent[T]) ShiftFocusUp() {
	c.shiftFocus(-1)
}

func (c *ListComponent[T]) ShiftFocusDown() {
	c.shiftFocus(+1)
}

func (c *ListComponent[T]) shiftFocus(amount int) {
	entryVisibilityMapKeys := maps.Keys(c.entryVisibilityMap)
	entryVisibilityMapValues := maps.Values(c.entryVisibilityMap)

	c.entryVisibilityMap = map[*T]bool{}
	for i, key := range entryVisibilityMapKeys {
		c.entryVisibilityMap[key] = entryVisibilityMapValues[i+amount%len(entryVisibilityMapValues)]
	}
}

func (c *ListComponent[T]) updateVisibleEntries() {
	// ensure we are displaying as many items as specified by MaxVisibleItems
	for _, entry := range c.entries {
		_, ok := c.entryVisibilityMap[entry]
		if !ok {
			if c.getVisibleEntriesCount() < MaxVisibleItems {
				c.entryVisibilityMap[entry] = true
			} else {
				c.entryVisibilityMap[entry] = false
			}
		}
	}

	c.layout.Clear()
	for row, entry := range c.entries {
		currentVisibility := c.entryVisibilityMap[entry]
		if currentVisibility {
			c.layout.AddItem(c.toLayout(row, entry), 0, 1, false)
		}
	}
}

func (c *ListComponent[T]) getVisibleEntriesCount() int {
	count := 0
	for _, isVisible := range c.entryVisibilityMap {
		if isVisible {
			count += 1
		}
	}
	return count
}
