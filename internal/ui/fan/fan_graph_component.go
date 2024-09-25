package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/util"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type FanGraphComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout         *tview.Flex
	bmScatterPlot  *tvxwidgets.Plot
	graphComponent *util.GraphComponent[client.Fan]
}

func NewFanGraphComponent(application *tview.Application, fan *client.Fan) *FanGraphComponent {

	graphComponent := util.NewGraphComponent[client.Fan](
		application,
		util.NewGraphComponentConfig().WithReversedOrder(),
		fan,
		func(c *client.Fan) float64 {
			return float64(c.Rpm)
		},
		func(c *client.Fan) float64 {
			return float64(c.Pwm)
		},
	)

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
