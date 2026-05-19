package graph

import (
	"fan2go-tui/internal/ui/theme"
	uiutil "fan2go-tui/internal/ui/util"
	coreutil "fan2go-tui/internal/util"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/qdm12/reprint"
	"github.com/rivo/tview"
	"golang.org/x/exp/slices"
)

type GraphComponent[T any] struct {
	application *tview.Application

	config *GraphComponentConfig[T]

	Data                *T
	fetchValueFunctions []func(*T) float64

	series []GraphSeries

	yMinValue *float64
	yMaxValue *float64

	layout          *tview.Flex
	plotLayout      *OverlayPlot[T]
	scatterPlotData [][]float64
	valueBufferSize int
}

// GraphSeries captures shared x-axis behavior for lines and bars.
type GraphSeries interface {
	GetXLabel(i int) string
	MapXtoI(x float64) int
	GetXAxisZoomFactor() float64
	SetXAxisZoomFactor(xAxisZoomFactor float64)
	GetXAxisShift() float64
	SetXAxisShift(xAxisShift float64)
	SetXRange(xMin, xMax float64)
	ResetXRange()
	GetXMin() *float64
	GetXMax() *float64
	GetYOffset() float64
	SetYOffset(offset float64)
	GetYAxisZoomFactor() float64
	SetYAxisZoomFactor(yAxisZoomFactor float64)
	GetYAxisShift() float64
	SetYAxisShift(yAxisShift float64)
	SetYRange(yMin, yMax float64)
	ResetYRange()
	GetYMin() *float64
	GetYMax() *float64
}

type GraphDataSource struct {
	Value float64
}

func NewGraphComponent[T any](
	application *tview.Application,
	config *GraphComponentConfig[T],
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

	uiutil.SetupWindow(layout, "")

	plotLayout := NewOverlayPlot[T]()
	c.plotLayout = plotLayout
	if len(c.config.Overlays) > 0 {
		plotLayout.SetOverlays(c.config.Overlays)
	}

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

	plotLayout.SetYAxisLabelDataType(c.config.YAxisLabelDataType)

	layout.AddItem(plotLayout, 0, 1, false)
	_, _, width, _ := plotLayout.GetRect()
	c.setValueBufferSize(width * 4)

	return layout
}

func (c *GraphComponent[T]) SetYMinValue(min *float64) {
	c.yMinValue = min
	if min != nil {
		c.plotLayout.SetYAxisAutoScaleMin(false)
		c.plotLayout.SetMinVal(*min)
	} else {
		c.plotLayout.SetYAxisAutoScaleMin(c.config.YAxisAutoScaleMin)
		c.plotLayout.SetMinVal(0)
	}
}

func (c *GraphComponent[T]) SetYMaxValue(max *float64) {
	c.yMaxValue = max
	if max != nil {
		c.plotLayout.SetYAxisAutoScaleMax(false)
		c.plotLayout.SetMaxVal(*max)
	} else {
		c.plotLayout.SetYAxisAutoScaleMax(c.config.YAxisAutoScaleMax)
	}
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

	c.UpdateValueBufferSize()

	c.updateViewPort()
	lineData := c.computeGraphLineData()
	for idx := range c.fetchValueFunctions {
		c.refreshPlot(idx)
	}
	combinedData := make([][]float64, 0, len(c.scatterPlotData)+len(lineData))
	combinedData = append(combinedData, c.scatterPlotData...)
	combinedData = append(combinedData, lineData...)
	c.plotLayout.SetData(combinedData)
	overlayYMin, overlayYMax := c.computeOverlayPointYRange(combinedData)
	c.plotLayout.SetOverlayContext(OverlayRenderContext[T]{
		XValueToIndex:      c.mapXValueToIndex,
		XValueToIndexFloat: c.mapXValueToIndexFloat,
		YMin:               overlayYMin,
		YMax:               overlayYMax,
		Data:               c.Data,
		Bars:               c.GetBars(),
		ValueBufferSize:    c.GetValueBufferSize(),
		Reversed:           c.config.Reversed,
	})

	// TODO: think about what to do with multiple lines
	for _, line := range c.GetLines() {
		xMax := line.xMax
		if xMax != nil {
			//c.ZoomToRangeX(0, *xMax)
		}
	}
}

func (c *GraphComponent[T]) mapXValueToIndex(x float64) int {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return -1
	}

	series := c.GetSeries()
	if len(series) == 0 {
		return -1
	}
	return series[0].MapXtoI(x)
}

func (c *GraphComponent[T]) mapXValueToIndexFloat(x float64) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return math.NaN()
	}

	series := c.GetSeries()
	if len(series) == 0 {
		return math.NaN()
	}

	first := series[0]
	xAxisZoomFactor := first.GetXAxisZoomFactor()
	if xAxisZoomFactor == 0 {
		return math.NaN()
	}

	return (x - first.GetXAxisShift()) / xAxisZoomFactor
}

func (c *GraphComponent[T]) getXAxisZoomFactor() float64 {
	series := c.GetSeries()
	if len(series) > 0 {
		return series[0].GetXAxisZoomFactor()
	}
	return 1.0
}

func (c *GraphComponent[T]) getXAxisShift() float64 {
	series := c.GetSeries()
	if len(series) > 0 {
		return series[0].GetXAxisShift()
	}
	return 0.0
}

func (c *GraphComponent[T]) getXMin() *float64 {
	series := c.GetSeries()
	if len(series) > 0 {
		return series[0].GetXMin()
	}
	return nil
}

func (c *GraphComponent[T]) getXMax() *float64 {
	series := c.GetSeries()
	if len(series) > 0 {
		return series[0].GetXMax()
	}
	return nil
}

func (c *GraphComponent[T]) updateViewPort() {
	maxYOffset := 0.0
	yAxisZoomFactor := 1.0
	yAxisShift := 0.0
	for _, series := range c.GetSeries() {
		maxYOffset = math.Max(maxYOffset, series.GetYOffset())
		yAxisZoomFactor = math.Max(yAxisZoomFactor, series.GetYAxisZoomFactor())
		yAxisShift = series.GetYAxisShift()
	}

	yMinValue := 0.0
	if c.yMinValue != nil {
		yMinValue = *c.yMinValue
	}

	yMaxValue := 0.0
	if c.yMaxValue != nil {
		yMaxValue = *c.yMaxValue
	}

	viewPortMin := (yMinValue + maxYOffset + yAxisShift) / yAxisZoomFactor
	viewPortMax := (yMaxValue + maxYOffset + yAxisShift) / yAxisZoomFactor
	c.plotLayout.SetYRange(viewPortMin, viewPortMax)
}

func (c *GraphComponent[T]) computeOverlayPointYRange(data [][]float64) (float64, float64) {
	minData := math.Inf(1)
	maxData := math.Inf(-1)
	hasFinite := false

	for _, series := range data {
		for _, value := range series {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				continue
			}

			hasFinite = true
			minData = math.Min(minData, value)
			maxData = math.Max(maxData, value)
		}
	}

	minVal := 0.0
	if c.yMinValue != nil {
		minVal = *c.yMinValue
	} else if c.config.YAxisAutoScaleMin && hasFinite {
		minVal = minData
	}

	maxVal := 0.0
	if c.yMaxValue != nil {
		maxVal = *c.yMaxValue
	} else if c.config.YAxisAutoScaleMax && hasFinite {
		maxVal = maxData
	}

	if maxVal <= minVal {
		if hasFinite && maxData > minVal {
			maxVal = maxData
		} else {
			maxVal = minVal + 1
		}
	}

	return minVal, maxVal
}

func (c *GraphComponent[T]) computeGraphLineData() [][]float64 {
	graphData := make([][]float64, len(c.GetLines()))

	bufferSize := c.GetValueBufferSize()
	for _, line := range c.GetLines() {
		n := bufferSize
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

func (c *GraphComponent[T]) ZoomToRangeX(minX, maxX float64) {
	for _, series := range c.GetSeries() {
		iAtXMin := series.MapXtoI(minX)
		iAtXMax := series.MapXtoI(maxX)

		_, _, width, _ := c.plotLayout.GetRect()
		availableSlots := width - 10

		xScaleFactorToGetXMaxAtEndOfBuffer := float64(1) / (float64(availableSlots) / float64(iAtXMax-iAtXMin))
		newFactor := series.GetXAxisZoomFactor() * xScaleFactorToGetXMaxAtEndOfBuffer
		series.SetXAxisZoomFactor(newFactor)
	}
}

func (c *GraphComponent[T]) refreshPlot(idx int) {
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
	c.Data = data

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
	c.setValueBufferSize(width)
}

func (c *GraphComponent[T]) isVisible() bool {
	return coreutil.IsTxViewVisible(c.layout.Box)
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

func (c *GraphComponent[T]) AddSeries(series GraphSeries) {
	switch series.(type) {
	case *GraphLine, *GraphBar:
		c.series = append(c.series, series)
	default:
		panic("unsupported graph series type")
	}

	c.setXAxisLabelFunc(series)
}

func (c *GraphComponent[T]) setXAxisLabelFunc(series GraphSeries) {
	c.plotLayout.SetXAxisLabelFunc(func(i int) string {
		return series.GetXLabel(i)
	})
}

func (c *GraphComponent[T]) GetSeries() []GraphSeries {
	return c.series
}

func (c *GraphComponent[T]) GetLines() []*GraphLine {
	lines := make([]*GraphLine, 0, len(c.series))
	for _, s := range c.series {
		if line, ok := s.(*GraphLine); ok {
			lines = append(lines, line)
		}
	}
	return lines
}

func (c *GraphComponent[T]) GetBars() []*GraphBar {
	bars := make([]*GraphBar, 0, len(c.series))
	for _, s := range c.series {
		if bar, ok := s.(*GraphBar); ok {
			bars = append(bars, bar)
		}
	}
	return bars
}

func (c *GraphComponent[T]) GetXAxisZoomFactor() float64 {
	return c.getXAxisZoomFactor()
}

func (c *GraphComponent[T]) SetXAxisZoomFactor(xAxisZoomFactor float64) {
	for _, series := range c.GetSeries() {
		series.SetXAxisZoomFactor(xAxisZoomFactor)
	}
	c.Refresh()
}

func (c *GraphComponent[T]) GetXAxisShift() float64 {
	return c.getXAxisShift()
}

func (c *GraphComponent[T]) SetXAxisShift(xAxisShift float64) {
	for _, series := range c.GetSeries() {
		series.SetXAxisShift(xAxisShift)
	}
	c.Refresh()
}

func (c *GraphComponent[T]) GetYAxisZoomFactor() float64 {
	series := c.GetSeries()
	if len(series) == 0 {
		return 1.0
	}
	return series[0].GetYAxisZoomFactor()
}

func (c *GraphComponent[T]) SetYAxisZoomFactor(yAxisZoomFactor float64) {
	for _, series := range c.GetSeries() {
		series.SetYAxisZoomFactor(yAxisZoomFactor)
	}
	c.Refresh()
}

func (c *GraphComponent[T]) GetYAxisShift() float64 {
	series := c.GetSeries()
	if len(series) == 0 {
		return 0.0
	}
	return series[0].GetYAxisShift()
}

func (c *GraphComponent[T]) SetYAxisShift(yAxisShift float64) {
	for _, series := range c.GetSeries() {
		series.SetYAxisShift(yAxisShift)
	}
	c.Refresh()
}

func (c *GraphComponent[T]) GetPlotRect() (int, int, int, int) {
	return c.plotLayout.GetPlotRect()
}

func (c *GraphComponent[T]) SetXRange(xMin, xMax float64) {
	for _, series := range c.GetSeries() {
		series.SetXRange(xMin, xMax)
	}
}

func (c *GraphComponent[T]) ResetXRange() {
	for _, series := range c.GetSeries() {
		series.ResetXRange()
	}
}

func (c *GraphComponent[T]) GetXMin() *float64 {
	return c.getXMin()
}

func (c *GraphComponent[T]) GetXMax() *float64 {
	return c.getXMax()
}
