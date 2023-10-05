package curve

import (
	"fan2go-tui/internal/client"
	"github.com/rivo/tview"
	"sort"
	"strings"
)

type CurvesPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Curves map[string]*client.Curve

	layout *tview.Flex

	curveComponents      []*CurveComponent
	curveGraphsComponent *CurveGraphsComponent
	curveGraphComponent  *CurveGraphComponent
}

func NewCurvesPage(application *tview.Application, client client.Fan2goApiClient) CurvesPage {

	curvesPage := CurvesPage{
		application: application,
		client:      client,
	}

	curvesPage.layout = curvesPage.createLayout()

	return curvesPage
}

func (c *CurvesPage) createLayout() *tview.Flex {

	curvesPageLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	curveInfoLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	curvesPageLayout.AddItem(curveInfoLayout, 0, 1, true)

	curves, err := c.client.GetCurves()
	if err != nil {
		// TODO: handle error
		//c.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return curvesPageLayout
	}
	var curveComponents []*CurveComponent
	curvesIds := []string{}
	for _, c := range *curves {
		curvesIds = append(curvesIds, c.Config.ID)
	}

	sort.SliceStable(curvesIds, func(i, j int) bool {
		a := curvesIds[i]
		b := curvesIds[j]

		result := strings.Compare(strings.ToLower(a), strings.ToLower(b))

		if result <= 0 {
			return true
		} else {
			return false
		}
	})

	for _, id := range curvesIds {
		curve := (*curves)[id]

		curveComponent := NewCurveComponent(c.application, curve)
		curveComponents = append(curveComponents, curveComponent)
		curveComponent.SetCurve(curve)
		curveComponent.Refresh()
		layout := curveComponent.GetLayout()
		curveInfoLayout.AddItem(layout, 0, 1, true)
	}
	c.curveComponents = curveComponents

	curveGraphsComponent := NewCurveGraphsComponent(c.application)
	c.curveGraphsComponent = curveGraphsComponent
	//curveComponents = append(curveComponents, curveGaphsComponent)

	// update overview
	curveList := []*client.Curve{}
	for _, f := range *curves {
		curveList = append(curveList, f)
	}

	curveGraphsComponent.SetCurves(curveList)
	curveGraphsComponent.Refresh()
	layout := curveGraphsComponent.GetLayout()
	curvesPageLayout.AddItem(layout, 0, 1, true)

	return curvesPageLayout
}

func (c *CurvesPage) fetchCurves() (*map[string]*client.Curve, error) {
	return c.client.GetCurves()
}

func (c *CurvesPage) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurvesPage) Refresh() {
	curves, err := c.client.GetCurves()
	if err != nil {
		return
	}

	for _, component := range c.curveComponents {
		curve, ok := (*curves)[component.Curve.Config.ID]
		if !ok {
			continue
		}
		component.SetCurve(curve)
		component.Refresh()
	}
}
