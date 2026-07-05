package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"

	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type CurveGraphComponent struct {
	application *tview.Application

	CurveState *state.CurveState

	layout         *tview.Flex
	graphComponent *graph.GraphComponent
	values         *[]float64
}

func NewCurveGraphComponent(application *tview.Application, curveState *state.CurveState) *CurveGraphComponent {
	values := &[]float64{}
	c := &CurveGraphComponent{
		application: application,
		CurveState:  curveState,
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
			newCurrentCurveYAxisLabelOverlay(c.getCurve),
		)

	graphComponent := graph.NewGraphComponent(
		application,
		graphConfig,
	)

	seriesValueProvider := graph.NewRoundedSliceSeriesValueProvider(values)
	line := graph.NewGraphLineFromSeriesValueProvider("Curve", seriesValueProvider)
	graphComponent.AddSeries(line, graph.WithLegend(graph.NewGraphSeriesLegend("Value")))

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

func (c *CurveGraphComponent) getCurve() *client.Curve {
	if c == nil || c.CurveState == nil {
		return nil
	}
	return c.CurveState.Curve
}

func (c *CurveGraphComponent) refresh() {
	if c.CurveState == nil || c.CurveState.Curve == nil {
		return
	}
	component := c.graphComponent
	bufferSize := component.GetValueBufferSize()

	if bufferSize > 0 {
		history := c.CurveState.Values
		if len(history) > bufferSize {
			*c.values = append([]float64(nil), history[len(history)-bufferSize:]...)
		} else {
			*c.values = append([]float64(nil), history...)
		}
	}

	component.Refresh()
}

func (c *CurveGraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurveGraphComponent) SetCurve(curveState *state.CurveState) {
	c.CurveState = curveState
	c.refresh()
}

func (c *CurveGraphComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
