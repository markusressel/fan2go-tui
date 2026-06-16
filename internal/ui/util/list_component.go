package util

import (
	"fan2go-tui/internal/util"
	"math"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
)

type ListComponent[T comparable] struct {
	application *tview.Application

	layout        *tview.Flex
	entriesLayout *tview.Flex

	entries      []*T
	entriesMutex sync.Mutex

	config             *ListComponentConfig
	entryVisibilityMap map[*T]bool
	startIndex         int
	selectedIndex      int

	getLayout                func(entry *T) (layout *tview.Flex)
	inputCapture             func(event *tcell.EventKey) *tcell.EventKey
	selectionChangedCallback func(selectedEntry *T)

	sortListEntries func(entries []*T, inverted bool) []*T

	sortInverted bool

	scrollbarComponent *ScrollbarComponent

	lastKnownHeight int
}

type HorizontalScrollable interface {
	ScrollHorizontal(delta int)
}

// NewListComponent creates a new ListComponent.
// The application is used to redraw the component.
// The config is used to configure the component.
// The getLayout function is used to create the layout for each entry.
// The sortListEntries function is used to sort the entries.
func NewListComponent[T comparable](
	application *tview.Application,
	config *ListComponentConfig,
	getLayout func(entry *T) (layout *tview.Flex),
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
	layout.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		_, _, _, innerHeight := layout.GetInnerRect()
		if innerHeight > 0 && innerHeight != c.lastKnownHeight {
			c.lastKnownHeight = innerHeight
			c.updateLayoutInternal()
		}
		return layout.GetInnerRect()
	})

	c.entriesLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	layout.AddItem(c.entriesLayout, 0, 1, true)

	c.scrollbarComponent = NewScrollbarComponent(c.application, ScrollBarVertical, 0, 1, 0, c.GetMaxVisibleItems())
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
		} else if key == tcell.KeyPgUp {
			c.scrollByPage(-1)
			return nil
		} else if key == tcell.KeyPgDn {
			c.scrollByPage(1)
			return nil
		} else if key == tcell.KeyLeft {
			if c.scrollSelectedEntryHorizontal(-4) {
				return nil
			}
			return event
		} else if key == tcell.KeyRight {
			if c.scrollSelectedEntryHorizontal(4) {
				return nil
			}
			return event
		} else if key == tcell.KeyEnter {
			return nil
		}

		switch event.Rune() {
		case 'h':
			if c.scrollSelectedEntryHorizontal(-4) {
				return nil
			}
		case 'l':
			if c.scrollSelectedEntryHorizontal(4) {
				return nil
			}
		}

		return event
	})

	c.layout = layout
}

func (c *ListComponent[T]) SetDirection(direction int) {
	c.layout.SetDirection(direction)
}

func (c *ListComponent[T]) updateLayout() {
	c.updateLayoutInternal()
	c.application.ForceDraw()
}

func (c *ListComponent[T]) updateLayoutInternal() {
	c.updateVisibleEntries()
	c.updateScrollBarInternal()
}

func (c *ListComponent[T]) GetLayout() *tview.Flex {
	return c.layout
}

func (c *ListComponent[T]) SetTitle(title string) {
	SetupWindow(c.layout, title)
}

func (c *ListComponent[T]) GetData() []*T {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()
	return c.sortListEntries(c.entries, c.sortInverted)
}

func (c *ListComponent[T]) SetData(entries []*T) {
	selectedEntryBefore := c.GetSelectedItem()

	c.entriesMutex.Lock()
	c.entries = entries
	c.entriesMutex.Unlock()

	if len(entries) == 0 {
		c.selectedIndex = -1
	}

	c.updateLayout()

	if len(entries) == 0 {
		return
	}

	if selectedEntryBefore != nil && slices.Contains(entries, selectedEntryBefore) {
		c.selectedIndex = slices.Index(entries, selectedEntryBefore)
		c.scrollTo(selectedEntryBefore)
		c.application.ForceDraw()
		return
	}

	if c.selectedIndex < 0 || c.selectedIndex >= len(entries) {
		c.SelectFirst()
		return
	}

	c.application.ForceDraw()
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
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()
	return c.entries
}

func (c *ListComponent[T]) IsEmpty() bool {
	c.entriesMutex.Lock()
	defer c.entriesMutex.Unlock()
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

func (c *ListComponent[T]) scrollByPage(direction int) {
	if len(c.entries) == 0 {
		return
	}

	data := c.GetData()
	pageSize := c.GetMaxVisibleItems()
	if pageSize < 1 {
		pageSize = 1
	}

	selected := c.GetSelectedItem()
	currentIndex := slices.Index(data, selected)
	if currentIndex < 0 {
		currentIndex = 0
	}

	visibleMin, visibleMax := c.GetVisibleRange()
	if visibleMin < 0 || visibleMin >= len(data) || visibleMax < visibleMin {
		visibleMin = 0
		visibleMax = int(math.Min(float64(len(data)-1), float64(pageSize-1)))
	}

	if direction > 0 && visibleMax < len(data)-1 {
		c.setVisibleWindow(visibleMax+1, pageSize)
	} else if direction < 0 && visibleMin > 0 {
		c.setVisibleWindow(int(math.Max(0, float64(visibleMin-pageSize))), pageSize)
	}

	targetIndex := currentIndex + (direction * pageSize)
	targetIndex = int(math.Max(0, math.Min(float64(targetIndex), float64(len(data)-1))))
	c.selectAtDataIndex(data, targetIndex)
	c.application.ForceDraw()
}

func (c *ListComponent[T]) setVisibleWindow(startIndex int, pageSize int) {
	c.startIndex = startIndex
	c.updateLayout()
}

func (c *ListComponent[T]) selectAtDataIndex(data []*T, targetIndex int) {
	if len(data) == 0 {
		return
	}
	if targetIndex < 0 {
		targetIndex = 0
	}
	if targetIndex >= len(data) {
		targetIndex = len(data) - 1
	}

	target := data[targetIndex]
	targetLayout := c.getLayout(target)
	c.selectedIndex = slices.Index(c.entries, target)
	c.application.SetFocus(targetLayout)
	c.selectionChangedCallback(target)
}

func (c *ListComponent[T]) scrollSelectedEntryHorizontal(delta int) bool {
	selected := c.GetSelectedItem()
	if selected == nil {
		return false
	}

	scrollable, ok := any(selected).(HorizontalScrollable)
	if !ok {
		return false
	}
	scrollable.ScrollHorizontal(delta)
	return true
}

func (c *ListComponent[T]) scroll(rows int) {
	c.startIndex += rows
	c.updateLayout()
}

func (c *ListComponent[T]) GetVisibleRange() (int, int) {
	data := c.GetData()
	if len(data) == 0 {
		return 0, 0
	}
	maxVisible := c.GetMaxVisibleItems()
	endIndex := c.startIndex + maxVisible - 1
	if endIndex >= len(data) {
		endIndex = len(data) - 1
	}
	return c.startIndex, endIndex
}

func (c *ListComponent[T]) updateVisibleEntries() {
	// ensure we are displaying as many items as specified by MaxVisibleItems
	maxVisibleItems := c.GetMaxVisibleItems()
	data := c.GetData()

	if len(data) == 0 {
		c.startIndex = 0
	} else {
		if c.startIndex+maxVisibleItems > len(data) {
			c.startIndex = len(data) - maxVisibleItems
		}
		if c.startIndex < 0 {
			c.startIndex = 0
		}
	}

	// rebuild the visibility map
	c.entryVisibilityMap = map[*T]bool{}
	for index, entry := range data {
		c.entryVisibilityMap[entry] = index >= c.startIndex && index < c.startIndex+maxVisibleItems
	}

	// cleanup the entries layout
	c.entriesLayout.Clear()
	// create a layout for each visible entry
	for _, entry := range data {
		currentVisibility := c.entryVisibilityMap[entry]
		if currentVisibility {
			c.entriesLayout.AddItem(c.getLayout(entry), 0, 1, false)
		}
	}
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
	data := c.GetData()
	for idx, entry := range data {
		entryLayout := c.getLayout(entry)
		if entryLayout.HasFocus() {
			nextEntryIndex := (len(data) + idx + rows) % len(data)
			nextEntry := data[nextEntryIndex]
			nextEntryLayout := c.getLayout(nextEntry)
			c.selectedIndex = slices.Index(c.entries, nextEntry)
			c.application.SetFocus(nextEntryLayout)
			c.selectionChangedCallback(nextEntry)
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
	data := c.GetData()
	index := slices.Index(data, selection)
	if index < 0 {
		return false
	}
	maxVisible := c.GetMaxVisibleItems()
	return index >= c.startIndex && index < c.startIndex+maxVisible
}

func (c *ListComponent[T]) determineScrollDistanceToEntry(selection *T) int {
	data := c.GetData()
	index := slices.Index(data, selection)
	if index < 0 {
		return 0
	}

	maxVisible := c.GetMaxVisibleItems()
	if index < c.startIndex {
		return index - c.startIndex
	} else if index >= c.startIndex+maxVisible {
		return index - (c.startIndex + maxVisible - 1)
	} else {
		return 0
	}
}

func (c *ListComponent[T]) SelectFirst() {
	data := c.GetData()
	if len(data) > 0 {
		c.SelectEntry(data[0])
	}
}

func (c *ListComponent[T]) updateScrollBar() {
	c.updateScrollBarInternal()
	c.scrollbarComponent.UpdateLayout()
}

func (c *ListComponent[T]) updateScrollBarInternal() {
	if len(c.entries) <= c.GetMaxVisibleItems() {
		c.hideScrollbar()
	} else {
		c.showScrollbar()
	}

	minScrollIndex := 0
	c.scrollbarComponent.SetMinInternal(minScrollIndex)
	maxScrollIndex := int(math.Max(0.0, float64(len(c.entries))))
	c.scrollbarComponent.SetMaxInternal(maxScrollIndex)
	visibleIndexMin, visibleIndexMax := c.GetVisibleRange()

	c.scrollbarComponent.SetPositionInternal(visibleIndexMin)

	width := (visibleIndexMax - visibleIndexMin) + 1
	c.scrollbarComponent.SetWidthInternal(width)
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

func (c *ListComponent[T]) GetMaxVisibleItems() int {
	configValue := c.config.MaxVisibleItems
	if configValue > 0 {
		return configValue
	}
	if c.layout == nil || c.entriesLayout == nil {
		return 1
	}

	_, _, width, height := c.layout.GetInnerRect()
	// During startup, primitives may report a placeholder 1x1 rect before real layout is applied.
	if height <= 1 || width <= 1 || (len(c.entryVisibilityMap) == 0 && height < c.config.MinHeightPerEntry) {
		if len(c.entries) > 0 {
			return len(c.entries)
		}
		return 1
	}

	dynamicMaxVisibleItems := util.Coerce(float64(height/c.config.MinHeightPerEntry), 1, float64(1000))
	return int(dynamicMaxVisibleItems)
}
