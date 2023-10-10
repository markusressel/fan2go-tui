package curve

import (
	"fan2go-tui/internal/client"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"slices"
	"sort"
	"strings"
)

const (
	MaxVisibleItems = 3
)

type CurvesPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Curves             map[string]*client.Curve
	entryVisibilityMap []string

	layout         *tview.Flex
	curveRowLayout *tview.Flex

	curveListItemComponents map[string]*CurveListItemComponent
}

func NewCurvesPage(application *tview.Application, client client.Fan2goApiClient) CurvesPage {

	curvesPage := CurvesPage{
		application:             application,
		client:                  client,
		entryVisibilityMap:      []string{},
		curveListItemComponents: map[string]*CurveListItemComponent{},
	}

	curvesPage.layout = curvesPage.createLayout()

	return curvesPage
}

func (c *CurvesPage) createLayout() *tview.Flex {
	curvesPageLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	c.curveRowLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	curvesPageLayout.AddItem(c.curveRowLayout, 0, 1, true)

	return curvesPageLayout
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

func (c *CurvesPage) Refresh() error {
	curves, curveIds, err := c.fetchCurves()
	if err != nil || curves == nil {
		curves = &map[string]*client.Curve{}
	}

	if len(c.entryVisibilityMap) < MaxVisibleItems {
		c.entryVisibilityMap = curveIds[0:MaxVisibleItems]
	}

	var visibleCurveIds []string
	// filter currently invisible curves
	for _, curveId := range curveIds {
		isCurveIdVisible := slices.ContainsFunc(visibleCurveIds, func(visibleCurveId string) bool {
			return curveId == visibleCurveId
		})
		if isCurveIdVisible {
			visibleCurveIds = append(visibleCurveIds, curveId)
		}
	}
	c.entryVisibilityMap = visibleCurveIds

	// remove now nonexisting entries
	oldCIds := maps.Keys(c.curveListItemComponents)
	for _, oldCId := range oldCIds {
		_, ok := (*curves)[oldCId]
		if !ok {
			curveListItemComponent := c.curveListItemComponents[oldCId]
			c.curveRowLayout.RemoveItem(curveListItemComponent.GetLayout())
			delete(c.curveListItemComponents, oldCId)
		}
	}

	// add new entries / update existing entries
	for _, cId := range curveIds {
		curve := (*curves)[cId]
		curveListItemComponent, ok := c.curveListItemComponents[cId]
		if ok {
			curveListItemComponent.SetCurve(curve)
		} else {
			curveListItemComponent = NewCurveListItemComponent(c.application, curve)
			c.curveListItemComponents[cId] = curveListItemComponent
			curveListItemComponent.SetCurve(curve)
			c.curveRowLayout.AddItem(curveListItemComponent.GetLayout(), 0, 1, true)
		}
	}

	// update visibility
	for curveId, oldCurveListItemComponent := range c.curveListItemComponents {
		curveIdVisible := slices.ContainsFunc(c.entryVisibilityMap, func(visibleCurveId string) bool {
			return curveId == visibleCurveId
		})

		if curveIdVisible == false {
			c.curveRowLayout.RemoveItem(oldCurveListItemComponent.GetLayout())
		}
	}

	return err
}
