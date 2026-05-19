package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"
	"math"
	"strconv"

	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type FanGraphComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout         *tview.Flex
	bmScatterPlot  *tvxwidgets.Plot
	graphComponent *graph.GraphComponent[client.Fan]
	rpmValues      *[]float64
	pwmValues      *[]float64
}

func NewFanGraphComponent(application *tview.Application, fan *client.Fan) *FanGraphComponent {
	graphConfig := graph.NewGraphComponentConfigFor(fan).
		WithReversedOrder().
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(true).
		WithYAxisLabelDataType(tvxwidgets.PlotYAxisLabelDataInt)

	graphComponent := graph.NewGraphComponent[client.Fan](
		application,
		graphConfig,
		fan,
	)

	rpmValues := &[]float64{}
	pwmValues := &[]float64{}
	xFunc := func(i int) float64 { return float64(i) }
	xLabelFunc := func(i int, x float64) string { return strconv.Itoa(int(math.Round(x))) }
	rpmLine := graph.NewGraphLine(
		"RPM",
		xFunc,
		func(x float64) float64 {
			idx := int(math.Round(x))
			if idx < 0 || idx >= len(*rpmValues) {
				return math.NaN()
			}
			return (*rpmValues)[idx]
		},
		xLabelFunc,
	)
	pwmLine := graph.NewGraphLine(
		"PWM",
		xFunc,
		func(x float64) float64 {
			idx := int(math.Round(x))
			if idx < 0 || idx >= len(*pwmValues) {
				return math.NaN()
			}
			return (*pwmValues)[idx]
		},
		xLabelFunc,
	)
	graphComponent.AddSeries(rpmLine)
	graphComponent.AddSeries(pwmLine)

	minVal := 0.0
	graphComponent.SetYMinValue(&minVal)

	c := &FanGraphComponent{
		application:    application,
		graphComponent: graphComponent,
		Fan:            fan,
		rpmValues:      rpmValues,
		pwmValues:      pwmValues,
	}

	c.layout = c.createLayout()
	c.layout.AddItem(graphComponent.GetLayout(), 0, 1, false)

	return c
}

func (c *FanGraphComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.SetBorder(false)

	return layout
}

func (c *FanGraphComponent) refresh() {
	fan := c.Fan
	if fan == nil {
		return
	}
	component := c.graphComponent
	*c.rpmValues = append(*c.rpmValues, float64(fan.Rpm))
	*c.pwmValues = append(*c.pwmValues, float64(fan.Pwm))
	bufferSize := component.GetValueBufferSize()
	if len(*c.rpmValues) > bufferSize {
		*c.rpmValues = (*c.rpmValues)[len(*c.rpmValues)-bufferSize:]
	}
	if len(*c.pwmValues) > bufferSize {
		*c.pwmValues = (*c.pwmValues)[len(*c.pwmValues)-bufferSize:]
	}
	component.Refresh()
}

func (c *FanGraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanGraphComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
	c.refresh()
}

func (c *FanGraphComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
