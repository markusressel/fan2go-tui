package util

import (
	"fan2go-tui/internal/ui/theme"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
)

type GraphComponent[T any] struct {
	application *tview.Application

	Data        *T
	fetchValue  func(*T) float64
	fetchValue2 func(*T) float64

	layout          *tview.Flex
	bmScatterPlot   *tvxwidgets.Plot
	scatterPlotData [][]float64
	valueBufferSize int
}

func NewGraphComponent[T any](application *tview.Application, data *T, fetchValue func(*T) float64, fetchValue2 func(*T) float64) *GraphComponent[T] {
	c := &GraphComponent[T]{
		application:     application,
		Data:            data,
		fetchValue:      fetchValue,
		fetchValue2:     fetchValue2,
		scatterPlotData: make([][]float64, 2),
	}

	c.layout = c.createLayout()

	return c
}

func (c *GraphComponent[T]) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	SetupWindow(layout, "")

	bmScatterPlot := tvxwidgets.NewPlot()
	c.bmScatterPlot = bmScatterPlot
	bmScatterPlot.SetLineColor([]tcell.Color{
		theme.Colors.Graphs.Rpm,
		theme.Colors.Graphs.Pwm,
	})
	bmScatterPlot.SetPlotType(tvxwidgets.PlotTypeLineChart)
	bmScatterPlot.SetMarker(tvxwidgets.PlotMarkerBraille)
	layout.AddItem(bmScatterPlot, 0, 1, false)
	_, _, width, _ := bmScatterPlot.GetRect()
	c.valueBufferSize = width * 4

	return layout
}

func (c *GraphComponent[T]) Refresh() {
	c.bmScatterPlot.SetDrawAxes(true)
	c.bmScatterPlot.SetData(c.scatterPlotData)

	_, _, width, _ := c.bmScatterPlot.GetRect()
	c.valueBufferSize = width - 5
}

func (c *GraphComponent[T]) GetLayout() *tview.Flex {
	return c.layout
}

func (c *GraphComponent[T]) SetTitle(title string) {
	titleText := theme.CreateTitleText(title)
	c.layout.SetTitle(titleText)
}

func (c *GraphComponent[T]) InsertValue(data *T) {
	if c.fetchValue != nil {
		value := c.fetchValue(data)
		c.scatterPlotData[0] = slices.Insert(c.scatterPlotData[0], len(c.scatterPlotData[0]), value)
		// limit data to visible data points
		if (len(c.scatterPlotData[0])) > c.valueBufferSize {
			overflow := len(c.scatterPlotData[0]) - c.valueBufferSize
			c.scatterPlotData[0] = c.scatterPlotData[0][overflow:]
		}
	}

	if c.fetchValue2 != nil {
		value2 := c.fetchValue2(data)
		c.scatterPlotData[1] = slices.Insert(c.scatterPlotData[1], len(c.scatterPlotData[1]), value2)
		// limit data to visible data points
		if (len(c.scatterPlotData[1])) > c.valueBufferSize {
			overflow := len(c.scatterPlotData[1]) - c.valueBufferSize
			c.scatterPlotData[1] = c.scatterPlotData[1][overflow:]
		}
	}

	c.Refresh()
}
