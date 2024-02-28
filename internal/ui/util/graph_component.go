package util

import (
	"fan2go-tui/internal/ui/theme"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
	"math"
)

type GraphComponent[T any] struct {
	application *tview.Application

	Data        *T
	fetchValue  func(*T) float64
	fetchValue2 func(*T) float64

	layout          *tview.Flex
	plotLayout      *tvxwidgets.Plot
	scatterPlotData [][]float64
	valueBufferSize int

	reversed bool
}

func NewGraphComponent[T any](
	application *tview.Application,
	data *T,
	fetchValue func(*T) float64,
	fetchValue2 func(*T) float64,
	reversed bool,
) *GraphComponent[T] {
	c := &GraphComponent[T]{
		application:     application,
		Data:            data,
		fetchValue:      fetchValue,
		fetchValue2:     fetchValue2,
		scatterPlotData: make([][]float64, 2),
		reversed:        reversed,
	}

	c.layout = c.createLayout()

	return c
}

func (c *GraphComponent[T]) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	SetupWindow(layout, "")

	plotLayout := tvxwidgets.NewPlot()
	c.plotLayout = plotLayout
	plotLayout.SetLineColor([]tcell.Color{
		theme.Colors.Graph.Rpm,
		theme.Colors.Graph.Pwm,
	})
	plotLayout.SetPlotType(tvxwidgets.PlotTypeLineChart)
	plotLayout.SetMarker(tvxwidgets.PlotMarkerBraille)
	layout.AddItem(plotLayout, 0, 1, false)
	_, _, width, _ := plotLayout.GetRect()
	c.valueBufferSize = width * 4

	return layout
}

func (c *GraphComponent[T]) Refresh() {
	c.plotLayout.SetDrawAxes(true)
	c.plotLayout.SetData(c.scatterPlotData)

	_, _, width, _ := c.plotLayout.GetRect()
	c.valueBufferSize = width - 5

	if c.fetchValue != nil {
		missingDataPoints := c.valueBufferSize - len(c.scatterPlotData[0])
		for i := 0; i < missingDataPoints; i++ {
			targetIndex := 0
			if c.reversed {
				targetIndex = len(c.scatterPlotData[0])
			}
			c.scatterPlotData[0] = slices.Insert(c.scatterPlotData[0], targetIndex, math.NaN())
		}

		// limit data to visible data points
		overflow := len(c.scatterPlotData[0]) - c.valueBufferSize
		if c.reversed {
			c.scatterPlotData[0] = c.scatterPlotData[0][:len(c.scatterPlotData[0])-overflow]
		} else {
			c.scatterPlotData[0] = c.scatterPlotData[0][overflow:]
		}
	}

	if c.fetchValue2 != nil {
		missingDataPoints := c.valueBufferSize - len(c.scatterPlotData[1])
		for i := 0; i < missingDataPoints; i++ {
			targetIndex := 0
			if c.reversed {
				targetIndex = len(c.scatterPlotData[1])
			}
			c.scatterPlotData[1] = slices.Insert(c.scatterPlotData[1], targetIndex, math.NaN())
		}

		// limit data to visible data points
		overflow := len(c.scatterPlotData[1]) - c.valueBufferSize
		if c.reversed {
			c.scatterPlotData[1] = c.scatterPlotData[1][:len(c.scatterPlotData[1])-overflow]
		} else {
			c.scatterPlotData[1] = c.scatterPlotData[1][overflow:]
		}
	}
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
		data1 := c.scatterPlotData[0]
		targetIndex := len(data1)
		if c.reversed {
			reversedCopy := slices.Clone(data1)
			slices.Reverse(reversedCopy)
			reversedCopy = slices.Insert(reversedCopy, targetIndex, value)
			slices.Reverse(reversedCopy)
			data1 = reversedCopy
		} else {
			data1 = slices.Insert(data1, targetIndex, value)
		}
		c.scatterPlotData[0] = data1
	}
	if c.fetchValue2 != nil {
		value2 := c.fetchValue2(data)
		data2 := c.scatterPlotData[1]
		targetIndex := len(data2)
		if c.reversed {
			reversedCopy := slices.Clone(data2)
			slices.Reverse(reversedCopy)
			reversedCopy = slices.Insert(reversedCopy, targetIndex, value2)
			slices.Reverse(reversedCopy)
			data2 = reversedCopy
		} else {
			data2 = slices.Insert(data2, targetIndex, value2)
		}
		c.scatterPlotData[1] = data2
	}

	c.Refresh()
}
