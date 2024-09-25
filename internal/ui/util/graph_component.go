package util

import (
	"fan2go-tui/internal/ui/theme"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
	"math"
)

type GraphComponent[T any] struct {
	application *tview.Application

	config *GraphComponentConfig

	Data                *T
	fetchValueFunctions []func(*T) float64

	layout          *tview.Flex
	plotLayout      *tvxwidgets.Plot
	scatterPlotData [][]float64
	valueBufferSize int
}

func NewGraphComponent[T any](
	application *tview.Application,
	config *GraphComponentConfig,
	data *T,
	fetchValueFunctions []func(*T) float64,
) *GraphComponent[T] {
	c := &GraphComponent[T]{
		application:         application,
		config:              config,
		Data:                data,
		fetchValueFunctions: fetchValueFunctions,
		scatterPlotData:     make([][]float64, len(fetchValueFunctions)),
	}

	c.layout = c.createLayout()

	return c
}

func (c *GraphComponent[T]) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	SetupWindow(layout, "")

	plotLayout := tvxwidgets.NewPlot()
	c.plotLayout = plotLayout

	if len(c.config.PlotColors) > 0 {
		plotLayout.SetLineColor(c.config.PlotColors)
	}
	plotLayout.SetPlotType(c.config.PlotType)
	plotLayout.SetMarker(c.config.MarkerType)

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

	for idx := range c.fetchValueFunctions {
		missingDataPoints := c.valueBufferSize - len(c.scatterPlotData[idx])
		for i := 0; i < missingDataPoints; i++ {
			targetIndex := 0
			if c.config.Reversed {
				targetIndex = len(c.scatterPlotData[idx])
			}
			c.scatterPlotData[idx] = slices.Insert(c.scatterPlotData[idx], targetIndex, math.NaN())
		}

		// limit data to visible data points
		overflow := len(c.scatterPlotData[idx]) - c.valueBufferSize
		if c.config.Reversed {
			c.scatterPlotData[idx] = c.scatterPlotData[idx][:len(c.scatterPlotData[idx])-overflow]
		} else {
			c.scatterPlotData[idx] = c.scatterPlotData[idx][overflow:]
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
	for idx, fetchValue := range c.fetchValueFunctions {
		value := fetchValue(data)
		data := c.scatterPlotData[idx]
		targetIndex := len(data)
		if c.config.Reversed {
			reversedCopy := slices.Clone(data)
			slices.Reverse(reversedCopy)
			reversedCopy = slices.Insert(reversedCopy, targetIndex, value)
			slices.Reverse(reversedCopy)
			data = reversedCopy
		} else {
			data = slices.Insert(data, targetIndex, value)
		}
		c.scatterPlotData[idx] = data
	}

	c.Refresh()
}
