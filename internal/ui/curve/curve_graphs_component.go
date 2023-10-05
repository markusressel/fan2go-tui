package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/util"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type CurveGraphsComponent struct {
	application *tview.Application

	Curves []*client.Curve

	layout          *tview.Flex
	bmScatterPlot   *tvxwidgets.Plot
	graphComponents map[string]*util.GraphComponent[client.Curve]
}

func NewCurveGraphsComponent(application *tview.Application) *CurveGraphsComponent {
	c := &CurveGraphsComponent{
		application:     application,
		Curves:          []*client.Curve{},
		graphComponents: map[string]*util.GraphComponent[client.Curve]{},
	}

	c.layout = c.createLayout()

	return c
}

func (c *CurveGraphsComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.SetBorder(false)

	return layout
}

func (c *CurveGraphsComponent) Refresh() {
	for _, curve := range c.Curves {
		component, ok := c.graphComponents[curve.Config.ID]
		if !ok {
			component = util.NewGraphComponent[client.Curve](c.application, curve, func(c *client.Curve) float64 {
				return float64(c.Value)
			}, func(c *client.Curve) float64 {
				return 0
			},
			)
			c.graphComponents[curve.Config.ID] = component
			c.layout.AddItem(component.GetLayout(), 0, 1, false)
			component.InsertValue(curve)
			component.SetTitle(curve.Config.ID)
			component.Refresh()
		} else {
			component.InsertValue(curve)
			component.Refresh()
		}
	}
}

func (c *CurveGraphsComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurveGraphsComponent) SetCurves(curves []*client.Curve) {
	c.Curves = curves
	c.Refresh()
}
