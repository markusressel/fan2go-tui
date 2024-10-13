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
	yFunc1 func(x float64) float64
}

func NewFanRpmCurveComponent(application *tview.Application, fan *client.Fan) *FanRpmCurveComponent {
	fanCurveData := *fan.FanCurveData
	pwmKeys := maps.Keys(fanCurveData)
	slices.Sort(pwmKeys)

	horizontalStretchFactor := 1.0
	verticalStretchFactor := 1.0
	xOffset := 0.0
	yOffset := 0.0

	// TODO: set xAxisZoomFactor based on available graph width
	xAxisZoomFactor := 0.7
	xAxisShift := 0.0

	xFunc1 := func(i int) float64 {
		if i >= len(pwmKeys) {
			return math.NaN()
		}
		return (float64(pwmKeys[i]) / xAxisZoomFactor) + xAxisShift
	}
	fFunc := func(x float64) float64 {
		val, ok := fanCurveData[int(math.Round(x))]
		if !ok {
			return math.NaN()
		} else {
			return val
		}
	}
	yFunc1 := func(x float64) float64 {
		return (fFunc((x+xOffset)/horizontalStretchFactor) + yOffset) * verticalStretchFactor
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

	c := &FanRpmCurveComponent{
		application:    application,
		graphComponent: graphComponent,
		Fan:            fan,

		pwmKeys: pwmKeys,

		xFunc1: xFunc1,
		fFunc:  fFunc,
		yFunc1: yFunc1,
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

	fanRpmGraphData, err := c.computeFanRpmGraphData(c.Fan)
	if err != nil {
		return
	}
	c.graphComponent.SetRawData(fanRpmGraphData)
}

func (c *FanRpmCurveComponent) computeFanRpmGraphData(fan *client.Fan) ([][]float64, error) {
	graphData := make([][]float64, 1)

	computeDataArray := func() [][]float64 {
		n := len(c.pwmKeys)
		data := make([][]float64, 1)
		data[0] = make([]float64, n)
		for i := 0; i < n; i++ {
			xVal := c.xFunc1(i)
			yVal := c.yFunc1(xVal)
			data[0][i] = yVal
		}

		return data
	}

	data := computeDataArray()

	c.graphComponent.UpdateValueBufferSize()
	//totalRange := c.graphComponent.GetValueBufferSize()
	graphData[0] = data[0]

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
