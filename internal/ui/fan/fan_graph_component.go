package fan

import (
	"fan2go-tui/internal/ui/theme"
	uiutil "fan2go-tui/internal/ui/util"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

type FanGraphComponent[T any] struct {
	application *tview.Application

	Data       *T
	fetchValue func(*T) float64

	layout          *tview.Flex
	bmScatterPlot   *tvxwidgets.Plot
	scatterPlotData [][]float64
}

func NewFanGraphComponent[T any](application *tview.Application, data *T, fetchValue func(*T) float64) *FanGraphComponent[T] {
	c := &FanGraphComponent[T]{
		application:     application,
		Data:            data,
		fetchValue:      fetchValue,
		scatterPlotData: make([][]float64, 2),
	}

	c.layout = c.createLayout()

	return c
}

func (c *FanGraphComponent[T]) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	layout.SetBorder(true)
	uiutil.SetupWindow(layout, "")

	bmScatterPlot := tvxwidgets.NewPlot()
	c.bmScatterPlot = bmScatterPlot
	bmScatterPlot.SetLineColor([]tcell.Color{
		tcell.ColorGold,
		tcell.ColorLightSkyBlue,
	})
	bmScatterPlot.SetPlotType(tvxwidgets.PlotTypeScatter)
	bmScatterPlot.SetMarker(tvxwidgets.PlotMarkerBraille)
	layout.AddItem(bmScatterPlot, 0, 1, false)

	return layout
}

func (c *FanGraphComponent[T]) Refresh() {
	c.bmScatterPlot.SetData(c.scatterPlotData)
}

func (c *FanGraphComponent[T]) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanGraphComponent[T]) SetTitle(title string) {
	titleText := theme.CreateTitleText(title)
	c.layout.SetTitle(titleText)
}

func (c *FanGraphComponent[T]) InsertValue(data *T) {
	value := c.fetchValue(data)
	c.scatterPlotData[0] = append(c.scatterPlotData[0], value)
	// limit data to 100 points
	if (len(c.scatterPlotData[0])) > 100 {
		c.scatterPlotData[0] = c.scatterPlotData[0][1:]
	}
	c.Refresh()
}
