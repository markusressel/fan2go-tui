package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
	uiutil "fan2go-tui/internal/ui/util"
	"fmt"
	"github.com/rivo/tview"
)

type SensorComponent struct {
	application *tview.Application

	Sensor *client.Sensor

	layout *tview.Flex

	configTextView *tview.TextView
	valueTextView  *tview.TextView
}

func NewSensorComponent(application *tview.Application, sensor *client.Sensor) *SensorComponent {
	c := &SensorComponent{
		application: application,
		Sensor:      sensor,
	}

	c.layout = c.createLayout()

	return c
}

func (c *SensorComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	titleText := fmt.Sprintf("Curve: %s", c.Sensor.Config.ID)

	layout.SetBorder(true)
	uiutil.SetupWindow(layout, titleText)

	configTextView := tview.NewTextView()
	layout.AddItem(configTextView, 0, 1, false)
	c.configTextView = configTextView

	curveValueTextView := tview.NewTextView().SetTextColor(theme.Colors.Graphs.Curve)
	layout.AddItem(curveValueTextView, 1, 0, false)
	c.valueTextView = curveValueTextView

	return layout
}

func (c *SensorComponent) Refresh() {
	// print basic info
	valueText := fmt.Sprintf("Avg: %f", c.Sensor.MovingAvg/1000)
	c.valueTextView.SetText(valueText)

	// print config
	config := c.Sensor.Config

	configText := ""
	// configText += fmt.Sprintf("Id: %s\n", config.Id)
	configText += fmt.Sprintf("Sensor: %s\n", config.ID)
	// value = strconv.FormatFloat(config.MinPwm, 'f', -1, 64)

	if config.HwMon != nil {
		configText += fmt.Sprintf("Type: HwMon\n")
		configText += fmt.Sprintf("  Index: %d\n", config.HwMon.Index)
		configText += fmt.Sprintf("  Platform: %s\n", config.HwMon.Platform)
		configText += fmt.Sprintf("  TempInput: %s\n", config.HwMon.TempInput)
	} else if config.File != nil {
		configText += fmt.Sprintf("Type: File\n")
		configText += fmt.Sprintf("  Path: %s\n", config.File.Path)
	} else if config.Cmd != nil {
		configText += fmt.Sprintf("Type: Cmd\n")
		configText += fmt.Sprintf("  Exec: %s\n", config.Cmd.Exec)
		configText += fmt.Sprintf("  Args: %s\n", config.Cmd.Args)
	}
	c.configTextView.SetText(configText)
}

func (c *SensorComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorComponent) SetSensor(sensor *client.Sensor) {
	c.Sensor = sensor
}
