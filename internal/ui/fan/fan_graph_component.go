package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"

	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type FanGraphComponent struct {
	application *tview.Application

	FanState *state.FanState

	layout         *tview.Flex
	graphComponent *graph.GraphComponent
	rpmValues      *[]float64
	pwmValues      *[]float64
}

func NewFanGraphComponent(application *tview.Application, fanState *state.FanState) *FanGraphComponent {
	rpmValues := &[]float64{}
	pwmValues := &[]float64{}
	c := &FanGraphComponent{
		application: application,
		FanState:    fanState,
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
	if c == nil || c.FanState == nil {
		return nil
	}
	return c.FanState.Fan
}

func (c *FanGraphComponent) refresh() {
	if c.FanState == nil || c.FanState.Fan == nil {
		return
	}
	component := c.graphComponent
	bufferSize := component.GetValueBufferSize()

	if bufferSize > 0 {
		rpmHistory := c.FanState.RpmValues
		if len(rpmHistory) > bufferSize {
			*c.rpmValues = append([]float64(nil), rpmHistory[len(rpmHistory)-bufferSize:]...)
		} else {
			*c.rpmValues = append([]float64(nil), rpmHistory...)
		}

		pwmHistory := c.FanState.PwmValues
		if len(pwmHistory) > bufferSize {
			*c.pwmValues = append([]float64(nil), pwmHistory[len(pwmHistory)-bufferSize:]...)
		} else {
			*c.pwmValues = append([]float64(nil), pwmHistory...)
		}
	}

	component.Refresh()
}

func (c *FanGraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanGraphComponent) SetFan(fanState *state.FanState) {
	c.FanState = fanState
	c.refresh()
}

func (c *FanGraphComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
