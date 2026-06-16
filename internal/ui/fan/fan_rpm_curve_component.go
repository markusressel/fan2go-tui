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
	fanRpmCurveMinX             = 0.0
	fanRpmCurveMaxX             = 255.0
	fanRpmCurveTrailHistorySize = 200
)

type FanRpmCurveComponent struct {
	application *tview.Application

	Fan     *client.Fan
	history []graph.XY

	layout              *tview.Flex
	graphComponent      *graph.GraphComponent
	seriesValueProvider *graph.DiscreteIntSeriesValueProvider
}

func NewFanRpmCurveComponent(application *tview.Application, fan *client.Fan) *FanRpmCurveComponent {
	fanCurveData := map[int]float64{}
	if fan != nil && fan.FanCurveData != nil {
		fanCurveData = *fan.FanCurveData
	}

	seriesValueProvider := graph.NewDiscreteIntSeriesValueProvider(fanCurveData)
	rpmGraphLine := graph.NewGraphLineFromSeriesValueProvider("RPM", seriesValueProvider)

	c := &FanRpmCurveComponent{
		application:         application,
		Fan:                 fan,
		history:             make([]graph.XY, 0, fanRpmCurveTrailHistorySize),
		seriesValueProvider: seriesValueProvider,
	}

	graphConfig := graph.NewGraphComponentConfig().
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(true).
		WithDrawXAxisLabel(true).
		WithLegendCorner(graph.LegendCornerBottomRight).
		WithYAxisLabelDataType(tvxwidgets.PlotYAxisLabelDataInt).
		WithOverlays(
			newCurrentRpmYAxisLabelOverlay(c.getFan),
			newCurrentPwmXAxisLabelOverlay(c.getFan),
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
			graph.Trail(c.getHistory).
				WithColor(theme.Colors.Graph.CurrentPwmLine).
				WithMaxPoints(fanRpmCurveTrailHistorySize),
		)

	graphComponent := graph.NewGraphComponent(
		application,
		graphConfig,
	)

	graphComponent.AddSeries(rpmGraphLine, graph.WithLegend(graph.NewGraphSeriesLegend("RPM / PWM")))
	graphComponent.SetXRange(fanRpmCurveMinX, fanRpmCurveMaxX)
	graphComponent.ZoomToRangeX(fanRpmCurveMinX, fanRpmCurveMaxX)
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

func (c *FanRpmCurveComponent) getFan() *client.Fan {
	if c == nil {
		return nil
	}
	return c.Fan
}

func (c *FanRpmCurveComponent) getHistory() []graph.XY {
	if c == nil || len(c.history) == 0 {
		return nil
	}
	historyCopy := make([]graph.XY, len(c.history))
	copy(historyCopy, c.history)
	return historyCopy
}

func (c *FanRpmCurveComponent) appendHistory(pwm, rpm float64) {
	c.history = append(c.history, graph.XY{X: pwm, Y: rpm})
	if len(c.history) > fanRpmCurveTrailHistorySize {
		c.history = c.history[len(c.history)-fanRpmCurveTrailHistorySize:]
	}
}

func (c *FanRpmCurveComponent) refresh() {
	c.syncCurveData()

	fan := c.Fan
	if fan == nil {
		return
	}
	if c.graphComponent == nil {
		return
	}
	c.appendHistory(float64(fan.Pwm), float64(fan.Rpm))

	c.graphComponent.Refresh()
}

func (c *FanRpmCurveComponent) syncCurveData() {
	if c == nil || c.seriesValueProvider == nil {
		return
	}

	curveData := map[int]float64{}
	if c.Fan != nil && c.Fan.FanCurveData != nil {
		curveData = *c.Fan.FanCurveData
	}

	c.seriesValueProvider.SetValues(curveData)
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
