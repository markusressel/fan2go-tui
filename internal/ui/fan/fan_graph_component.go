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
	bmScatterPlot  *tvxwidgets.Plot
	graphComponent *graph.GraphComponent[client.Fan]
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
		[]func(*client.Fan) float64{
			func(c *client.Fan) float64 {
				return float64(c.Rpm)
			},
			func(c *client.Fan) float64 {
				return float64(c.Pwm)
			},
		},
	)

	minVal := 0.0
	graphComponent.SetYMinValue(&minVal)

	c := &FanGraphComponent{
		application:    application,
		graphComponent: graphComponent,
		Fan:            fan,
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
	component.InsertValue(fan)
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
