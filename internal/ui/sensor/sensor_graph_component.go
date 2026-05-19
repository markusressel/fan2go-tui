package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"
	"math"
	"strconv"

	"github.com/gdamore/tcell/v2"
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
