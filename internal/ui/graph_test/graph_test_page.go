package fan

import (
	"fan2go-tui/internal/ui/util"
	"math"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GraphTestPage struct {
	application *tview.Application

	layout *tview.Flex

	graphComponent *util.GraphComponent[util.GraphDataSource]
}

func NewGraphTestPage(application *tview.Application) GraphTestPage {

	graphTestPage := GraphTestPage{
		application: application,
	}

	graphTestPage.layout = graphTestPage.createLayout()

	return graphTestPage
}

func (c *GraphTestPage) createLayout() *tview.Flex {
	graphTestPageLayout := tview.NewFlex()
	c.layout = graphTestPageLayout

	graphDataSource := &util.GraphDataSource{
		Value: 0.0,
	}
	graphComponent := util.NewGraphComponent[util.GraphDataSource](
		c.application,
		util.NewGraphComponentConfig().
			WithYAxisAutoScaleMin(false).
			WithYAxisAutoScaleMax(false),
		graphDataSource,
		[]func(val *util.GraphDataSource) float64{
			func(val *util.GraphDataSource) float64 {
				return val.Value
			},
		},
	)
	c.graphComponent = graphComponent

	data := make([][]float64, 1)
	for i := 0; i < 100; i++ {
		data[0] = append(data[0], float64(i))
	}

	xFunc := func(i int) float64 {
		return float64(i)
	}
	fFunc := func(x float64) float64 {
		if x < 0 {
			return math.NaN()
		} else if x > 100 {
			return math.NaN()
		} else {
			return x
		}
	}
	xLabelFunc := func(i int, x float64) string {
		return strconv.Itoa(int(x))
	}

	graphComponent.AddLine(util.NewGraphLine(
		"Test",
		xFunc,
		fFunc,
		xLabelFunc,
	))

	graphComponent.SetYRange(-100, 200)
	graphComponent.SetXRange(0, 100)

	graphTestPageLayout.AddItem(c.graphComponent.GetLayout(), 0, 1, true)

	c.layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()

		// Shift Controls
		if event.Modifiers() == tcell.ModShift && key == tcell.KeyUp {
			c.graphComponent.SetYAxisShift(c.graphComponent.GetYAxisShift() + 1)
		} else if event.Modifiers() == tcell.ModShift && key == tcell.KeyDown {
			c.graphComponent.SetYAxisShift(c.graphComponent.GetYAxisShift() - 1)
		} else if event.Modifiers() == tcell.ModShift && key == tcell.KeyRight {
			c.graphComponent.SetXAxisShift(c.graphComponent.GetXAxisShift() + 1)
		} else if event.Modifiers() == tcell.ModShift && key == tcell.KeyLeft {
			c.graphComponent.SetXAxisShift(c.graphComponent.GetXAxisShift() - 1)
			// Zoom Controls
		} else if event.Modifiers() == tcell.ModCtrl && key == tcell.KeyUp {
			c.graphComponent.SetYAxisZoomFactor(c.graphComponent.GetYAxisZoomFactor() + 0.1)
		} else if event.Modifiers() == tcell.ModCtrl && key == tcell.KeyDown {
			c.graphComponent.SetYAxisZoomFactor(c.graphComponent.GetYAxisZoomFactor() - 0.1)
		} else if event.Modifiers() == tcell.ModCtrl && key == tcell.KeyRight {
			c.graphComponent.SetXAxisZoomFactor(c.graphComponent.GetXAxisZoomFactor() + 0.1)
		} else if event.Modifiers() == tcell.ModCtrl && key == tcell.KeyLeft {
			c.graphComponent.SetXAxisZoomFactor(c.graphComponent.GetXAxisZoomFactor() - 0.1)
			// Reset Controls to default
		} else if key == tcell.KeyCtrlR {
			c.graphComponent.SetYAxisShift(0)
			c.graphComponent.SetXAxisShift(0)
		}
		return event
	})

	return c.layout
}

func (c *GraphTestPage) GetLayout() *tview.Flex {
	return c.layout
}

func (c *GraphTestPage) Refresh() error {
	c.graphComponent.Refresh()
	c.graphComponent.ZoomToRangeX(0, 100)
	return nil
}
