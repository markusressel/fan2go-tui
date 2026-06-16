package graph

import (
	"fan2go-tui/internal/ui/theme"
	uiutil "fan2go-tui/internal/ui/util"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/qdm12/reprint"
	"github.com/rivo/tview"
)

type GraphComponent struct {
	application *tview.Application

	config *GraphComponentConfig

	series []GraphSeries
	legend []*GraphSeriesLegend

	yMinValue *float64
	yMaxValue *float64

	zoomMinX float64
	zoomMaxX float64

	layout          *tview.Flex
	plotLayout      *OverlayPlot
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

type GraphSeriesOption func(*graphSeriesOptions)

type graphSeriesOptions struct {
	legend *GraphSeriesLegend
}

func WithLegend(legend *GraphSeriesLegend) GraphSeriesOption {
	return func(opts *graphSeriesOptions) {
		opts.legend = legend
	}
}

func NewGraphComponent(
	application *tview.Application,
	config *GraphComponentConfig,
) *GraphComponent {
	c := &GraphComponent{
		application: application,
		config:      config,
		zoomMinX:    math.NaN(),
		zoomMaxX:    math.NaN(),
	}

	c.layout = c.createLayout()

	return c
}

func (c *GraphComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	uiutil.SetupWindow(layout, "")

	plotLayout := NewOverlayPlot()
	c.plotLayout = plotLayout
	if len(c.config.Overlays) > 0 {
		plotLayout.SetOverlays(c.config.Overlays)
	}

	if len(c.config.PlotColors) > 0 {
		// Ensure we have enough default colors for multi-series rendering.
		plotColors := reprint.This(c.config.PlotColors).([]tcell.Color)

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

	uiutil.SetupReactiveResize(c.application, plotLayout.Box, c.Refresh)

	layout.AddItem(plotLayout, 0, 1, false)
	_, _, width, _ := plotLayout.GetRect()
	c.setValueBufferSize(width * 4)

	return layout
}

func (c *GraphComponent) SetYMinValue(min *float64) {
	c.yMinValue = min
	if min != nil {
		c.plotLayout.SetYAxisAutoScaleMin(false)
		c.plotLayout.SetMinVal(*min)
	} else {
		c.plotLayout.SetYAxisAutoScaleMin(c.config.YAxisAutoScaleMin)
		c.plotLayout.SetMinVal(0)
	}
}

func (c *GraphComponent) SetYMaxValue(max *float64) {
	c.yMaxValue = max
	if max != nil {
		c.plotLayout.SetYAxisAutoScaleMax(false)
		c.plotLayout.SetMaxVal(*max)
	} else {
		c.plotLayout.SetYAxisAutoScaleMax(c.config.YAxisAutoScaleMax)
	}
}

func (c *GraphComponent) SetYRange(min, max float64) {
	c.yMinValue = &min
	c.yMaxValue = &max
	c.plotLayout.SetYRange(min, max)
}

func (c *GraphComponent) Refresh() {
	c.plotLayout.SetDrawAxes(true)
	if c.yMinValue != nil {
		c.plotLayout.SetMinVal(*c.yMinValue)
	}
	if c.yMaxValue != nil {
		c.plotLayout.SetMaxVal(*c.yMaxValue)
	}

	c.UpdateValueBufferSize()
	c.applyZoom()

	c.updateViewPort()
	combinedData := c.computePlotSeriesData()
	placeholderData := c.createPlaceholderSeriesData()
	c.plotLayout.SetData(placeholderData)
	c.applyAutoScaleFromData(combinedData)
	overlayYMin, overlayYMax := c.computeOverlayPointYRange(combinedData)
	lineColors := c.getPlotColors(len(combinedData))
	legendOverlay := c.buildLegendOverlay(c.getPlotColors(len(c.series)))
	c.plotLayout.SetOverlayContext(OverlayRenderContext{
		XValueToIndex:      c.mapXValueToIndex,
		XValueToIndexFloat: c.mapXValueToIndexFloat,
		YMin:               overlayYMin,
		YMax:               overlayYMax,
		Bars:               c.GetBars(),
		ValueBufferSize:    c.GetValueBufferSize(),
		Reversed:           c.config.Reversed,
		SeriesData:         combinedData,
		SeriesColors:       lineColors,
		YAxisLabelsAreInts: c.config.YAxisLabelDataType == tvxwidgets.PlotYAxisLabelDataInt,
		LegendOverlay:      legendOverlay,
	})

	// TODO: think about what to do with multiple lines
	for _, line := range c.GetLines() {
		xMax := line.xMax
		if xMax != nil {
			//c.ZoomToRangeX(0, *xMax)
		}
	}
}

func (c *GraphComponent) createPlaceholderSeriesData() [][]float64 {
	bufferSize := c.GetValueBufferSize()
	placeholder := make([]float64, bufferSize)
	for i := 0; i < bufferSize; i++ {
		placeholder[i] = math.NaN()
	}
	return [][]float64{placeholder}
}

func (c *GraphComponent) applyAutoScaleFromData(data [][]float64) {
	if len(data) == 0 {
		return
	}

	if c.yMinValue == nil && c.config.YAxisAutoScaleMin {
		min := math.Inf(1)
		has := false
		for _, s := range data {
			for _, v := range s {
				if math.IsNaN(v) || math.IsInf(v, 0) {
					continue
				}
				has = true
				min = math.Min(min, v)
			}
		}
		if has {
			c.plotLayout.SetMinVal(min)
		}
	}

	if c.yMaxValue == nil && c.config.YAxisAutoScaleMax {
		max := math.Inf(-1)
		has := false
		for _, s := range data {
			for _, v := range s {
				if math.IsNaN(v) || math.IsInf(v, 0) {
					continue
				}
				has = true
				max = math.Max(max, v)
			}
		}
		if has {
			c.plotLayout.SetMaxVal(max)
		}
	}
}

func (c *GraphComponent) getPlotColors(required int) []tcell.Color {
	colors := reprint.This(c.config.PlotColors).([]tcell.Color)
	for i := len(colors); i < required; i++ {
		colors = append(colors, theme.Colors.Graph.Default)
	}
	return colors
}

func (c *GraphComponent) mapXValueToIndex(x float64) int {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return -1
	}

	series := c.GetSeries()
	if len(series) == 0 {
		return -1
	}
	return series[0].MapXtoI(x)
}

func (c *GraphComponent) mapXValueToIndexFloat(x float64) float64 {
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

func (c *GraphComponent) getXAxisZoomFactor() float64 {
	series := c.GetSeries()
	if len(series) > 0 {
		return series[0].GetXAxisZoomFactor()
	}
	return 1.0
}

func (c *GraphComponent) getXAxisShift() float64 {
	series := c.GetSeries()
	if len(series) > 0 {
		return series[0].GetXAxisShift()
	}
	return 0.0
}

func (c *GraphComponent) getXMin() *float64 {
	series := c.GetSeries()
	if len(series) > 0 {
		return series[0].GetXMin()
	}
	return nil
}

func (c *GraphComponent) getXMax() *float64 {
	series := c.GetSeries()
	if len(series) > 0 {
		return series[0].GetXMax()
	}
	return nil
}

func (c *GraphComponent) updateViewPort() {
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

func (c *GraphComponent) computeOverlayPointYRange(data [][]float64) (float64, float64) {
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

func (c *GraphComponent) computePlotSeriesData() [][]float64 {
	lines := c.GetLines()
	lineSeriesData := make([][]float64, 0, len(lines))

	bufferSize := c.GetValueBufferSize()
	for _, line := range lines {
		data := make([]float64, bufferSize)
		for i := 0; i < bufferSize; i++ {
			xVal := line.GetX(i)
			yVal := line.GetY(xVal)
			data[i] = yVal
		}

		lineSeriesData = append(lineSeriesData, data)
	}

	return lineSeriesData
}

func (c *GraphComponent) ZoomToRangeX(minX, maxX float64) {
	c.zoomMinX = minX
	c.zoomMaxX = maxX
	c.applyZoom()
}

func (c *GraphComponent) applyZoom() {
	if math.IsNaN(c.zoomMinX) || math.IsNaN(c.zoomMaxX) {
		return
	}

	minX := c.zoomMinX
	maxX := c.zoomMaxX

	span := maxX - minX
	if span <= 0 {
		return
	}

	_, _, width, _ := c.plotLayout.GetInnerRect()
	availableSlots := width - 10
	if availableSlots <= 0 {
		return
	}

	newFactor := span / float64(availableSlots)
	if math.IsNaN(newFactor) || math.IsInf(newFactor, 0) || newFactor <= 0 {
		return
	}

	for _, series := range c.GetSeries() {
		series.SetXAxisZoomFactor(newFactor)
	}
}

func (c *GraphComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *GraphComponent) SetTitle(title string) {
	titleText := theme.CreateTitleText(title)
	c.layout.SetTitle(titleText)
}

func (c *GraphComponent) UpdateValueBufferSize() {
	if !c.isVisible() {
		return
	}

	_, _, width, _ := c.plotLayout.GetInnerRect()
	if width <= 0 {
		return
	}
	c.setValueBufferSize(width)
}

func (c *GraphComponent) isVisible() bool {
	_, _, width, height := c.layout.GetRect()
	return width > 0 && height > 0
}

func (c *GraphComponent) setValueBufferSize(i int) {
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

func (c *GraphComponent) GetValueBufferSize() int {
	return c.valueBufferSize
}

func (c *GraphComponent) AddSeries(series GraphSeries, options ...GraphSeriesOption) {
	opts := graphSeriesOptions{}
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	switch series.(type) {
	case *GraphLine, *GraphBar:
		c.series = append(c.series, series)
		c.legend = append(c.legend, opts.legend)
	default:
		panic("unsupported graph series type")
	}

	c.setXAxisLabelFunc(series)
}

func (c *GraphComponent) AddSeriesWithLegend(series GraphSeries, legend *GraphSeriesLegend) {
	c.AddSeries(series, WithLegend(legend))
}

func (c *GraphComponent) buildLegendOverlay(seriesColors []tcell.Color) *LegendOverlay {
	if len(c.legend) == 0 {
		return nil
	}

	entries := make([]LegendOverlayEntry, 0, len(c.legend))
	for i, legend := range c.legend {
		if legend == nil {
			continue
		}

		text := legend.displayText()
		if text == "" {
			continue
		}

		bgColor := tcell.ColorWhite
		if i < len(seriesColors) {
			bgColor = seriesColors[i]
		}
		if legend.backgroundColor != nil {
			bgColor = *legend.backgroundColor
		}

		glyph := rune(0x25CF) // ●
		if legend.glyph != nil {
			glyph = *legend.glyph
		}

		textColor := tcell.ColorDefault
		autoTextColor := true
		if legend.textColor != nil {
			textColor = *legend.textColor
			autoTextColor = false
		}

		entries = append(entries, LegendOverlayEntry{
			Text:            text,
			TextColor:       textColor,
			AutoTextColor:   autoTextColor,
			BackgroundColor: bgColor,
			Glyph:           glyph,
		})
	}

	if len(entries) == 0 {
		return nil
	}

	return &LegendOverlay{
		Corner:  c.config.LegendCorner,
		Entries: entries,
	}
}

func pickReadableTextColor(bg tcell.Color) tcell.Color {
	r, g, b := bg.RGB()
	bgLum := relativeLuminance(r, g, b)

	// Choose the better contrast between pure black and pure white text.
	blackContrast := contrastRatio(bgLum, 0.0)
	whiteContrast := contrastRatio(bgLum, 1.0)
	if whiteContrast > blackContrast {
		return tcell.ColorWhite
	}
	return tcell.ColorBlack
}

func relativeLuminance(r, g, b int32) float64 {
	rf := linearizedSrgb(float64(r) / 255.0)
	gf := linearizedSrgb(float64(g) / 255.0)
	bf := linearizedSrgb(float64(b) / 255.0)
	return (0.2126 * rf) + (0.7152 * gf) + (0.0722 * bf)
}

func linearizedSrgb(v float64) float64 {
	if v <= 0.04045 {
		return v / 12.92
	}
	return math.Pow((v+0.055)/1.055, 2.4)
}

func contrastRatio(a, b float64) float64 {
	if a < b {
		a, b = b, a
	}
	return (a + 0.05) / (b + 0.05)
}

func (c *GraphComponent) setXAxisLabelFunc(series GraphSeries) {
	c.plotLayout.SetXAxisLabelFunc(func(i int) string {
		return series.GetXLabel(i)
	})
}

func (c *GraphComponent) GetSeries() []GraphSeries {
	return c.series
}

func (c *GraphComponent) GetLines() []*GraphLine {
	lines := make([]*GraphLine, 0, len(c.series))
	for _, s := range c.series {
		if line, ok := s.(*GraphLine); ok {
			lines = append(lines, line)
		}
	}
	return lines
}

func (c *GraphComponent) GetBars() []*GraphBar {
	bars := make([]*GraphBar, 0, len(c.series))
	for _, s := range c.series {
		if bar, ok := s.(*GraphBar); ok {
			bars = append(bars, bar)
		}
	}
	return bars
}

func (c *GraphComponent) GetXAxisZoomFactor() float64 {
	return c.getXAxisZoomFactor()
}

func (c *GraphComponent) SetXAxisZoomFactor(xAxisZoomFactor float64) {
	for _, series := range c.GetSeries() {
		series.SetXAxisZoomFactor(xAxisZoomFactor)
	}
	c.Refresh()
}

func (c *GraphComponent) GetXAxisShift() float64 {
	return c.getXAxisShift()
}

func (c *GraphComponent) SetXAxisShift(xAxisShift float64) {
	for _, series := range c.GetSeries() {
		series.SetXAxisShift(xAxisShift)
	}
	c.Refresh()
}

func (c *GraphComponent) GetYAxisZoomFactor() float64 {
	series := c.GetSeries()
	if len(series) == 0 {
		return 1.0
	}
	return series[0].GetYAxisZoomFactor()
}

func (c *GraphComponent) SetYAxisZoomFactor(yAxisZoomFactor float64) {
	for _, series := range c.GetSeries() {
		series.SetYAxisZoomFactor(yAxisZoomFactor)
	}
	c.Refresh()
}

func (c *GraphComponent) GetYAxisShift() float64 {
	series := c.GetSeries()
	if len(series) == 0 {
		return 0.0
	}
	return series[0].GetYAxisShift()
}

func (c *GraphComponent) SetYAxisShift(yAxisShift float64) {
	for _, series := range c.GetSeries() {
		series.SetYAxisShift(yAxisShift)
	}
	c.Refresh()
}

func (c *GraphComponent) GetPlotRect() (int, int, int, int) {
	return c.plotLayout.GetPlotRect()
}

func (c *GraphComponent) SetXRange(xMin, xMax float64) {
	for _, series := range c.GetSeries() {
		series.SetXRange(xMin, xMax)
	}
}

func (c *GraphComponent) ResetXRange() {
	for _, series := range c.GetSeries() {
		series.ResetXRange()
	}
}

func (c *GraphComponent) GetXMin() *float64 {
	return c.getXMin()
}

func (c *GraphComponent) GetXMax() *float64 {
	return c.getXMax()
}
