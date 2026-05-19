package graph

import (
	"math"
	"sort"

	"github.com/gdamore/tcell/v2"
)

type GraphBarGradientStop struct {
	YValue float64
	Color  tcell.Color
}

type GraphBarGradientFunc func(yMin, yMax float64) []GraphBarGradientStop

// GraphBar describes a function rendered as a filled bar from the baseline up to f(x).
type GraphBar struct {
	horizontalStretchFactor float64
	verticalStretchFactor   float64
	xOffset                 float64
	yOffset                 float64

	xAxisZoomFactor float64
	xAxisShift      float64

	xMax *float64
	xMin *float64

	yAxisZoomFactor float64
	yAxisShift      float64

	yMax *float64
	yMin *float64

	name       string
	x          func(i int) float64
	f          func(float64) float64
	xLabelFunc func(i int, x float64) string

	color        tcell.Color
	gradientFunc GraphBarGradientFunc
}

func NewGraphBar(
	name string,
	xFunc func(i int) float64,
	fFunc func(float64) float64,
	xLabelFunc func(i int, x float64) string,
) *GraphBar {
	return &GraphBar{
		name: name,

		horizontalStretchFactor: 1.0,
		verticalStretchFactor:   1.0,
		xOffset:                 0.0,
		yOffset:                 0.0,

		xAxisZoomFactor: 1.0,
		xAxisShift:      0.0,
		xMax:            nil,
		xMin:            nil,

		yAxisZoomFactor: 1.0,
		yAxisShift:      0.0,
		yMax:            nil,
		yMin:            nil,

		x:          xFunc,
		f:          fFunc,
		xLabelFunc: xLabelFunc,

		color: tcell.ColorWhite,
	}
}

func (b *GraphBar) SetHorizontalStretchFactor(factor float64) { b.horizontalStretchFactor = factor }
func (b *GraphBar) SetVerticalStretchFactor(factor float64)   { b.verticalStretchFactor = factor }
func (b *GraphBar) SetXOffset(offset float64)                 { b.xOffset = offset }
func (b *GraphBar) SetYOffset(offset float64)                 { b.yOffset = offset }

func (b *GraphBar) SetXAxisZoomFactor(xAxisZoomFactor float64) { b.xAxisZoomFactor = xAxisZoomFactor }
func (b *GraphBar) SetXAxisShift(xAxisShift float64)           { b.xAxisShift = xAxisShift }
func (b *GraphBar) SetYAxisZoomFactor(yAxisZoomFactor float64) { b.yAxisZoomFactor = yAxisZoomFactor }
func (b *GraphBar) SetYAxisShift(yAxisShift float64)           { b.yAxisShift = yAxisShift }

func (b *GraphBar) GetX(i int) float64 {
	x := b.MapItoX(i)
	if math.IsNaN(x) {
		return math.NaN()
	}
	return x
}

func (b *GraphBar) GetY(x float64) float64 {
	targetXVal := (x + b.xOffset) / b.horizontalStretchFactor
	fVal := b.f(targetXVal)
	return (fVal + b.yOffset) * b.verticalStretchFactor
}

func (b *GraphBar) GetXLabel(i int) string {
	xVal := b.GetX(i)
	if math.IsNaN(xVal) {
		return ""
	}

	xMax := b.GetXMax()
	if xMax != nil && xVal > *xMax {
		return ""
	}

	return b.xLabelFunc(i, xVal)
}

func (b *GraphBar) GetXAxisZoomFactor() float64 { return b.xAxisZoomFactor }
func (b *GraphBar) GetXAxisShift() float64      { return b.xAxisShift }
func (b *GraphBar) GetYAxisZoomFactor() float64 { return b.yAxisZoomFactor }
func (b *GraphBar) GetYAxisShift() float64      { return b.yAxisShift }
func (b *GraphBar) GetYOffset() float64         { return b.yOffset }

func (b *GraphBar) SetXRange(xMin, xMax float64) {
	b.xAxisShift = -1 * xMin
	b.xMax = &xMax
}

func (b *GraphBar) ResetXRange() {
	b.xAxisShift = 0
	b.xMax = nil
}

func (b *GraphBar) SetYRange(yMin, yMax float64) {
	b.yAxisShift = -1 * yMin
	b.yMax = &yMax
}

func (b *GraphBar) ResetYRange() {
	b.yAxisShift = 0
	b.yMax = nil
}

func (b *GraphBar) GetXMax() *float64 { return b.xMax }
func (b *GraphBar) GetXMin() *float64 { return b.xMin }
func (b *GraphBar) GetYMax() *float64 { return b.yMax }
func (b *GraphBar) GetYMin() *float64 { return b.yMin }

func (b *GraphBar) GetColor() tcell.Color { return b.color }
func (b *GraphBar) SetColor(color tcell.Color) {
	b.color = color
}

func (b *GraphBar) SetGradient(gradientFunc GraphBarGradientFunc) {
	b.gradientFunc = gradientFunc
}

func (b *GraphBar) WithGradient(gradientFunc GraphBarGradientFunc) *GraphBar {
	b.SetGradient(gradientFunc)
	return b
}

func (b *GraphBar) GetGradientStops(yMin, yMax float64) []GraphBarGradientStop {
	if b.gradientFunc == nil {
		return nil
	}

	stops := append([]GraphBarGradientStop{}, b.gradientFunc(yMin, yMax)...)
	sort.Slice(stops, func(i, j int) bool {
		return stops[i].YValue < stops[j].YValue
	})

	return stops
}

func (b *GraphBar) MapItoX(i int) float64 {
	return float64(i)*b.xAxisZoomFactor + b.xAxisShift
}

func (b *GraphBar) MapXtoI(x float64) int {
	return int((x - b.xAxisShift) / b.xAxisZoomFactor)
}
