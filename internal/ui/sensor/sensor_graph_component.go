package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SensorGraphComponent struct {
	application *tview.Application

	SensorState *state.SensorState

	layout         *tview.Flex
	graphComponent *graph.GraphComponent
	graphBar       *graph.GraphBar
	values         *[]float64
}

func NewSensorGraphComponent(application *tview.Application, sensorState *state.SensorState) *SensorGraphComponent {
	values := &[]float64{}
	c := &SensorGraphComponent{
		application: application,
		SensorState: sensorState,
		values:      values,
	}

	graphConfig := graph.NewGraphComponentConfig().
		WithReversedOrder().
		WithPlotColors(
			theme.Colors.Graph.Sensor,
			theme.Colors.Graph.SensorMin,
			theme.Colors.Graph.SensorMax,
		).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(false).
		WithOverlays(
			newCurrentSensorYAxisLabelOverlay(c.getSensor),
		)

	graphComponent := graph.NewGraphComponent(
		application,
		graphConfig,
	)

	seriesValueProvider := graph.NewRoundedSliceSeriesValueProvider(values)
	bar := graph.NewGraphBarFromSeriesValueProvider("Sensor", seriesValueProvider)
	bar.SetColor(theme.Colors.Graph.Sensor)
	bar.WithGradient(func(yMin, yMax float64) []graph.GraphBarGradientStop {
		rangeY := yMax - yMin
		if rangeY <= 0 {
			rangeY = 1
		}

		return []graph.GraphBarGradientStop{
			{YValue: yMin, Color: tcell.NewRGBColor(0, 120, 255)},
			{YValue: yMin + rangeY*(1.0/3.0), Color: tcell.NewRGBColor(0, 200, 80)},
			{YValue: yMin + rangeY*(2.0/3.0), Color: tcell.NewRGBColor(240, 220, 0)},
			{YValue: yMax, Color: tcell.NewRGBColor(220, 40, 30)},
		}
	})
	graphComponent.AddSeries(bar, graph.WithLegend(graph.NewGraphSeriesLegend("Temperature").WithUnit("°C")))

	graphComponent.SetYRange(0, 100)

	c.graphComponent = graphComponent
	c.graphBar = bar

	c.layout = c.createLayout()
	c.layout.AddItem(graphComponent.GetLayout(), 0, 1, false)

	return c
}

func (c *SensorGraphComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	return layout
}

func (c *SensorGraphComponent) getSensor() *client.Sensor {
	if c == nil || c.SensorState == nil {
		return nil
	}
	return c.SensorState.Sensor
}

func (c *SensorGraphComponent) refresh() {
	if c.SensorState == nil || c.SensorState.Sensor == nil {
		return
	}
	component := c.graphComponent
	bufferSize := component.GetValueBufferSize()

	if bufferSize > 0 {
		history := c.SensorState.Values
		if len(history) > bufferSize {
			*c.values = append([]float64(nil), history[len(history)-bufferSize:]...)
		} else {
			*c.values = append([]float64(nil), history...)
		}
	}

	component.Refresh()
}

func (c *SensorGraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorGraphComponent) SetSensor(sensorState *state.SensorState) {
	c.SensorState = sensorState
	c.refresh()
}

func (c *SensorGraphComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
