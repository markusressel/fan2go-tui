package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/util"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type CurveGraphComponent struct {
	application *tview.Application

	Curve *client.Curve

	layout         *tview.Flex
	bmScatterPlot  *tvxwidgets.Plot
	graphComponent *util.GraphComponent[client.Curve]
}

func NewCurveGraphComponent(application *tview.Application, curve *client.Curve) *CurveGraphComponent {

	graphComponent := util.NewGraphComponent[client.Curve](application, curve, func(c *client.Curve) float64 {
		return c.Value
	}, nil,
	)

	c := &CurveGraphComponent{
		application:    application,
		graphComponent: graphComponent,
		Curve:          curve,
	}

	c.layout = c.createLayout()
	c.layout.AddItem(graphComponent.GetLayout(), 0, 1, false)

	return c
}

func (c *CurveGraphComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.SetBorder(false)

	return layout
}

func (c *CurveGraphComponent) refresh() {
	curve := c.Curve
	if curve == nil {
		return
	}
	component := c.graphComponent
	component.InsertValue(curve)
	component.Refresh()
}

func (c *CurveGraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurveGraphComponent) SetCurve(curve *client.Curve) {
	c.Curve = curve
	c.refresh()
}

func (c *CurveGraphComponent) InsertValue(curve *client.Curve) {
	c.graphComponent.InsertValue(curve)
	c.refresh()
}

func (c *CurveGraphComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
