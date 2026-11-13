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

	warningsComponent *txwidgets.WarningsComponent
	configComponent   *txwidgets.ConfigInfoComponent
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

	// hidden until there are warnings, see refresh
	c.warningsComponent = txwidgets.NewWarningsComponent()
	layout.AddItem(c.warningsComponent.GetPrimitive(), 0, 0, false)

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
		c.setWarnings(nil)
		c.configComponent.SetSections(nil)
		return
	}

	c.setWarnings(c.Sensor.Warnings)

	config := c.Sensor.Config
	c.configComponent.SetSections(txwidgets.SensorConfigSections(config))
}

func (c *SensorInfoComponent) ScrollHorizontal(delta int) {
	if c.configComponent != nil {
		c.configComponent.ScrollHorizontal(delta)
	}
}

// setWarnings shows the warnings of the sensor above its configuration, and hides them if there are none.
func (c *SensorInfoComponent) setWarnings(warnings []client.Warning) {
	c.warningsComponent.SetWarnings(warnings)
	proportion := 0
	if c.warningsComponent.HasWarnings() {
		proportion = 1
	}
	c.layout.ResizeItem(c.warningsComponent.GetPrimitive(), 0, proportion)
}
