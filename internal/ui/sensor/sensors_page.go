package sensor

import (
	"fan2go-tui/internal/client"
	uiutil "fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
	"sort"
	"strings"
)

type SensorsPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Sensors map[string]*client.Sensor

	layout *tview.Flex

	sensorInfoComponents  []*SensorInfoComponent
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

	sensorRowLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	sensorsPageLayout.AddItem(sensorRowLayout, 0, 1, true)

	sensors, sensorIds, err := c.fetchSensors()
	if err != nil {
		//c.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return sensorsPageLayout
	}
	for _, sId := range sensorIds {
		sensorColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
		uiutil.SetupWindow(sensorColumnLayout, sId)
		sensorColumnLayout.SetTitleAlign(tview.AlignLeft)
		sensorColumnLayout.SetBorder(true)
		sensorRowLayout.AddItem(sensorColumnLayout, 0, 1, true)

		s := (*sensors)[sId]
		sensorInfoComponent := NewSensorInfoComponent(c.application, s)
		c.sensorInfoComponents = append(c.sensorInfoComponents, sensorInfoComponent)
		sensorInfoComponent.SetSensor(s)
		sensorInfoComponent.Refresh()
		layout := sensorInfoComponent.GetLayout()
		sensorColumnLayout.AddItem(layout, 0, 1, true)

		sensorGraphComponent := NewSensorGraphComponent(c.application, s)
		c.sensorGraphComponents = append(c.sensorGraphComponents, sensorGraphComponent)
		sensorGraphComponent.SetSensor(s)
		sensorGraphComponent.Refresh()
		layout = sensorGraphComponent.GetLayout()
		sensorColumnLayout.AddItem(layout, 0, 3, false)
	}

	return sensorsPageLayout
}

func (c *SensorsPage) fetchSensors() (*map[string]*client.Sensor, []string, error) {
	result, err := c.client.GetSensors()
	if err != nil {
		return nil, nil, err
	}

	var sensorIds []string
	for _, s := range *result {
		sensorIds = append(sensorIds, s.Config.ID)
	}
	sort.SliceStable(sensorIds, func(i, j int) bool {
		a := sensorIds[i]
		b := sensorIds[j]

		result := strings.Compare(strings.ToLower(a), strings.ToLower(b))

		if result <= 0 {
			return true
		} else {
			return false
		}
	})

	return result, sensorIds, err
}

func (c *SensorsPage) GetLayout() *tview.Flex {
	return c.layout
}
func (c *SensorsPage) Refresh() {
	sensors, _, err := c.fetchSensors()
	if err != nil {
		return
	}
	c.Sensors = *sensors

	for _, component := range c.sensorInfoComponents {
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
