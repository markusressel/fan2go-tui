package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
	"fan2go-tui/internal/ui/util"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"slices"
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
		WithPlotColors(theme.Colors.Graph.Rpm, theme.Colors.Graph.Pwm).
		WithXMax(256)

	graphComponent := util.NewGraphComponent[client.Fan](
		application,
		graphConfig,
		fan,
		[]func(*client.Fan) float64{
			func(c *client.Fan) float64 {
				return 0
			},
		},
	)

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

	fanRpmGraphData, err := computeFanRpmGraphData(c.Fan)
	if err != nil {
		return
	}
	c.graphComponent.SetRawData(fanRpmGraphData)
}

func computeFanRpmGraphData(fan *client.Fan) ([][]float64, error) {
	fanCurveData := *fan.FanCurveData

	graphData := make([][]float64, 1)

	rpmValues := []float64{}

	pwmKeys := maps.Keys(fanCurveData)
	slices.Sort(pwmKeys)

	for _, pwmKey := range pwmKeys {
		//pwmValue := float64(pwm)
		currentRpmValue := fanCurveData[pwmKey]
		rpmValues = append(rpmValues, currentRpmValue)
	}

	graphData[0] = rpmValues

	return graphData, nil
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
