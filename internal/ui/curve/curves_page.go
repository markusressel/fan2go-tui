package curve

import (
	"fan2go-tui/internal/client"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"sort"
	"strings"
)

type CurvesPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Curves map[string]*client.Curve

	layout         *tview.Flex
	curveRowLayout *tview.Flex

	curveListItemComponents map[string]*CurveListItemComponent
}

func NewCurvesPage(application *tview.Application, client client.Fan2goApiClient) CurvesPage {

	curvesPage := CurvesPage{
		application:             application,
		client:                  client,
		curveListItemComponents: map[string]*CurveListItemComponent{},
	}

	curvesPage.layout = curvesPage.createLayout()

	return curvesPage
}

func (c *CurvesPage) createLayout() *tview.Flex {
	fansPageLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	c.curveRowLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	fansPageLayout.AddItem(c.curveRowLayout, 0, 1, true)

	return fansPageLayout
}

func (c *CurvesPage) fetchCurves() (*map[string]*client.Curve, []string, error) {
	result, err := c.client.GetCurves()
	if err != nil {
		return nil, nil, err
	}

	var curveIds []string
	for _, c := range *result {
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

	return result, curveIds, err
}

func (c *CurvesPage) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurvesPage) Refresh() {
	fans, curveIds, err := c.fetchCurves()
	if err != nil || fans == nil {
		fans = &map[string]*client.Curve{}
	}

	oldCIds := maps.Keys(c.curveListItemComponents)
	// remove now nonexisting entries
	for _, oldFId := range oldCIds {
		_, ok := (*fans)[oldFId]
		if !ok {
			fanListItemComponent := c.curveListItemComponents[oldFId]
			c.curveRowLayout.RemoveItem(fanListItemComponent.GetLayout())
			delete(c.curveListItemComponents, oldFId)
		}
	}

	// add new entries / update existing entries
	for _, cId := range curveIds {
		curve := (*fans)[cId]
		fanListItemComponent, ok := c.curveListItemComponents[cId]
		if ok {
			fanListItemComponent.SetCurve(curve)
			fanListItemComponent.Refresh()
		} else {
			fanListItemComponent = NewCurveListItemComponent(c.application, curve)
			c.curveListItemComponents[cId] = fanListItemComponent
			fanListItemComponent.SetCurve(curve)
			fanListItemComponent.Refresh()
			c.curveRowLayout.AddItem(fanListItemComponent.GetLayout(), 0, 1, true)
		}
	}
}
