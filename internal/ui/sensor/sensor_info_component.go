package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/txwidgets"

	"github.com/rivo/tview"
)

type SensorInfoComponent struct {
	application *tview.Application

	Sensor *client.Sensor

	layout *tview.Flex

	configComponent *txwidgets.ConfigInfoComponent
}

func NewSensorInfoComponent(application *tview.Application, sensor *client.Sensor) *SensorInfoComponent {
	c := &SensorInfoComponent{
		application: application,
		Sensor:      sensor,
	}

	c.layout = c.createLayout()

	return c
}

func (c *SensorInfoComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	configComponent := txwidgets.NewConfigInfoComponent()
	layout.AddItem(configComponent.GetPrimitive(), 0, 1, false)
	c.configComponent = configComponent

	return layout
}

func (c *SensorInfoComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorInfoComponent) SetSensor(sensor *client.Sensor) {
	c.Sensor = sensor
	c.refresh()
}

func (c *SensorInfoComponent) refresh() {
	if c.Sensor == nil {
		c.configComponent.SetSections(nil)
		return
	}

	config := c.Sensor.Config
	c.configComponent.SetSections(txwidgets.SensorConfigSections(config))
}

func (c *SensorInfoComponent) ScrollHorizontal(delta int) {
	if c.configComponent != nil {
		c.configComponent.ScrollHorizontal(delta)
	}
}
