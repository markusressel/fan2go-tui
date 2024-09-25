package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
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
	graphConfig := util.NewGraphComponentConfig().
		WithReversedOrder().
		WithPlotColors(theme.Colors.Graph.Curve)
	graphComponent := util.NewGraphComponent[client.Curve](
		application,
		graphConfig,
		curve,
		[]func(*client.Curve) float64{
			func(c *client.Curve) float64 {
				return c.Value
			},
		},
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
}

func (c *CurveGraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurveGraphComponent) SetCurve(curve *client.Curve) {
	c.Curve = curve
	c.refresh()
}

func (c *CurveGraphComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
