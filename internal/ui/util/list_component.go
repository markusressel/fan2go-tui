package util

import (
	"fan2go-tui/internal/logging"
	"fan2go-tui/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
	"math"
	"sort"
	"sync"
)

const (
	MaxVisibleItems = 3
)

type IListComponent interface {
	Compare()
}

type ListComponent[T comparable] struct {
	application *tview.Application

	layout *tview.Flex

	entries      []*T
	entriesMutex sync.Mutex

	entryVisibilityMap map[*T]bool

	toLayout                 func(entry *T) (layout tview.Primitive)
	inputCapture             func(event *tcell.EventKey) *tcell.EventKey
	selectionChangedCallback func(selectedEntry *T)

	compare      func(a, b *T) bool
	sortInverted bool
}

func NewListComponent[T comparable](
	application *tview.Application,
	toLayout func(entry *T) (layout tview.Primitive),
	compare func(a, b *T) bool,
) *ListComponent[T] {
	listComponent := &ListComponent[T]{
		application:        application,
		entries:            nil,
		entriesMutex:       sync.Mutex{},
		entryVisibilityMap: map[*T]bool{},
		toLayout:           toLayout,
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		selectionChangedCallback: func(selectedEntry *T) {},
		compare:                  compare,
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
			c.selectPreviousEntry()
			return event
		} else if key == tcell.KeyDown {
			c.selectNextEntry()
			return event
		} else if key == tcell.KeyLeft {
			return nil
		} else if key == tcell.KeyRight {
			return nil
		} else if key == tcell.KeyEnter {
			return nil
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
	sort.SliceStable(c.entries, func(i, j int) bool {
		a := c.entries[i]
		b := c.entries[j]
		return c.compare(a, b)
	})
	return c.entries
}

func (c *ListComponent[T]) SetData(entries []*T) {
	c.entriesMutex.Lock()
	selectFirst := c.entries == nil
	c.entries = entries
	c.entriesMutex.Unlock()
	c.updateLayout()
	if selectFirst {
		c.SelectFirst()
	}
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

func (c *ListComponent[T]) scrollUp() {
	c.scroll(-1)
	c.selectPreviousEntry()
	c.application.ForceDraw()
}

func (c *ListComponent[T]) scrollDown() {
	c.scroll(+1)
	c.selectNextEntry()
	c.application.ForceDraw()
}

func (c *ListComponent[T]) scroll(rows int) {
	var entryVisibilityMapKeys []*T
	var entryVisibilityMapValues []bool

	var keys = c.GetData()

	for _, key := range keys {
		value := c.entryVisibilityMap[key]
		entryVisibilityMapKeys = append(entryVisibilityMapKeys, key)
		entryVisibilityMapValues = append(entryVisibilityMapValues, value)
	}

	if rows < 0 && entryVisibilityMapValues[0] == false || rows > 0 && entryVisibilityMapValues[len(entryVisibilityMapValues)-1] == false {
		entryVisibilityMapValues = util.RotateSliceBy(entryVisibilityMapValues, rows)
	}

	c.entryVisibilityMap = map[*T]bool{}
	for i, key := range entryVisibilityMapKeys {
		c.entryVisibilityMap[key] = entryVisibilityMapValues[i]
	}
	c.updateLayout()
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
	for _, entry := range c.entries {
		currentVisibility := c.entryVisibilityMap[entry]
		if currentVisibility {
			c.layout.AddItem(c.toLayout(entry), 0, 1, false)
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

func (c *ListComponent[T]) selectPreviousEntry() {
	newSelection := c.shiftSelection(-1)
	c.scrollTo(newSelection)
	c.application.ForceDraw()
}

func (c *ListComponent[T]) selectNextEntry() {
	newSelection := c.shiftSelection(+1)
	c.scrollTo(newSelection)
	c.application.ForceDraw()
}

func (c *ListComponent[T]) shiftSelection(rows int) *T {
	for idx, entry := range c.entries {
		entryLayout := c.toLayout(entry)
		if entryLayout.HasFocus() {
			nextEntryIndex := (len(c.entries) + idx + rows) % len(c.entries)
			nextEntry := c.entries[nextEntryIndex]
			nextEntryLayout := c.toLayout(nextEntry)
			c.application.SetFocus(nextEntryLayout)
			c.selectionChangedCallback(entry)
			return nextEntry
		}
	}
	return nil
}

func (c *ListComponent[T]) scrollTo(selection *T) {
	if !c.isInVisibleScrollRange(selection) {
		distance := c.determineScrollDistanceToEntry(selection)
		c.scroll(distance)
	}
}

func (c *ListComponent[T]) scrollToSelection() {
	newSelection := c.shiftSelection(0)
	c.scrollTo(newSelection)
}

func (c *ListComponent[T]) isInVisibleScrollRange(selection *T) bool {
	for _, entry := range c.GetData() {
		isVisible := c.entryVisibilityMap[entry]
		if entry == selection && isVisible {
			return true
		}
	}
	return false
}

func (c *ListComponent[T]) determineScrollDistanceToEntry(selection *T) int {
	data := c.GetData()

	index := slices.Index(data, selection)

	// find the min/max indices of currently visible items
	minIndex := len(c.entryVisibilityMap)
	maxIndex := 0
	for idx, entry := range data {
		isVisible := c.entryVisibilityMap[entry]
		if isVisible {
			minIndex = int(math.Min(float64(minIndex), float64(idx)))
			maxIndex = int(math.Max(float64(maxIndex), float64(idx)))
		}
	}

	if index < minIndex {
		return index - minIndex
	} else if index > maxIndex {
		return index - maxIndex
	} else {
		return 0
	}
}

func (c *ListComponent[T]) SelectFirst() {
	for _, entry := range c.GetData() {
		entryLayout := c.toLayout(entry)
		c.application.SetFocus(entryLayout)
		focused := c.application.GetFocus()
		if focused != entryLayout {
			logging.Error("sdsadsa")
		}
		c.selectionChangedCallback(entry)
		return
	}
	c.scrollToSelection()
}
