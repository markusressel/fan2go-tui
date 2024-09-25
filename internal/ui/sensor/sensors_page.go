package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"sort"
	"strings"
)

type SensorsPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Sensors map[string]*client.Sensor

	layout          *tview.Flex
	sensorRowLayout *tview.Flex

	sensorList *util.ListComponent[SensorListItemComponent]

	sensorListItemComponents map[string]*SensorListItemComponent
}

func NewSensorsPage(application *tview.Application, client client.Fan2goApiClient) SensorsPage {

	sensorsPage := SensorsPage{
		application:              application,
		client:                   client,
		sensorListItemComponents: map[string]*SensorListItemComponent{},
	}

	sensorsPage.layout = sensorsPage.createLayout()

	return sensorsPage
}

func (c *SensorsPage) createLayout() *tview.Flex {
	sensorsPageLayout := tview.NewFlex()

	listConfig := util.NewListComponentConfig()
	sensorListComponent := util.NewListComponent[SensorListItemComponent](
		c.application,
		listConfig,
		func(entry *SensorListItemComponent) (layout *tview.Flex) {
			return entry.GetLayout()
		},
		//func(a, b *SensorListItemComponent) bool {
		//	return strings.Compare(a.Sensor.Config.ID, b.Sensor.Config.ID) <= 0
		//},
		func(entries []*SensorListItemComponent, inverted bool) []*SensorListItemComponent {
			sort.SliceStable(entries, func(i, j int) bool {
				a := entries[i].Sensor.Config.ID
				b := entries[j].Sensor.Config.ID

				result := strings.Compare(strings.ToLower(a), strings.ToLower(b))

				if result <= 0 {
					return true
				} else {
					return false
				}
			})

			if inverted {
				slices.Reverse(entries)
			}

			return entries
		},
	)
	c.sensorList = sensorListComponent
	sensorsPageLayout.AddItem(c.sensorList.GetLayout(), 0, 1, true)

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

func (c *SensorsPage) Refresh() error {
	sensors, sensorIds, err := c.fetchSensors()
	if err != nil || sensors == nil {
		sensors = &map[string]*client.Sensor{}
	}

	var sensorListItemsComponents []*SensorListItemComponent

	oldSIds := maps.Keys(c.sensorListItemComponents)
	// remove now nonexisting entries
	for _, oldSId := range oldSIds {
		_, ok := (*sensors)[oldSId]
		if !ok {
			delete(c.sensorListItemComponents, oldSId)
		}
	}

	// add new entries / update existing entries
	for _, sId := range sensorIds {
		sensor := (*sensors)[sId]
		sensorListItemComponent, ok := c.sensorListItemComponents[sId]
		if ok {
			sensorListItemComponent.SetSensor(sensor)
			sensorListItemsComponents = append(sensorListItemsComponents, sensorListItemComponent)
		} else {
			sensorListItemComponent = NewSensorListItemComponent(c.application, sensor)
			c.sensorListItemComponents[sId] = sensorListItemComponent
			sensorListItemComponent.SetSensor(sensor)
			sensorListItemsComponents = append(sensorListItemsComponents, sensorListItemComponent)
		}
	}

	c.sensorList.SetData(sensorListItemsComponents)

	return err
}

func (c *SensorsPage) ScrollToItem() {
	c.sensorList.SelectEntry(c.sensorList.GetSelectedItem())
}
