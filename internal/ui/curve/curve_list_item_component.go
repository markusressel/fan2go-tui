package curve

import (
	"fan2go-tui/internal/client"
	uiutil "fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
)

type CurveListItemComponent struct {
	application *tview.Application

	Curve *client.Curve

	layout *tview.Flex

	curveInfoComponent  *CurveInfoComponent
	curveGraphComponent *CurveGraphComponent
}

func NewCurveListItemComponent(application *tview.Application, curve *client.Curve) *CurveListItemComponent {
	c := &CurveListItemComponent{
		application: application,
		Curve:       curve,
	}

	c.layout = c.createLayout()

	return c
}

func (c *CurveListItemComponent) createLayout() *tview.Flex {
	rootLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	curveColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	uiutil.SetupWindow(curveColumnLayout, c.Curve.Config.ID)
	curveColumnLayout.SetTitleAlign(tview.AlignLeft)
	curveColumnLayout.SetBorder(true)
	rootLayout.AddItem(curveColumnLayout, 0, 1, true)

	f := c.Curve
	curveInfoComponent := NewCurveInfoComponent(c.application, f)
	c.curveInfoComponent = curveInfoComponent
	curveInfoComponent.SetCurve(f)
	layout := curveInfoComponent.GetLayout()
	curveColumnLayout.AddItem(layout, 0, 1, true)

	curveGraphComponent := NewCurveGraphComponent(c.application, f)
	c.curveGraphComponent = curveGraphComponent
	curveGraphComponent.SetCurve(f)
	layout = curveGraphComponent.GetLayout()
	curveColumnLayout.AddItem(layout, 0, 3, true)

	return rootLayout
}

func (c *CurveListItemComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurveListItemComponent) SetCurve(curve *client.Curve) {
	c.Curve = curve
	c.refresh()
}

func (c *CurveListItemComponent) refresh() {
	c.curveInfoComponent.SetCurve(c.Curve)
	c.curveGraphComponent.SetCurve(c.Curve)
}
