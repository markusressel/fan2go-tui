package curve

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

type CurvesPage struct {
	application *tview.Application

	store *state.Store

	Curves             map[string]*client.Curve
	entryVisibilityMap []string

	layout         *tview.Flex
	curveRowLayout *tview.Flex

	curveList *util.ListComponent[CurveListItemComponent]

	curveListItemComponents map[string]*CurveListItemComponent
	onOpenSensor            func(sensorID string)
	onOpenCurve             func(curveID string)
}

func NewCurvesPage(application *tview.Application, store *state.Store, onOpenSensor func(sensorID string), onOpenCurve func(curveID string)) CurvesPage {

	curvesPage := CurvesPage{
		application:             application,
		store:                   store,
		curveListItemComponents: map[string]*CurveListItemComponent{},
		onOpenSensor:            onOpenSensor,
		onOpenCurve:             onOpenCurve,
	}

	curvesPage.layout = curvesPage.createLayout()

	return curvesPage
}

func (c *CurvesPage) createLayout() *tview.Flex {
	curvesPageLayout := tview.NewFlex()

	listConfig := util.NewListComponentConfig()
	curveListComponent := util.NewListComponent[CurveListItemComponent](
		c.application,
		listConfig,
		func(entry *CurveListItemComponent) (layout *tview.Flex) {
			return entry.GetLayout()
		},
		//func(a, b *CurveListItemComponent) bool {
		//	return strings.Compare(a.Curve.Config.ID, b.Curve.Config.ID) <= 0
		//},
		func(entries []*CurveListItemComponent, inverted bool) []*CurveListItemComponent {
			sort.SliceStable(entries, func(i, j int) bool {
				a := entries[i].CurveState.Curve.Config.ID
				b := entries[j].CurveState.Curve.Config.ID

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
	c.curveList = curveListComponent
	curvesPageLayout.AddItem(c.curveList.GetLayout(), 0, 1, true)

	return curvesPageLayout
}

func (c *CurvesPage) getCurveIds(curves map[string]*client.Curve) []string {
	var curveIds []string
	for _, c := range curves {
		curveIds = append(curveIds, c.Config.ID)
	}
	sort.SliceStable(curveIds, func(i, j int) bool {
		a := curveIds[i]
		b := curveIds[j]

		result := strings.Compare(strings.ToLower(a), strings.ToLower(b))

		if result <= 0 {
			return true
		} else {
			return false
		}
	})

	return curveIds
}

func (c *CurvesPage) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurvesPage) Refresh() error {
	curves := c.store.GetCurves()
	curveIds := c.getCurveIds(curves)

	var curveListItemsComponents []*CurveListItemComponent

	oldCIds := maps.Keys(c.curveListItemComponents)
	// remove now non-existing entries
	for _, oldCId := range oldCIds {
		_, ok := curves[oldCId]
		if !ok {
			delete(c.curveListItemComponents, oldCId)
		}
	}

	// add new entries / update existing entries
	for _, cId := range curveIds {
		curveState := c.store.GetCurveState(cId)
		curveListItemComponent, ok := c.curveListItemComponents[cId]
		if ok {
			curveListItemComponent.SetCurve(curveState)
			curveListItemsComponents = append(curveListItemsComponents, curveListItemComponent)
		} else {
			curveListItemComponent = NewCurveListItemComponent(c.application, curveState, c.onOpenSensor, c.onOpenCurve)
			curveListItemComponent.SetCurve(curveState)
			c.curveListItemComponents[cId] = curveListItemComponent
			curveListItemsComponents = append(curveListItemsComponents, curveListItemComponent)
		}
	}

	c.curveList.SetData(curveListItemsComponents)

	return nil
}

func (c *CurvesPage) ScrollToItem() {
	c.curveList.SelectEntry(c.curveList.GetSelectedItem())
}

func (c *CurvesPage) SelectCurveByID(curveID string) bool {
	curveListItem, ok := c.curveListItemComponents[curveID]
	if !ok {
		return false
	}
	c.curveList.SelectEntry(curveListItem)
	return true
}

func (c *CurvesPage) GetShortcutMap() []shortcut_helper.ShortcutEntry {
	return []shortcut_helper.ShortcutEntry{
		{KeyCombo: []string{"↑", "↓"}, Name: "Select"},
		{KeyCombo: []string{"PgUp", "PgDn"}, Name: "Scroll"},
	}
}
