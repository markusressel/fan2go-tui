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
	configTextView   *tview.TextView
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
	layout.AddItem(pwmValueTextView, 1, 0, false)
	c.pwmValueTextView = pwmValueTextView

	rpmValueTextView := tview.NewTextView()
	layout.AddItem(rpmValueTextView, 1, 0, false)
	c.rpmValueTextView = rpmValueTextView

	configTextView := tview.NewTextView()
	layout.AddItem(configTextView, 0, 1, false)
	c.configTextView = configTextView

	tableContainer := c.tableContainer.GetLayout()
	layout.AddItem(tableContainer, 3, 0, true)

	return layout
}

func (c *FanComponent) Refresh() {
	// print basic info
	pwmText := fmt.Sprintf("PWM: %d", c.Fan.Pwm)
	c.pwmValueTextView.SetText(pwmText)

	rpmText := fmt.Sprintf("RPM: %d", c.Fan.Rpm)
	c.rpmValueTextView.SetText(rpmText)

	// print config
	config := c.Fan.Config

	configText := ""
	configText += fmt.Sprintf("Id: %s\n", config.Id)
	configText += fmt.Sprintf("Curve: %s\n", config.Curve)
	configText += fmt.Sprintf("Pwm:\n")
	configText += fmt.Sprintf("  Min: %d\n", *config.MinPwm)
	configText += fmt.Sprintf("  Start: %d\n", *config.StartPwm)
	configText += fmt.Sprintf("  Max: %d\n", *config.MaxPwm)
	configText += fmt.Sprintf("NeverStop: %v\n", config.NeverStop)

	// value = strconv.FormatFloat(config.MinPwm, 'f', -1, 64)

	if config.File != nil {
		configText += fmt.Sprintf("Type: File\n")
		configText += fmt.Sprintf("  PwmPath: %s\n", config.File.Path)
		configText += fmt.Sprintf("  RpmPath: %s\n", config.File.RpmPath)
	} else if config.HwMon != nil {
		configText += fmt.Sprintf("Type: HwMon\n")
		configText += fmt.Sprintf("  Platform: %s\n", config.HwMon.Platform)
		configText += fmt.Sprintf("  Index: %d\n", config.HwMon.Index)
		configText += fmt.Sprintf("  PwmChannel: %d\n", config.HwMon.PwmChannel)
		configText += fmt.Sprintf("  RpmChannel: %d\n", config.HwMon.RpmChannel)
		configText += fmt.Sprintf("  SysfsPath: %s\n", config.HwMon.SysfsPath)
		configText += fmt.Sprintf("  PwmPath: %s\n", config.HwMon.PwmPath)
		configText += fmt.Sprintf("  PwmEnablePath: %s\n", config.HwMon.PwmEnablePath)
		configText += fmt.Sprintf("  RpmInputPath: %s\n", config.HwMon.RpmInputPath)
	} else if config.Cmd != nil {
		configText += fmt.Sprintf("Type: Cmd\n")
		configText += fmt.Sprintf("  GetPwm: %s\n", config.Cmd.GetPwm)
		configText += fmt.Sprintf("  SetPwm: %s\n", config.Cmd.SetPwm)
		configText += fmt.Sprintf("  GetRpm: %s\n", config.Cmd.GetRpm)
	}
	c.configTextView.SetText(configText)
}

func (c *FanComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
}
