package util

import (
	"fan2go-tui/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"math"
	"sync"
)

type ListComponent[T comparable] struct {
	application *tview.Application

	layout        *tview.Flex
	entriesLayout *tview.Flex

	entries      []*T
	entriesMutex sync.Mutex

	config             *ListComponentConfig
	entryVisibilityMap map[*T]bool
	selectedIndex      int

	getLayout                func(entry *T) (layout *tview.Flex)
	inputCapture             func(event *tcell.EventKey) *tcell.EventKey
	selectionChangedCallback func(selectedEntry *T)

	sortListEntries func(entries []*T, inverted bool) []*T

	//compare      func(a, b *T) bool
	sortInverted bool

	scrollbarComponent *ScrollbarComponent
}

// NewListComponent creates a new ListComponent.
// The application is used to redraw the component.
// The config is used to configure the component.
// The getLayout function is used to create the layout for each entry.
// The compare function is used to sort the entries.
// The sortListEntries function is used to sort the entries.
func NewListComponent[T comparable](
	application *tview.Application,
	config *ListComponentConfig,
	getLayout func(entry *T) (layout *tview.Flex),
//compare func(a, b *T) bool,
	sortListEntries func(entries []*T, inverted bool) []*T,
) *ListComponent[T] {
	listComponent := &ListComponent[T]{
		application:        application,
		config:             config,
		entries:            []*T{},
		entriesMutex:       sync.Mutex{},
		entryVisibilityMap: map[*T]bool{},
		getLayout:          getLayout,
		sortListEntries:    sortListEntries,
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		selectionChangedCallback: func(selectedEntry *T) {},
		//compare:                  compare,
		selectedIndex: -1,
	}
	listComponent.createLayout()
	listComponent.SetDirection(tview.FlexColumn)
	return listComponent
}

func (c *ListComponent[T]) createLayout() {
	layout := tview.NewFlex()

	c.entriesLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(c.entriesLayout, 0, 1, true)

	c.scrollbarComponent = NewScrollbarComponent(c.application, ScrollBarVertical, 0, 1, 0, c.config.MaxVisibleItems)
	layout.AddItem(c.scrollbarComponent.GetLayout(), 1, 0, false)

	c.entriesLayout.SetFocusFunc(func() {
		// ensure the first item is automatically selected, if there is any
		data := c.GetData()
		if data != nil && len(data) > 0 {
			layout.Blur()
			if c.selectedIndex == -1 {
				c.selectedIndex = 0
			}
			itemLayout := c.getLayout(data[0])

			c.SelectEntry(c.GetSelectedItem())
			c.application.SetFocus(itemLayout)
		}
	})

	c.entriesLayout.Focus(func(item tview.Primitive) {
		for idx, entry := range c.entries {
			if item == c.getLayout(entry) {
				c.selectedIndex = idx
			}
		}
	})

	c.entriesLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
	c.updateScrollBar()
	c.application.ForceDraw()
}

func (c *ListComponent[T]) GetLayout() *tview.Flex {
	return c.layout
}

func (c *ListComponent[T]) SetTitle(title string) {
	SetupWindow(c.layout, title)
}

func (c *ListComponent[T]) GetData() []*T {
	return c.sortListEntries(c.entries, c.sortInverted)
	//sort.SliceStable(c.entries, func(i, j int) bool {
	//	a := c.entries[i]
	//	b := c.entries[j]
	//	return c.compare(a, b)
	//})
	//return c.entries
}

func (c *ListComponent[T]) SetData(entries []*T) {
	c.entriesMutex.Lock()
	selectFirst := c.entries == nil
	c.entries = entries
	c.entriesMutex.Unlock()
	c.updateLayout()
	c.application.ForceDraw()
	if selectFirst {
		c.SelectFirst()
	}
}

func (c *ListComponent[comparable]) SortBy(inverted bool) {
	c.entriesMutex.Lock()
	c.sortInverted = inverted
	c.entries = c.sortListEntries(c.entries, c.sortInverted)
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

	if len(entryVisibilityMapValues) > 0 && rows < 0 && entryVisibilityMapValues[0] == false || rows > 0 && entryVisibilityMapValues[len(entryVisibilityMapValues)-1] == false {
		entryVisibilityMapValues = util.RotateSliceBy(entryVisibilityMapValues, rows)
	}

	c.entryVisibilityMap = map[*T]bool{}
	for i, key := range entryVisibilityMapKeys {
		c.entryVisibilityMap[key] = entryVisibilityMapValues[i]
	}
	c.updateLayout()
}

func (c *ListComponent[T]) GetVisibleRange() (int, int) {
	minIndex := len(c.entryVisibilityMap)
	maxIndex := 0
	for idx, entry := range c.entries {
		isVisible := c.entryVisibilityMap[entry]
		if isVisible {
			minIndex = int(math.Min(float64(minIndex), float64(idx)))
			maxIndex = int(math.Max(float64(maxIndex), float64(idx)))
		}
	}
	return minIndex, maxIndex
}

func (c *ListComponent[T]) updateVisibleEntries() {
	// ensure we are displaying as many items as specified by MaxVisibleItems

	// cleanup the visibility map (remove entries that are not in the dataset anymore)
	c.cleanupVisibilityMap()

	for _, entry := range c.entries {
		_, ok := c.entryVisibilityMap[entry]
		if !ok {
			if c.getVisibleEntriesCount() < c.config.MaxVisibleItems {
				c.entryVisibilityMap[entry] = true
			} else {
				c.entryVisibilityMap[entry] = false
			}
		}
	}

	c.entriesLayout.Clear()
	for _, entry := range c.entries {
		currentVisibility := c.entryVisibilityMap[entry]
		if currentVisibility {
			c.entriesLayout.AddItem(c.getLayout(entry), 0, 1, false)
		}
	}
}

func (c *ListComponent[T]) cleanupVisibilityMap() {
	keys := maps.Keys(c.entryVisibilityMap)
	for _, key := range keys {
		ok := slices.Contains(c.entries, key)
		if !ok {
			delete(c.entryVisibilityMap, key)
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

func (c *ListComponent[T]) SelectEntry(entry *T) {
	indexToSelect := slices.Index(c.entries, entry)
	if indexToSelect == -1 {
		return
	}
	entryToSelect := c.entries[indexToSelect]
	entryLayout := c.getLayout(entryToSelect)
	c.selectedIndex = indexToSelect
	c.application.SetFocus(entryLayout)
	c.selectionChangedCallback(entry)
	c.scrollTo(entryToSelect)
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
		entryLayout := c.getLayout(entry)
		if entryLayout.HasFocus() {
			nextEntryIndex := (len(c.entries) + idx + rows) % len(c.entries)
			nextEntry := c.entries[nextEntryIndex]
			nextEntryLayout := c.getLayout(nextEntry)
			c.selectedIndex = nextEntryIndex
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
	for idx, entry := range c.GetData() {
		entryLayout := c.getLayout(entry)
		c.application.SetFocus(entryLayout)
		c.selectedIndex = idx
		c.selectionChangedCallback(entry)
		return
	}
	c.scrollToSelection()
}

func (c *ListComponent[T]) updateScrollBar() {
	if len(c.entries) <= c.config.MaxVisibleItems {
		c.hideScrollbar()
	} else {
		c.showScrollbar()
	}

	c.scrollbarComponent.SetMin(0)
	c.scrollbarComponent.SetMax(int(math.Max(0.0, float64(len(c.entries)))))
	visibleIndexMin, visibleIndexMax := c.GetVisibleRange()

	//newPosition := c.GetSelectedIndex()
	c.scrollbarComponent.SetPosition(visibleIndexMin)

	width := (visibleIndexMax - visibleIndexMin) + 1
	c.scrollbarComponent.SetWidth(width)
}

func (c *ListComponent[T]) GetSelectedIndex() int {
	return c.selectedIndex
}

func (c *ListComponent[T]) GetSelectedItem() *T {
	if c.selectedIndex == -1 {
		return nil
	}
	return c.entries[c.selectedIndex]
}

func (c *ListComponent[T]) hideScrollbar() {
	c.layout.RemoveItem(c.scrollbarComponent.GetLayout())
}

func (c *ListComponent[T]) showScrollbar() {
	if c.layout.GetItemCount() <= 1 {
		c.layout.AddItem(c.scrollbarComponent.GetLayout(), 1, 0, false)
	}
}
