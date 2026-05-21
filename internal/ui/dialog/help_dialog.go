package dialog

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type HelpPage struct {
	layout *tview.Flex
}

func NewHelpPage() *HelpPage {
	helpPage := &HelpPage{}

	helpPage.createLayout()

	return helpPage
}

type TableEntry struct {
	Key   string
	Value string
}

var (
	emptyEntry = &TableEntry{Key: "", Value: ""}
)

func (p *HelpPage) createLayout() {
	helpTable := tview.NewTable()

	helpTableEntries := []*TableEntry{
		{Key: "F1, ?", Value: "Opens help dialog"},
		{Key: "1, 2, 3", Value: "Switches to the given page"},
		{Key: "Tab, Shift+Tab", Value: "Cycles forwards/backwards through all pages"},
		{Key: "PgUp, PgDn", Value: "Scroll list views up/down"},
		{Key: "Left/Right or h/l", Value: "Horizontal scroll in config info views"},
		{Key: "Mousescroll", Value: "Scroll within scrollable views"},
		{Key: "ctrl+r", Value: "Refreshes all data"},
		emptyEntry,
		{Key: "esc", Value: "Closes any currently open dialog"},
		{Key: "ctrl+q", Value: "Quits fan2go-tui"},
	}

	columns, rows := 2, len(helpTableEntries)
	for row := 0; row < rows; row++ {
		for column := 0; column < columns; column++ {
			entry := helpTableEntries[row]

			for col := 0; col < columns; col++ {
				var text string
				var cellAlignment int
				var cellColor = tcell.ColorWhite
				if col == 0 && entry != emptyEntry {
					text = fmt.Sprintf("%s:", entry.Key)
					cellAlignment = tview.AlignRight
					cellColor = tcell.ColorSteelBlue
				} else {
					text = entry.Value
					cellAlignment = tview.AlignLeft
				}
				helpTable.SetCell(
					row, col,
					tview.NewTableCell(text).SetAlign(cellAlignment).SetTextColor(cellColor),
				)
			}
		}
	}

	p.layout = createModal(" Help ", helpTable, 60, 14)
}

func (p *HelpPage) GetLayout() *tview.Flex {
	return p.layout
}
