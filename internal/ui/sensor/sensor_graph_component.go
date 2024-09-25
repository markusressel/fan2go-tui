package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/util"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type SensorGraphComponent struct {
	application *tview.Application

	Sensor *client.Sensor

	layout         *tview.Flex
	bmScatterPlot  *tvxwidgets.Plot
	graphComponent *util.GraphComponent[client.Sensor]
}

func NewSensorGraphComponent(application *tview.Application, sensor *client.Sensor) *SensorGraphComponent {

	graphComponent := util.NewGraphComponent[client.Sensor](
		application,
		util.NewGraphComponentConfig().WithReversedOrder(),
		sensor,
		[]func(*client.Sensor) float64{
			func(c *client.Sensor) float64 {
				return c.MovingAvg / 1000
			},
		},
	)

	c := &SensorGraphComponent{
		application:    application,
		graphComponent: graphComponent,
		Sensor:         sensor,
	}

	c.layout = c.createLayout()
	c.layout.AddItem(graphComponent.GetLayout(), 0, 1, false)

	return c
}

func (c *SensorGraphComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	return layout
}

func (c *SensorGraphComponent) refresh() {
	sensor := c.Sensor
	if sensor == nil {
		return
	}
	component := c.graphComponent
	component.InsertValue(sensor)
}

func (c *SensorGraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorGraphComponent) SetSensor(sensor *client.Sensor) {
	c.Sensor = sensor
	c.refresh()
}

func (c *SensorGraphComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
