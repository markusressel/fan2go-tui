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
	sensorGraphComponent  *SensorGraphComponent
	sensorGraphsComponent *SensorGraphsComponent
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

	sensors, err := c.client.GetSensors()
	if err != nil {
		//c.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return sensorsPageLayout
	}
	var sensorComponents []*SensorComponent
	for _, s := range *sensors {
		sensorComponent := NewSensorComponent(c.application, s)
		sensorComponents = append(sensorComponents, sensorComponent)
		sensorComponent.SetSensor(s)
		sensorComponent.Refresh()
		layout := sensorComponent.GetLayout()
		sensorInfoLayout.AddItem(layout, 0, 1, true)
	}
	c.sensorComponents = sensorComponents

	sensorGraphsComponent := NewSensorGraphsComponent(c.application)
	c.sensorGraphsComponent = sensorGraphsComponent
	// sensorComponents = append(sensorComponents, sensorGaphsComponent)

	// update overview
	sensorList := []*client.Sensor{}
	for _, f := range *sensors {
		sensorList = append(sensorList, f)
	}

	sensorGraphsComponent.SetSensors(sensorList)
	sensorGraphsComponent.Refresh()
	layout := sensorGraphsComponent.GetLayout()
	sensorsPageLayout.AddItem(layout, 0, 1, true)

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

	for _, s := range c.Sensors {
		for _, component := range c.sensorComponents {
			if component.Sensor.Config.ID == s.Config.ID {
				component.SetSensor(s)
				component.Refresh()
			}
		}
	}

	//c.sensorGraphsComponent.SetSensors(sensors)
	c.sensorGraphsComponent.Refresh()
}
