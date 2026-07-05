package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/shortcut_helper"
	"fan2go-tui/internal/ui/util"
	"sort"
	"strings"

	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type SensorsPage struct {
	application *tview.Application

	store *state.Store

	Sensors map[string]*client.Sensor

	layout          *tview.Flex
	sensorRowLayout *tview.Flex

	sensorList *util.ListComponent[SensorListItemComponent]

	sensorListItemComponents map[string]*SensorListItemComponent
}

func NewSensorsPage(application *tview.Application, store *state.Store) SensorsPage {

	sensorsPage := SensorsPage{
		application:              application,
		store:                    store,
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
				a := entries[i].SensorState.Sensor.Config.ID
				b := entries[j].SensorState.Sensor.Config.ID

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

func (c *SensorsPage) getSensorIds(sensors map[string]*client.Sensor) []string {
	var sensorIds []string
	for _, s := range sensors {
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

	return sensorIds
}

func (c *SensorsPage) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorsPage) Refresh() error {
	sensors := c.store.GetSensors()
	sensorIds := c.getSensorIds(sensors)

	var sensorListItemsComponents []*SensorListItemComponent

	oldSIds := maps.Keys(c.sensorListItemComponents)
	// remove now nonexisting entries
	for _, oldSId := range oldSIds {
		_, ok := sensors[oldSId]
		if !ok {
			delete(c.sensorListItemComponents, oldSId)
		}
	}

	// add new entries / update existing entries
	for _, sId := range sensorIds {
		sensorState := c.store.GetSensorState(sId)
		sensorListItemComponent, ok := c.sensorListItemComponents[sId]
		if ok {
			sensorListItemComponent.SetSensor(sensorState)
			sensorListItemsComponents = append(sensorListItemsComponents, sensorListItemComponent)
		} else {
			sensorListItemComponent = NewSensorListItemComponent(c.application, sensorState)
			c.sensorListItemComponents[sId] = sensorListItemComponent
			sensorListItemComponent.SetSensor(sensorState)
			sensorListItemsComponents = append(sensorListItemsComponents, sensorListItemComponent)
		}
	}

	c.sensorList.SetData(sensorListItemsComponents)

	return nil
}

func (c *SensorsPage) ScrollToItem() {
	c.sensorList.SelectEntry(c.sensorList.GetSelectedItem())
}

func (c *SensorsPage) SelectSensorByID(sensorID string) bool {
	sensorListItem, ok := c.sensorListItemComponents[sensorID]
	if !ok {
		return false
	}
	c.sensorList.SelectEntry(sensorListItem)
	return true
}

func (c *SensorsPage) GetShortcutMap() []shortcut_helper.ShortcutEntry {
	return []shortcut_helper.ShortcutEntry{
		{KeyCombo: []string{"↑", "↓"}, Name: "Select"},
		{KeyCombo: []string{"PgUp", "PgDn"}, Name: "Scroll"},
	}
}
