package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"

	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type FanGraphComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout         *tview.Flex
	graphComponent *graph.GraphComponent
	rpmValues      *[]float64
	pwmValues      *[]float64
}

func NewFanGraphComponent(application *tview.Application, fan *client.Fan) *FanGraphComponent {
	rpmValues := &[]float64{}
	pwmValues := &[]float64{}
	c := &FanGraphComponent{
		application: application,
		Fan:         fan,
		rpmValues:   rpmValues,
		pwmValues:   pwmValues,
	}

	graphConfig := graph.NewGraphComponentConfig().
		WithReversedOrder().
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(true).
		WithYAxisLabelDataType(tvxwidgets.PlotYAxisLabelDataInt).
		WithOverlays(
			newCurrentRpmYAxisLabelOverlay(c.getFan),
		)

	graphComponent := graph.NewGraphComponent(
		application,
		graphConfig,
	)

	rpmProvider := graph.NewRoundedSliceSeriesValueProvider(rpmValues)
	pwmProvider := graph.NewRoundedSliceSeriesValueProvider(pwmValues)
	rpmLine := graph.NewGraphLineFromSeriesValueProvider("RPM", rpmProvider)
	pwmLine := graph.NewGraphLineFromSeriesValueProvider("PWM", pwmProvider)
	graphComponent.AddSeries(rpmLine, graph.WithLegend(graph.NewGraphSeriesLegend("RPM")))
	graphComponent.AddSeries(pwmLine, graph.WithLegend(graph.NewGraphSeriesLegend("PWM")))

	minVal := 0.0
	graphComponent.SetYMinValue(&minVal)

	c.graphComponent = graphComponent

	c.layout = c.createLayout()
	c.layout.AddItem(graphComponent.GetLayout(), 0, 1, false)

	return c
}

func (c *FanGraphComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.SetBorder(false)

	return layout
}

func (c *FanGraphComponent) getFan() *client.Fan {
	if c == nil {
		return nil
	}
	return c.Fan
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
