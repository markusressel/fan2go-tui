package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/util"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type SensorGraphsComponent struct {
	application *tview.Application

	Sensors []*client.Sensor

	layout          *tview.Flex
	bmScatterPlot   *tvxwidgets.Plot
	graphComponents map[string]*util.GraphComponent[client.Sensor]
}

func NewSensorGraphsComponent(application *tview.Application) *SensorGraphsComponent {
	c := &SensorGraphsComponent{
		application:     application,
		Sensors:         []*client.Sensor{},
		graphComponents: map[string]*util.GraphComponent[client.Sensor]{},
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
			component = util.NewGraphComponent[client.Sensor](
				c.application,
				sensor,
				func(c *client.Sensor) float64 {
					return c.MovingAvg / 1000
				},
				nil,
			)
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
