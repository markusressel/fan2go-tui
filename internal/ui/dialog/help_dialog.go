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
		{Key: "up, k", Value: "Moves cursor up"},
		{Key: "down, j", Value: "Moves cursor down"},
		{Key: "left, h", Value: "Opens parent directory"},
		{Key: "right", Value: "Enters selected directory"},
		{Key: "enter", Value: "Opens file action dialog"},
		{Key: "tab, shift+tab", Value: "Cycles window focus forwards"},
		{Key: "shift+tab", Value: "Cycles window focus backwards"},
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

	p.layout = createModal(" Help ", helpTable, 60, 13)
}

func (p *HelpPage) GetLayout() *tview.Flex {
	return p.layout
}
