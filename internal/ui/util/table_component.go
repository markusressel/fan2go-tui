package util

import (
	"fmt"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
)

type ColumnId int

type Column struct {
	Id        ColumnId
	Title     string
	Alignment int
}

type RowSelectionTable[T any] struct {
	application *tview.Application

	layout *tview.Table

	entries      []*T
	entriesMutex sync.Mutex

	sortByColumn             *Column
	sortTableEntries         func(entries []*T, column *Column, inverted bool) []*T
	toTableCells             func(row int, columns []*Column, entry *T) (cells []*tview.TableCell)
	inputCapture             func(event *tcell.EventKey) *tcell.EventKey
	doubleClickCallback      func()
	selectionChangedCallback func(selectedEntry *T)

	columnSpec   []*Column
	sortInverted bool
}

func NewTableContainer[T any](
	application *tview.Application,
	toTableCells func(row int, columns []*Column, entry *T) (cells []*tview.TableCell),
	sortTableEntries func(entries []*T, column *Column, inverted bool) []*T,
) *RowSelectionTable[T] {
	tableContainer := &RowSelectionTable[T]{
		application:      application,
		entriesMutex:     sync.Mutex{},
		toTableCells:     toTableCells,
		sortTableEntries: sortTableEntries,
		inputCapture: func(event *tcell.EventKey) *tcell.EventKey {
			return event
		},
		doubleClickCallback:      func() {},
		selectionChangedCallback: func(selectedEntry *T) {},
	}
	tableContainer.createLayout()
	return tableContainer
}

func (c *RowSelectionTable[T]) createLayout() {
	table := tview.NewTable()

	table.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		switch action {
		case tview.MouseLeftDoubleClick:
			go func() {
				c.doubleClickCallback()
			}()
			return action, nil
		}
		return action, event
	})

	table.SetBorder(true)
	table.SetBorders(false)
	table.SetBorderPadding(0, 0, 1, 1)

	// fixed header row
	table.SetFixed(1, 0)

	SetupWindow(table, "")

	table.SetSelectable(true, false)
	table.SetSelectionChangedFunc(func(row, column int) {
		selectedEntry := c.GetSelectedEntry()
		c.selectionChangedCallback(selectedEntry)
	})

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		event = c.inputCapture(event)
		if event == nil {
			return event
		}
		key := event.Key()
		if c.GetSelectedEntry() == nil {
			if key == tcell.KeyRight {
				c.nextSortOrder()
				return nil
			} else if key == tcell.KeyLeft {
				c.previousSortOrder()
				return nil
			} else if key == tcell.KeyEnter {
				c.toggleSortDirection()
				return nil
			}
		}
		return event
	})

	c.layout = table
}

func (c *RowSelectionTable[T]) GetLayout() *tview.Table {
	return c.layout
}

func (c *RowSelectionTable[T]) SetTitle(title string) {
	SetupWindow(c.layout, title)
}

func (c *RowSelectionTable[T]) SetColumnSpec(columns []*Column, defaultSortColumn *Column, inverted bool) {
	c.columnSpec = columns
	c.SortBy(defaultSortColumn, inverted)
	c.updateTableContents()
}

func (c *RowSelectionTable[T]) SetData(entries []*T) {
	c.entriesMutex.Lock()
	c.entries = entries
	c.entriesMutex.Unlock()
	c.SortBy(c.sortByColumn, c.sortInverted)
	c.updateTableContents()
}

func (c *RowSelectionTable[T]) SetDoubleClickCallback(f func()) {
	c.doubleClickCallback = f
}

func (c *RowSelectionTable[T]) SortBy(sortOption *Column, inverted bool) {
	c.entriesMutex.Lock()
	c.sortByColumn = sortOption
	c.sortInverted = inverted
	c.entries = c.sortTableEntries(c.entries, c.sortByColumn, c.sortInverted)
	c.entriesMutex.Unlock()
}

func (c *RowSelectionTable[T]) nextSortOrder() {
	currentIndex := slices.Index(c.columnSpec, c.sortByColumn)
	nextIndex := (currentIndex + 1) % len(c.columnSpec)
	column := c.columnSpec[nextIndex]
	c.SortBy(column, c.sortInverted)
	c.updateTableContents()
}

func (c *RowSelectionTable[T]) previousSortOrder() {
	currentIndex := slices.Index(c.columnSpec, c.sortByColumn)
	nextIndex := (len(c.columnSpec) + currentIndex - 1) % len(c.columnSpec)
	column := c.columnSpec[nextIndex]
	c.SortBy(column, c.sortInverted)
	c.updateTableContents()
}

func (c *RowSelectionTable[T]) toggleSortDirection() {
	c.sortInverted = !c.sortInverted
	c.SortBy(c.sortByColumn, c.sortInverted)
	c.updateTableContents()
}

func (c *RowSelectionTable[T]) updateTableContents() {

	table := c.layout
	if table == nil {
		return
	}

	table.Clear()

	// Table Header
	for column, tableColumn := range c.columnSpec {
		cellColor := tcell.ColorWhite
		cellAlignment := tableColumn.Alignment
		cellExpansion := 0

		cellText := tableColumn.Title
		if tableColumn == c.sortByColumn {
			var sortDirectionIndicator = "↓"
			if !c.sortInverted {
				sortDirectionIndicator = "↑"
			}
			cellText = fmt.Sprintf("%s %s", cellText, sortDirectionIndicator)
		}

		cell := tview.NewTableCell(cellText).
			SetTextColor(cellColor).
			SetAlign(cellAlignment).
			SetExpansion(cellExpansion)
		table.SetCell(0, column, cell)
	}

	// Table Content
	for row, entry := range c.entries {
		cells := c.toTableCells(row, c.columnSpec, entry)
		for column, cell := range cells {
			table.SetCell(row+1, column, cell)
		}
	}
}

func (c *RowSelectionTable[T]) Select(entry *T) {
	index := 0
	if entry != nil {
		index = slices.Index(c.entries, entry)
		if index < 0 {
			return
		} else {
			index += 1
		}
	}
	if index <= 1 {
		c.layout.ScrollToBeginning()
	}
	c.layout.Select(index, 0)
	c.application.ForceDraw()
}

func (c *RowSelectionTable[T]) HasFocus() bool {
	return c.layout.HasFocus()
}

func (c *RowSelectionTable[T]) GetEntries() []*T {
	return c.entries
}

func (c *RowSelectionTable[T]) GetSelectedEntry() *T {
	row, _ := c.layout.GetSelection()
	row -= 1
	if row >= 0 && row < len(c.entries) {
		return c.entries[row]
	} else {
		return nil
	}
}

func (c *RowSelectionTable[T]) IsEmpty() bool {
	return len(c.entries) <= 0
}

func (c *RowSelectionTable[T]) SetInputCapture(inputCapture func(event *tcell.EventKey) *tcell.EventKey) {
	c.inputCapture = inputCapture
}

func (c *RowSelectionTable[T]) SetSelectionChangedCallback(f func(selectedEntry *T)) {
	c.selectionChangedCallback = f
}

func (c *RowSelectionTable[T]) SelectHeader() {
	row, col := c.layout.GetSelection()
	if row != 0 || col != 0 {
		c.layout.Select(0, 0)
	}
}

func (c *RowSelectionTable[T]) SelectFirstIfExists() {
	if len(c.entries) > 0 {
		c.Select(c.entries[0])
	}
}
