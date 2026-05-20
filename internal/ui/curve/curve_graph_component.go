package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"

	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type CurveGraphComponent struct {
	application *tview.Application

	Curve *client.Curve

	layout         *tview.Flex
	graphComponent *graph.GraphComponent
	values         *[]float64
}

func NewCurveGraphComponent(application *tview.Application, curve *client.Curve) *CurveGraphComponent {
	values := &[]float64{}
	c := &CurveGraphComponent{
		application: application,
		Curve:       curve,
		values:      values,
	}

	graphConfig := graph.NewGraphComponentConfig().
		WithReversedOrder().
		WithPlotColors(
			theme.Colors.Graph.Curve,
			theme.Colors.Graph.CurveMin,
			theme.Colors.Graph.CurveMax,
		).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(false).
		WithYAxisLabelDataType(tvxwidgets.PlotYAxisLabelDataInt).
		WithOverlays(
			newCurrentCurveYAxisLabelOverlay(func() *client.Curve { return c.Curve }),
		)

	graphComponent := graph.NewGraphComponent(
		application,
		graphConfig,
	)

	seriesValueProvider := graph.NewRoundedSliceSeriesValueProvider(values)
	line := graph.NewGraphLineFromSeriesValueProvider("Curve", seriesValueProvider)
	graphComponent.AddSeries(line)

	graphComponent.SetYRange(0, 255)

	c.graphComponent = graphComponent

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
	*c.values = append(*c.values, curve.Value)
	bufferSize := component.GetValueBufferSize()
	if len(*c.values) > bufferSize {
		*c.values = (*c.values)[len(*c.values)-bufferSize:]
	}
	component.Refresh()
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
