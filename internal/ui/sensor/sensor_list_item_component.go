package sensor

import (
	"fan2go-tui/internal/client"
	uiutil "fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
)

type SensorListItemComponent struct {
	application *tview.Application

	Sensor *client.Sensor

	layout *tview.Flex

	sensorInfoComponent  *SensorInfoComponent
	sensorGraphComponent *SensorGraphComponent
}

func NewSensorListItemComponent(application *tview.Application, sensor *client.Sensor) *SensorListItemComponent {
	c := &SensorListItemComponent{
		application: application,
		Sensor:      sensor,
	}

	c.layout = c.createLayout()

	return c
}

func (c *SensorListItemComponent) createLayout() *tview.Flex {
	rootLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	sensorColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	uiutil.SetupWindow(sensorColumnLayout, c.Sensor.Config.ID)
	sensorColumnLayout.SetTitleAlign(tview.AlignLeft)
	sensorColumnLayout.SetBorder(true)
	rootLayout.AddItem(sensorColumnLayout, 0, 1, true)

	f := c.Sensor
	sensorInfoComponent := NewSensorInfoComponent(c.application, f)
	c.sensorInfoComponent = sensorInfoComponent
	sensorInfoComponent.Refresh()
	layout := sensorInfoComponent.GetLayout()
	sensorColumnLayout.AddItem(layout, 0, 1, true)

	sensorGraphComponent := NewSensorGraphComponent(c.application, f)
	c.sensorGraphComponent = sensorGraphComponent
	sensorGraphComponent.Refresh()
	layout = sensorGraphComponent.GetLayout()
	sensorColumnLayout.AddItem(layout, 0, 3, true)

	return rootLayout
}

func (c *SensorListItemComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorListItemComponent) SetSensor(sensor *client.Sensor) {
	c.Sensor = sensor
	c.Refresh()
}

func (c *SensorListItemComponent) Refresh() {
	c.sensorInfoComponent.SetSensor(c.Sensor)
	c.sensorGraphComponent.InsertValue(c.Sensor)
}
