package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
	"fan2go-tui/internal/ui/util"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"math"
	"slices"
	"strconv"
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
		if i >= len(pwmKeys) {
			return math.NaN()
		}
		return float64(pwmKeys[i])
	}
	fFunc := func(x float64) float64 {
		val, ok := fanCurveData[int(math.Round(x))]
		if !ok {
			return math.NaN()
		} else {
			return val
		}
	}

	xLabelFunc1 := func(i int, x float64) string {
		labelVal := math.Round(x)
		if labelVal > 255 {
			return ""
		}
		label := strconv.Itoa(int(labelVal))
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
		WithDrawXAxisLabel(true)

	graphComponent := util.NewGraphComponent[client.Fan](
		application,
		graphConfig,
		fan,
		[]func(*client.Fan) float64{},
	)

	graphComponent.AddLine(rpmGraphLine)

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

	c.graphComponent.UpdateValueBufferSize()
	_, _, width, _ := c.graphComponent.GetLayout().GetInnerRect()
	//totalRange := c.graphComponent.GetValueBufferSize()
	totalRange := math.Max(1, float64(width-5))
	newXAxisZoomFactor := 1 / (255.0 / float64(totalRange))
	//newXAxisZoomFactor = 1.0
	c.graphComponent.SetXAxisZoomFactor(newXAxisZoomFactor)
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
