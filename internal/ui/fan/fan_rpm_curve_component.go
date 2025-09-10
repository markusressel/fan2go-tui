package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
	"fan2go-tui/internal/ui/util"
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
	bmScatterPlot  *tvxwidgets.Plot
	graphComponent *util.GraphComponent[client.Fan]

	xFunc1 func(i int) float64
	fFunc  func(x float64) float64
}

func NewFanRpmCurveComponent(application *tview.Application, fan *client.Fan) *FanRpmCurveComponent {
	fanCurveData := *fan.FanCurveData
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

	rpmGraphLine := util.NewGraphLine(
		"RPM",
		xFunc1,
		fFunc,
		xLabelFunc1,
	)

	graphConfig := util.NewGraphComponentConfig().
		WithReversedOrder().
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(true).
		WithDrawXAxisLabel(true).
		WithYAxisLabelDataType(tvxwidgets.PlotYAxisLabelDataInt)

	graphComponent := util.NewGraphComponent[client.Fan](
		application,
		graphConfig,
		fan,
		[]func(*client.Fan) float64{},
	)

	graphComponent.AddLine(rpmGraphLine)
	graphComponent.SetXRange(0, 255)

	c := &FanRpmCurveComponent{
		application:    application,
		graphComponent: graphComponent,
		Fan:            fan,

		pwmKeys: pwmKeys,

		xFunc1: xFunc1,
		fFunc:  fFunc,
	}

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

	c.graphComponent.ZoomToRangeX(0, 255)
	c.graphComponent.Refresh()
}

func (c *FanRpmCurveComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanRpmCurveComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
	c.refresh()
}

func (c *FanRpmCurveComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
