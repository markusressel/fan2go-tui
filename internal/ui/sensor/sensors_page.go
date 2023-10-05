package sensor

import (
	"fan2go-tui/internal/client"
	"github.com/rivo/tview"
)

type SensorsPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Sensors map[string]*client.Sensor

	layout *tview.Flex

	sensorComponents      []*SensorComponent
	sensorGraphComponents []*SensorGraphComponent
}

func NewSensorsPage(application *tview.Application, client client.Fan2goApiClient) SensorsPage {

	sensorsPage := SensorsPage{
		application: application,
		client:      client,
	}

	sensorsPage.layout = sensorsPage.createLayout()

	return sensorsPage
}

func (c *SensorsPage) createLayout() *tview.Flex {
	sensorsPageLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	sensorInfoLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	sensorsPageLayout.AddItem(sensorInfoLayout, 0, 1, true)
	sensorGraphsLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	sensorsPageLayout.AddItem(sensorGraphsLayout, 0, 1, false)

	sensors, err := c.client.GetSensors()
	if err != nil {
		//c.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return sensorsPageLayout
	}
	for _, s := range *sensors {
		sensorComponent := NewSensorComponent(c.application, s)
		c.sensorComponents = append(c.sensorComponents, sensorComponent)
		sensorComponent.SetSensor(s)
		sensorComponent.Refresh()
		layout := sensorComponent.GetLayout()
		sensorInfoLayout.AddItem(layout, 0, 1, true)

		sensorGraphComponent := NewSensorGraphComponent(c.application, s)
		c.sensorGraphComponents = append(c.sensorGraphComponents, sensorGraphComponent)
		sensorGraphComponent.SetTitle(s.Config.ID)
		sensorGraphComponent.SetSensor(s)
		sensorGraphComponent.Refresh()
		layout = sensorGraphComponent.GetLayout()
		sensorGraphsLayout.AddItem(layout, 0, 1, false)
	}

	return sensorsPageLayout
}

func (c *SensorsPage) fetchSensors() (*map[string]*client.Sensor, error) {
	return c.client.GetSensors()
}

func (c *SensorsPage) GetLayout() *tview.Flex {
	return c.layout
}
func (c *SensorsPage) Refresh() {
	sensors, err := c.fetchSensors()
	if err != nil {
		return
	}
	c.Sensors = *sensors

	for _, component := range c.sensorComponents {
		sensor, ok := (*sensors)[component.Sensor.Config.ID]
		if !ok {
			continue
		}
		component.SetSensor(sensor)
		component.Refresh()
	}

	for _, component := range c.sensorGraphComponents {
		if component.Sensor == nil {
			continue
		}
		sensor, ok := (*sensors)[component.Sensor.Config.ID]
		if !ok || sensor == nil {
			continue
		}
		component.SetSensor(sensor)
		component.Refresh()
	}
}
