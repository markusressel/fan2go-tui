package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type SensorGraphComponent struct {
	application *tview.Application

	Sensor *client.Sensor

	layout         *tview.Flex
	graphComponent *graph.GraphComponent
	graphBar       *graph.GraphBar
	values         *[]float64
}

func NewSensorGraphComponent(application *tview.Application, sensor *client.Sensor) *SensorGraphComponent {
	values := &[]float64{}
	c := &SensorGraphComponent{
		application: application,
		Sensor:      sensor,
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
	graphComponent.AddSeries(bar)

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
	if c == nil {
		return nil
	}
	return c.Sensor
}

func (c *SensorGraphComponent) refresh() {
	sensor := c.Sensor
	if sensor == nil {
		return
	}
	component := c.graphComponent
	value := sensor.MovingAvg / 1000
	*c.values = append(*c.values, value)

	bufferSize := component.GetValueBufferSize()
	if len(*c.values) > bufferSize {
		trimmed := (*c.values)[len(*c.values)-bufferSize:]
		*c.values = trimmed
	}

	component.Refresh()
}

func (c *SensorGraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *SensorGraphComponent) SetSensor(sensor *client.Sensor) {
	c.Sensor = sensor
	c.refresh()
}

func (c *SensorGraphComponent) SetTitle(label string) {
	c.graphComponent.SetTitle(label)
}
