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

	xLabelFunc1 := func(i int) string {
		xVal := xFunc1(i)
		labelVal := math.Round(xVal)
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
	)

	graphConfig := util.NewGraphComponentConfig().
		WithReversedOrder().
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(true).
		WithXMax(255).
		WithDrawXAxisLabel(true).
		WithXAxisLabelFunc(xLabelFunc1)

	graphComponent := util.NewGraphComponent[client.Fan](
		application,
		graphConfig,
		fan,
		[]func(*client.Fan) float64{
			func(c *client.Fan) float64 {
				return 0
			},
		},
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

	fanRpmGraphData, err := c.computeFanRpmGraphData()
	if err != nil {
		return
	}

	c.graphComponent.UpdateValueBufferSize()
	totalRange := c.graphComponent.GetValueBufferSize()
	c.graphComponent.SetXAxisZoomFactor(float64(totalRange) / 255.0)

	c.graphComponent.SetRawData(fanRpmGraphData)
}

func (c *FanRpmCurveComponent) computeFanRpmGraphData() ([][]float64, error) {
	graphData := make([][]float64, 1)

	for _, line := range c.graphComponent.GetLines() {
		n := 200
		data := make([]float64, n)
		for i := 0; i < n; i++ {
			xVal := line.GetX(i)
			yVal := line.GetY(xVal)
			data[i] = yVal
		}
		graphData = append(graphData, data)
	}

	return graphData, nil
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
