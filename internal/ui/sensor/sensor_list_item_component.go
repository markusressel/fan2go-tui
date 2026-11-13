package sensor

import (
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/txwidgets"
	uiutil "fan2go-tui/internal/ui/util"

	"github.com/rivo/tview"
)

type SensorListItemComponent struct {
	application *tview.Application

	SensorState *state.SensorState

	layout *tview.Flex
	// window is the bordered box of the sensor, its title shows whether the sensor has warnings
	window *tview.Flex

	sensorInfoComponent  *SensorInfoComponent
	sensorGraphComponent *SensorGraphComponent
}

func NewSensorListItemComponent(application *tview.Application, sensorState *state.SensorState) *SensorListItemComponent {
	c := &SensorListItemComponent{
		application: application,
		SensorState: sensorState,
	}

	c.layout = c.createLayout()

	return c
}

func (c *SensorListItemComponent) createLayout() *tview.Flex {
	rootLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	sensorColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	uiutil.SetupWindow(sensorColumnLayout, c.SensorState.Sensor.Config.ID)
	c.window = sensorColumnLayout
	c.refreshWindowTitle()
	sensorColumnLayout.SetTitleAlign(tview.AlignLeft)
	sensorColumnLayout.SetBorder(true)
	rootLayout.AddItem(sensorColumnLayout, 0, 1, true)

	f := c.SensorState.Sensor
	sensorInfoComponent := NewSensorInfoComponent(c.application, f)
	c.sensorInfoComponent = sensorInfoComponent
	sensorInfoComponent.SetSensor(f)
	layout := sensorInfoComponent.GetLayout()
	sensorColumnLayout.AddItem(layout, 0, 1, true)
	sensorColumnLayout.AddItem(tview.NewBox(), 1, 0, false)

	sensorGraphComponent := NewSensorGraphComponent(c.application, c.SensorState)
	c.sensorGraphComponent = sensorGraphComponent
	sensorGraphComponent.SetSensor(c.SensorState)
	layout = sensorGraphComponent.GetLayout()
	sensorColumnLayout.AddItem(layout, 0, 3, true)

	return rootLayout
}

func (c *SensorListItemComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorListItemComponent) SetSensor(sensorState *state.SensorState) {
	c.SensorState = sensorState
	c.refresh()
}

func (c *SensorListItemComponent) refresh() {
	c.refreshWindowTitle()
	c.sensorInfoComponent.SetSensor(c.SensorState.Sensor)
	c.sensorGraphComponent.SetSensor(c.SensorState)
}

func (c *SensorListItemComponent) ScrollHorizontal(delta int) {
	if c.sensorInfoComponent != nil {
		c.sensorInfoComponent.ScrollHorizontal(delta)
	}
}

// refreshWindowTitle marks the title of the window if the sensor has warnings.
func (c *SensorListItemComponent) refreshWindowTitle() {
	if c.window == nil || c.SensorState == nil || c.SensorState.Sensor == nil {
		return
	}
	txwidgets.SetWindowWarnings(c.window.Box, c.SensorState.Sensor.Config.ID, c.SensorState.Sensor.Warnings)
}
