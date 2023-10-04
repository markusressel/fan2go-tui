package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/data"
	"fan2go-tui/internal/ui/table"
	uiutil "fan2go-tui/internal/ui/util"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
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

type FanOverviewComponent struct {
	application *tview.Application

	Fans []*client.Fan

	layout                       *tview.Flex
	tableContainer               *table.RowSelectionTable[data.FanTableEntry]
	selectedEntryChangedCallback func(fileEntry *data.FanTableEntry)
	bmScatterPlot                *tvxwidgets.Plot
	graphComponents              map[string]*FanGraphComponent[client.Fan]
}

func NewFanOverviewComponent(application *tview.Application) *FanOverviewComponent {
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
				result = b.Pwm - a.Pwm
			case columnRpm:
				result = b.Rpm - a.Rpm
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

	c := &FanOverviewComponent{
		application:                  application,
		Fans:                         []*client.Fan{},
		tableContainer:               tableContainer,
		selectedEntryChangedCallback: func(fileEntry *data.FanTableEntry) {},
		graphComponents:              map[string]*FanGraphComponent[client.Fan]{},
	}

	c.layout = c.createLayout()

	tableContainer.SetColumnSpec(tableColumns, columnLabel, false)
	tableContainer.SetSelectionChangedCallback(func(selectedEntry *data.FanTableEntry) {
		c.selectedEntryChangedCallback(selectedEntry)
	})

	return c
}

func (c *FanOverviewComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	titleText := fmt.Sprintf("Data")

	layout.SetBorder(true)
	uiutil.SetupWindow(layout, titleText)

	return layout
}

func (c *FanOverviewComponent) Refresh() {
	for _, fan := range c.Fans {
		component, ok := c.graphComponents[fan.Label]
		if !ok {
			component = NewFanGraphComponent[client.Fan](c.application, fan, func(c *client.Fan) float64 {
				return float64(c.Rpm)
			})
			c.graphComponents[fan.Label] = component
			c.layout.AddItem(component.GetLayout(), 0, 1, false)
			component.InsertValue(fan)
			component.SetTitle(fan.Label)
			component.Refresh()
		} else {
			component.InsertValue(fan)
			component.Refresh()
		}
	}

	c.updateTableEntries()
}

func (c *FanOverviewComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanOverviewComponent) SetFans(fans []*client.Fan) {
	c.Fans = fans
	c.Refresh()
}

func (c *FanOverviewComponent) updateTableEntries() {
	var tableEntries []*data.FanTableEntry
	for _, fan := range c.Fans {
		entry := data.FanTableEntry{
			Label: fan.Label,
			Pwm:   fan.Pwm,
			Rpm:   fan.Rpm,
		}
		tableEntries = append(tableEntries, &entry)
	}
	c.tableContainer.SetData(tableEntries)
}
