package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/data"
	"fan2go-tui/internal/ui/table"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sort"
	"strings"
)

var (
	columnLabel = &table.Column{
		Id:        0,
		Title:     "Label",
		Alignment: tview.AlignLeft,
	}
	columnPwm = &table.Column{
		Id:        1,
		Title:     "PWM",
		Alignment: tview.AlignLeft,
	}
	columnRpm = &table.Column{
		Id:        2,
		Title:     "RPM",
		Alignment: tview.AlignCenter,
	}

	tableColumns = []*table.Column{
		columnLabel,
		columnPwm,
		columnRpm,
	}
)

type FanComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout         *tview.Flex
	tableContainer *table.RowSelectionTable[data.FanTableEntry]

	selectedEntryChangedCallback func(fileEntry *data.FanTableEntry)

	pwmValueTextView *tview.TextView
	rpmValueTextView *tview.TextView
}

func NewFanComponent(application *tview.Application, fan *client.Fan) *FanComponent {
	toTableCellsFunction := func(row int, columns []*table.Column, entry *data.FanTableEntry) (cells []*tview.TableCell) {
		for _, column := range columns {
			var cellColor = tcell.ColorWhite
			var cellText string
			var cellAlignment = tview.AlignLeft
			var cellExpansion = 0

			if column == columnLabel {
				cellText = entry.Label
			} else if column == columnPwm {
				cellText = fmt.Sprintf("%d", entry.Pwm)
				cellAlignment = tview.AlignCenter
			} else if column == columnRpm {
				cellText = fmt.Sprintf("%d", entry.Rpm)
				cellAlignment = tview.AlignCenter
			} else {
				panic("Unknown column")
			}

			cell := tview.NewTableCell(cellText).
				SetTextColor(cellColor).
				SetAlign(cellAlignment).
				SetExpansion(cellExpansion)
			cells = append(cells, cell)
		}

		return cells
	}

	tableEntrySortFunction := func(entries []*data.FanTableEntry, columnToSortBy *table.Column, inverted bool) []*data.FanTableEntry {
		sort.SliceStable(entries, func(i, j int) bool {
			a := entries[i]
			b := entries[j]

			result := 0
			switch columnToSortBy {
			case columnLabel:
				result = strings.Compare(strings.ToLower(a.Label), strings.ToLower(b.Label))
			case columnPwm:
				result = int(b.Pwm - a.Pwm)
			case columnRpm:
				result = int(b.Rpm - a.Rpm)
			}

			if inverted {
				result *= -1
			}

			if result != 0 {
				if result <= 0 {
					return true
				} else {
					return false
				}
			}

			result = strings.Compare(strings.ToLower(a.Label), strings.ToLower(b.Label))
			if result != 0 {
				if result <= 0 {
					return true
				} else {
					return false
				}
			}

			if result <= 0 {
				return true
			} else {
				return false
			}
		})
		return entries
	}

	tableContainer := table.NewTableContainer[data.FanTableEntry](
		application,
		toTableCellsFunction,
		tableEntrySortFunction,
	)

	c := &FanComponent{
		application:                  application,
		Fan:                          fan,
		tableContainer:               tableContainer,
		selectedEntryChangedCallback: func(fileEntry *data.FanTableEntry) {},
	}

	c.layout = c.createLayout()

	tableContainer.SetColumnSpec(tableColumns, columnRpm, true)
	tableContainer.SetSelectionChangedCallback(func(selectedEntry *data.FanTableEntry) {
		c.selectedEntryChangedCallback(selectedEntry)
	})

	return c
}

func (c *FanComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	pwmValueTextView := tview.NewTextView()
	layout.AddItem(pwmValueTextView, 1, 0, true)
	c.pwmValueTextView = pwmValueTextView

	rpmValueTextView := tview.NewTextView()
	layout.AddItem(rpmValueTextView, 1, 0, true)
	c.rpmValueTextView = rpmValueTextView

	tableContainer := c.tableContainer.GetLayout()
	layout.AddItem(tableContainer, 0, 1, false)

	return layout
}

func (c *FanComponent) Refresh() {
	pwmText := fmt.Sprintf("PWM: %d", c.Fan.Pwm)
	c.pwmValueTextView.SetText(pwmText)

	rpmText := fmt.Sprintf("RPM: %d", c.Fan.Rpm)
	c.rpmValueTextView.SetText(rpmText)
}

func (c *FanComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
}
