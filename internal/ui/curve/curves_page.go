package curve

import (
	"fan2go-tui/internal/client"
	uiutil "fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
	"sort"
	"strings"
)

type CurvesPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Curves map[string]*client.Curve

	layout *tview.Flex

	curveComponents     []*CurveInfoComponent
	curveGraphComponent []*CurveGraphComponent
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

	curveRowLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	curvesPageLayout.AddItem(curveRowLayout, 0, 1, true)

	curves, curveIds, err := c.fetchCurves()
	if err != nil {
		// TODO: handle error
		//c.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return curvesPageLayout
	}

	for _, id := range curveIds {
		curveColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
		uiutil.SetupWindow(curveColumnLayout, id)
		curveColumnLayout.SetTitleAlign(tview.AlignLeft)
		curveColumnLayout.SetBorder(true)
		curveRowLayout.AddItem(curveColumnLayout, 0, 1, true)

		curve := (*curves)[id]
		curveInfoComponent := NewCurveInfoComponent(c.application, curve)
		c.curveComponents = append(c.curveComponents, curveInfoComponent)
		curveInfoComponent.SetCurve(curve)
		curveInfoComponent.Refresh()
		layout := curveInfoComponent.GetLayout()
		curveColumnLayout.AddItem(layout, 0, 1, true)

		curveGraphComponent := NewCurveGraphComponent(c.application, curve)
		c.curveGraphComponent = append(c.curveGraphComponent, curveGraphComponent)
		curveGraphComponent.SetCurve(curve)
		curveGraphComponent.Refresh()
		layout = curveGraphComponent.GetLayout()
		curveColumnLayout.AddItem(layout, 0, 3, true)
	}

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

	for _, component := range c.curveGraphComponent {
		if component.Curve == nil {
			continue
		}
		curve, ok := (*curves)[component.Curve.Config.ID]
		if !ok || curve == nil {
			continue
		}
		component.SetCurve(curve)
		component.Refresh()
	}
}
