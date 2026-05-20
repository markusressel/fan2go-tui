package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"
	"math"

	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

const (
	fanRpmCurveMinX = 0.0
	fanRpmCurveMaxX = 255.0
)

type FanRpmCurveComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout         *tview.Flex
	graphComponent *graph.GraphComponent
}

func NewFanRpmCurveComponent(application *tview.Application, fan *client.Fan) *FanRpmCurveComponent {
	fanCurveData := map[int]float64{}
	if fan.FanCurveData != nil {
		fanCurveData = *fan.FanCurveData
	}

	seriesValueProvider := graph.NewDiscreteIntSeriesValueProvider(fanCurveData)
	rpmGraphLine := graph.NewGraphLineFromSeriesValueProvider("RPM", seriesValueProvider)

	c := &FanRpmCurveComponent{
		application: application,
		Fan:         fan,
	}

	graphConfig := graph.NewGraphComponentConfig().
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(true).
		WithDrawXAxisLabel(true).
		WithYAxisLabelDataType(tvxwidgets.PlotYAxisLabelDataInt).
		WithOverlays(
			newCurrentRpmYAxisLabelOverlay(func() *client.Fan { return c.Fan }),
			graph.VLine(
				func() float64 {
					fan := c.Fan
					if fan == nil {
						return math.NaN()
					}
					return float64(fan.Pwm)
				},
			).WithColor(theme.Colors.Graph.CurrentPwmLine),
			graph.Marker(
				func() graph.XY {
					fan := c.Fan
					if fan == nil {
						return graph.XY{X: math.NaN(), Y: math.NaN()}
					}
					return graph.XY{X: float64(fan.Pwm), Y: float64(fan.Rpm)}
				},
			).WithColor(theme.Colors.Graph.CurrentRpmMarker),
		)

	graphComponent := graph.NewGraphComponent(
		application,
		graphConfig,
	)

	graphComponent.AddSeries(rpmGraphLine)
	graphComponent.SetXRange(fanRpmCurveMinX, fanRpmCurveMaxX)
	c.graphComponent = graphComponent

	c.layout = c.createLayout()
	c.layout.AddItem(graphComponent.GetLayout(), 0, 1, false)

	return c
}

func (c *FanRpmCurveComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	layout.SetBorder(false)
	return layout
}

func (c *FanRpmCurveComponent) refresh() {
	fan := c.Fan
	if fan == nil {
		return
	}
	if c.graphComponent == nil {
		return
	}

	// First refresh updates buffer/layout, second uses fixed x-range zoom for final draw.
	c.graphComponent.Refresh()
	c.graphComponent.ZoomToRangeX(fanRpmCurveMinX, fanRpmCurveMaxX)
	c.graphComponent.Refresh()
}

func (c *FanRpmCurveComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanRpmCurveComponent) SetFan(fan *client.Fan) {
	if c == nil {
		return
	}
	c.Fan = fan
	c.refresh()
}

func (c *FanRpmCurveComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
