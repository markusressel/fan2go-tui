package graph

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

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

	name       string
	x          func(i int) float64
	f          func(float64) float64
	xLabelFunc func(i int, x float64) string

	color tcell.Color
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

func (b *GraphBar) SetXRange(xMin, xMax float64) {
	b.xAxisShift = -1 * xMin
	b.xMax = &xMax
}

func (b *GraphBar) ResetXRange() {
	b.xAxisShift = 0
	b.xMax = nil
}

func (b *GraphBar) GetXMax() *float64 { return b.xMax }
func (b *GraphBar) GetXMin() *float64 { return b.xMin }

func (b *GraphBar) GetColor() tcell.Color { return b.color }
func (b *GraphBar) SetColor(color tcell.Color) {
	b.color = color
}

func (b *GraphBar) MapItoX(i int) float64 {
	return float64(i)*b.xAxisZoomFactor + b.xAxisShift
}

func (b *GraphBar) MapXtoI(x float64) int {
	return int((x - b.xAxisShift) / b.xAxisZoomFactor)
}
