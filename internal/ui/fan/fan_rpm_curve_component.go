package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
	"fan2go-tui/internal/ui/util"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type FanRpmCurveComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout         *tview.Flex
	bmScatterPlot  *tvxwidgets.Plot
	graphComponent *util.GraphComponent[client.Fan]
}

func NewFanRpmCurveComponent(application *tview.Application, fan *client.Fan) *FanRpmCurveComponent {
	graphConfig := util.NewGraphComponentConfig().
		WithReversedOrder().
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm)

	fanCurveData := fan.FanCurveData

	currentRpmValue := 0.0

	graphComponent := util.NewGraphComponent[client.Fan](
		application,
		graphConfig,
		fan,
		[]func(*client.Fan) float64{
			func(c *client.Fan) float64 {
				return currentRpmValue
			},
		},
	)

	for _, curveData := range *fanCurveData {
		currentRpmValue = float64(curveData)
		graphComponent.InsertValue(fan)
	}

	c := &FanRpmCurveComponent{
		application:    application,
		graphComponent: graphComponent,
		Fan:            fan,
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
	component := c.graphComponent
	component.InsertValue(fan)
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
