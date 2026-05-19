package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"
	"math"
	"strconv"

	"github.com/rivo/tview"
)

type SensorGraphComponent struct {
	application *tview.Application

	Sensor *client.Sensor

	layout         *tview.Flex
	graphComponent *graph.GraphComponent[client.Sensor]
	graphBar       *graph.GraphBar
	values         *[]float64
}

func NewSensorGraphComponent(application *tview.Application, sensor *client.Sensor) *SensorGraphComponent {
	graphConfig := graph.NewGraphComponentConfigFor(sensor).
		WithReversedOrder().
		WithPlotColors(
			theme.Colors.Graph.Sensor,
			theme.Colors.Graph.SensorMin,
			theme.Colors.Graph.SensorMax,
		).
		WithYAxisAutoScaleMin(false).
		WithYAxisAutoScaleMax(false)

	graphComponent := graph.NewGraphComponent[client.Sensor](
		application,
		graphConfig,
		sensor,
		[]func(*client.Sensor) float64{},
	)

	values := &[]float64{}
	xFunc := func(i int) float64 {
		return float64(i)
	}
	fFunc := func(x float64) float64 {
		idx := int(math.Round(x))
		if idx < 0 || idx >= len(*values) {
			return math.NaN()
		}
		return (*values)[idx]
	}
	xLabelFunc := func(i int, x float64) string {
		return strconv.Itoa(int(math.Round(x)))
	}

	bar := graph.NewGraphBar("Sensor", xFunc, fFunc, xLabelFunc)
	bar.SetColor(theme.Colors.Graph.Sensor)
	graphComponent.AddSeries(bar)

	graphComponent.SetYRange(0, 100)

	c := &SensorGraphComponent{
		application:    application,
		graphComponent: graphComponent,
		graphBar:       bar,
		Sensor:         sensor,
		values:         values,
	}

	c.layout = c.createLayout()
	c.layout.AddItem(graphComponent.GetLayout(), 0, 1, false)

	return c
}

func (c *SensorGraphComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	return layout
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
