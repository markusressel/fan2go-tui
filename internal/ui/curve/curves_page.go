package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"sort"
	"strings"
)

type CurvesPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Curves             map[string]*client.Curve
	entryVisibilityMap []string

	layout         *tview.Flex
	curveRowLayout *tview.Flex

	curveList *util.ListComponent[CurveListItemComponent]

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
	curvesPageLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	curveListComponent := util.NewListComponent[CurveListItemComponent](
		c.application,
		func(row int, entry *CurveListItemComponent) (layout tview.Primitive) {
			return entry.GetLayout()
		},
	)
	c.curveList = curveListComponent
	curvesPageLayout.AddItem(c.curveList.GetLayout(), 0, 1, true)

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

	var curveListItemsComponents []*CurveListItemComponent

	oldFIds := maps.Keys(c.curveListItemComponents)
	// remove now nonexisting entries
	for _, oldFId := range oldFIds {
		_, ok := (*curves)[oldFId]
		if !ok {
			curveListItemComponent := c.curveListItemComponents[oldFId]
			c.curveRowLayout.RemoveItem(curveListItemComponent.GetLayout())
			delete(c.curveListItemComponents, oldFId)
		}
	}

	// add new entries / update existing entries
	for _, fId := range curveIds {
		curve := (*curves)[fId]
		curveListItemComponent, ok := c.curveListItemComponents[fId]
		if ok {
			curveListItemComponent.SetCurve(curve)
			curveListItemsComponents = append(curveListItemsComponents, curveListItemComponent)
		} else {
			curveListItemComponent = NewCurveListItemComponent(c.application, curve)
			c.curveListItemComponents[fId] = curveListItemComponent
			curveListItemComponent.SetCurve(curve)
			curveListItemsComponents = append(curveListItemsComponents, curveListItemComponent)
		}
	}

	c.curveList.SetData(curveListItemsComponents)

	return err
}
