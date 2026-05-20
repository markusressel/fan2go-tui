package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"
	"math"
	"slices"
	"strconv"

	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
)

type FanRpmCurveComponent struct {
	application *tview.Application

	Fan     *client.Fan
	pwmKeys []int

	layout         *tview.Flex
	graphComponent *graph.GraphComponent
}

func NewFanRpmCurveComponent(application *tview.Application, fan *client.Fan) *FanRpmCurveComponent {
	fanCurveData := map[int]float64{}
	if fan.FanCurveData != nil {
		fanCurveData = *fan.FanCurveData
	}
	pwmKeys := maps.Keys(fanCurveData)
	slices.Sort(pwmKeys)

	xFunc1 := func(i int) float64 {
		if i < 0 || i >= len(pwmKeys) {
			return math.NaN()
		}
		return float64(pwmKeys[i])
	}
	fFunc := func(x float64) float64 {
		val, ok := fanCurveData[int(math.Floor(x))]
		if !ok {
			return math.NaN()
		} else {
			return val
		}
	}

	xLabelFunc1 := func(i int, x float64) string {
		labelVal := int(math.Round(x))
		label := strconv.Itoa(labelVal)
		return label
	}

	rpmGraphLine := graph.NewGraphLine(
		"RPM",
		xFunc1,
		fFunc,
		xLabelFunc1,
	)

	c := &FanRpmCurveComponent{
		application: application,
		Fan:         fan,

		pwmKeys: pwmKeys,
	}

	graphConfig := graph.NewGraphComponentConfig().
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(true).
		WithDrawXAxisLabel(true).
		WithYAxisLabelDataType(tvxwidgets.PlotYAxisLabelDataInt).
		WithOverlays(
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
	graphComponent.SetXRange(0, 255)
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

	c.graphComponent.ZoomToRangeX(0, 255)
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
