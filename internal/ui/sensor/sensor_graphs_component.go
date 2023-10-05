package sensor

import (
	"fan2go-tui/internal/client"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type SensorGraphsComponent struct {
	application *tview.Application

	Sensors []*client.Sensor

	layout          *tview.Flex
	bmScatterPlot   *tvxwidgets.Plot
	graphComponents map[string]*SensorGraphComponent
}

func NewSensorGraphsComponent(application *tview.Application) *SensorGraphsComponent {
	c := &SensorGraphsComponent{
		application:     application,
		Sensors:         []*client.Sensor{},
		graphComponents: map[string]*SensorGraphComponent{},
	}

	c.layout = c.createLayout()

	return c
}

func (c *SensorGraphsComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.SetBorder(false)

	return layout
}

func (c *SensorGraphsComponent) Refresh() {
	for _, sensor := range c.Sensors {
		component, ok := c.graphComponents[sensor.Config.ID]
		if !ok {
			component = NewSensorGraphComponent(c.application, sensor)
			c.graphComponents[sensor.Config.ID] = component
			c.layout.AddItem(component.GetLayout(), 0, 1, false)
			component.InsertValue(sensor)
			component.SetTitle(sensor.Config.ID)
			component.Refresh()
		} else {
			component.InsertValue(sensor)
			component.Refresh()
		}
	}
}

func (c *SensorGraphsComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorGraphsComponent) SetSensors(sensors []*client.Sensor) {
	c.Sensors = sensors
	c.Refresh()
}
