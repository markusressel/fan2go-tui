package util

import (
	"fan2go-tui/internal/ui/theme"
	"fan2go-tui/internal/util"
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/qdm12/reprint"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
	"math"
)

type GraphComponent[T any] struct {
	application *tview.Application

	config *GraphComponentConfig

	Data                *T
	fetchValueFunctions []func(*T) float64

	graphLines []*GraphLine

	yMinValue *float64
	yMaxValue *float64

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
		// Ensure that the number of plot colors matches the number of fetch value functions
		plotColors := reprint.This(c.config.PlotColors).([]tcell.Color)
		for i := len(plotColors); i < len(c.fetchValueFunctions); i++ {
			plotColors = append(plotColors, theme.Colors.Graph.Default)
		}

		// add default color a couple of times to make sure we have enough colors
		for i := len(plotColors); i < 5; i++ {
			plotColors = append(plotColors, theme.Colors.Graph.Default)
		}

		plotLayout.SetLineColor(plotColors)
	}
	plotLayout.SetPlotType(c.config.PlotType)
	plotLayout.SetMarker(c.config.MarkerType)

	plotLayout.SetDrawXAxisLabel(c.config.DrawXAxisLabel)
	plotLayout.SetDrawYAxisLabel(c.config.DrawYAxisLabel)
	plotLayout.SetYAxisAutoScaleMin(c.config.YAxisAutoScaleMin)
	plotLayout.SetYAxisAutoScaleMax(c.config.YAxisAutoScaleMax)
	plotLayout.SetXAxisLabelFunc(c.config.XAxisLabelFunc)

	layout.AddItem(plotLayout, 0, 1, false)
	_, _, width, _ := plotLayout.GetRect()
	c.setValueBufferSize(width * 4)

	return layout
}

func (c *GraphComponent[T]) SetYMinValue(min float64) {
	c.yMinValue = &min
	c.plotLayout.SetMinVal(min)
}

func (c *GraphComponent[T]) SetYMaxValue(max float64) {
	c.yMaxValue = &max
	c.plotLayout.SetMaxVal(max)
}

func (c *GraphComponent[T]) SetYRange(min, max float64) {
	c.yMinValue = &min
	c.yMaxValue = &max
	c.plotLayout.SetYRange(min, max)
}

func (c *GraphComponent[T]) Refresh() {
	c.plotLayout.SetDrawAxes(true)
	if c.yMinValue != nil {
		c.plotLayout.SetMinVal(*c.yMinValue)
	}
	if c.yMaxValue != nil {
		c.plotLayout.SetMaxVal(*c.yMaxValue)
	}

	c.updateCamera()
	lineData := c.computeGraphLineData()
	combinedData := make([][]float64, 0, len(c.scatterPlotData)+len(lineData))
	combinedData = append(combinedData, c.scatterPlotData...)
	combinedData = append(combinedData, lineData...)
	c.plotLayout.SetData(combinedData)

	c.UpdateValueBufferSize()

	for idx := range c.fetchValueFunctions {
		c.refreshPlot(idx)
	}
}

func (c *GraphComponent[T]) updateCamera() {
	maxYOffset := 0.0
	yAxisZoomFactor := 1.0
	yAxisShift := 0.0
	for _, line := range c.graphLines {
		maxYOffset = math.Max(maxYOffset, line.GetYOffset())
		yAxisZoomFactor = math.Max(yAxisZoomFactor, line.GetYAxisZoomFactor())
		yAxisShift = math.Max(yAxisShift, line.GetYAxisShift())
	}

	c.plotLayout.SetYRange(
		(-1+maxYOffset+yAxisShift)/yAxisZoomFactor,
		(1+maxYOffset+yAxisShift)/yAxisZoomFactor,
	)
}

func (c *GraphComponent[T]) computeGraphLineData() [][]float64 {
	graphData := make([][]float64, len(c.GetLines()))

	for _, line := range c.GetLines() {
		n := 200
		data := make([]float64, n)
		for i := 0; i < n; i++ {
			xVal := line.GetX(i)
			yVal := line.GetY(xVal)
			data[i] = yVal
		}
		graphData = append(graphData, data)
	}

	return graphData
}

func (c *GraphComponent[T]) refreshPlot(idx int) {
	missingDataPoints := c.valueBufferSize - len(c.scatterPlotData[idx])

	lastDataPoint := math.NaN()
	hasDataPoints := len(c.scatterPlotData[idx]) > 0
	if hasDataPoints {
		if c.config.Reversed {
			lastDataPoint = c.scatterPlotData[idx][0]
		} else {
			lastDataPoint = c.scatterPlotData[idx][len(c.scatterPlotData[idx])-1]
		}
	}

	for i := 0; i < missingDataPoints; i++ {
		targetIndex := 0
		if c.config.Reversed {
			targetIndex = len(c.scatterPlotData[idx])
		}
		c.scatterPlotData[idx] = slices.Insert(c.scatterPlotData[idx], targetIndex, lastDataPoint)
	}

	// limit data to visible data points
	overflow := len(c.scatterPlotData[idx]) - c.valueBufferSize
	if c.config.Reversed {
		c.scatterPlotData[idx] = c.scatterPlotData[idx][:len(c.scatterPlotData[idx])-overflow]
	} else {
		c.scatterPlotData[idx] = c.scatterPlotData[idx][overflow:]
	}
}

func (c *GraphComponent[T]) GetLayout() *tview.Flex {
	return c.layout
}

func (c *GraphComponent[T]) SetTitle(title string) {
	titleText := theme.CreateTitleText(title)
	c.layout.SetTitle(titleText)
}

// SetRawData sets the raw data for the graph component.
func (c *GraphComponent[T]) SetRawData(data [][]float64) {
	c.scatterPlotData = data
	c.Refresh()
}

func (c *GraphComponent[T]) InsertValue(data *T) {
	for idx, fetchValue := range c.fetchValueFunctions {
		value := fetchValue(data)
		plotDataPoints := c.scatterPlotData[idx]
		targetIndex := len(plotDataPoints)

		if c.config.Reversed {
			reversedCopy := slices.Clone(plotDataPoints)
			slices.Reverse(reversedCopy)
			reversedCopy = slices.Insert(reversedCopy, targetIndex, value)
			slices.Reverse(reversedCopy)
			plotDataPoints = reversedCopy
		} else {
			plotDataPoints = slices.Insert(plotDataPoints, targetIndex, value)
		}
		c.scatterPlotData[idx] = plotDataPoints
	}

	c.Refresh()
}

func (c *GraphComponent[T]) UpdateValueBufferSize() {
	if !c.isVisible() {
		c.setValueBufferSize(500)
		return
	}

	_, _, width, _ := c.plotLayout.GetInnerRect()
	c.setValueBufferSize(width - 5)
}

func (c *GraphComponent[T]) isVisible() bool {
	return util.IsTxViewVisible(c.layout.Box)
}

func (c *GraphComponent[T]) setValueBufferSize(i int) {
	if c.config.XMax > 0 {
		if i > c.config.XMax {
			i = c.config.XMax
		}
	}
	if i < 1 {
		i = 1
	}
	c.valueBufferSize = i
}

func (c *GraphComponent[T]) GetValueBufferSize() int {
	return c.valueBufferSize
}

func (c *GraphComponent[T]) AddLine(graphLineConfig *GraphLine) *GraphLine {
	c.graphLines = append(c.graphLines, graphLineConfig)
	return graphLineConfig
}

func (c *GraphComponent[T]) GetLines() []*GraphLine {
	return c.graphLines
}

func (c *GraphComponent[T]) SetXAxisZoomFactor(xAxisZoomFactor float64) {
	for _, line := range c.graphLines {
		line.SetXAxisZoomFactor(xAxisZoomFactor)
	}
	c.Refresh()
}

func (c *GraphComponent[T]) SetXAxisShift(xAxisShift float64) {
	for _, line := range c.graphLines {
		line.SetXAxisShift(xAxisShift)
	}
	c.Refresh()
}

func (c *GraphComponent[T]) SetYAxisZoomFactor(yAxisZoomFactor float64) {
	for _, line := range c.graphLines {
		line.SetYAxisZoomFactor(yAxisZoomFactor)
	}
	c.Refresh()
}

func (c *GraphComponent[T]) SetYAxisShift(yAxisShift float64) {
	for _, line := range c.graphLines {
		line.SetYAxisShift(yAxisShift)
	}
	c.Refresh()
}
