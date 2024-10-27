package util

import (
	"github.com/gdamore/tcell/v2"
	"math"
)

type GraphLine struct {
	// Graph function manipulation settings
	horizontalStretchFactor float64
	verticalStretchFactor   float64
	xOffset                 float64
	yOffset                 float64

	// ViewPort settings
	xAxisZoomFactor float64
	xAxisShift      float64
	xMax            *float64

	yAxisZoomFactor float64
	yAxisShift      float64

	name       string
	x          func(i int) float64
	f          func(float64) float64
	xLabelFunc func(i int, x float64) string

	color tcell.Color
}

func NewGraphLine(
	name string,
	xFunc func(i int) float64,
	fFunc func(float64) float64,
	xLabelFunc func(i int, x float64) string,
) *GraphLine {
	return &GraphLine{
		name: name,

		horizontalStretchFactor: 1.0,
		verticalStretchFactor:   1.0,
		xOffset:                 0.0,
		yOffset:                 0.0,

		xAxisZoomFactor: 1.0,
		xAxisShift:      0.0,
		xMax:            nil,

		yAxisZoomFactor: 1.0,
		yAxisShift:      0.0,

		x:          xFunc,
		f:          fFunc,
		xLabelFunc: xLabelFunc,

		color: tcell.ColorWhite,
	}
}

func (l *GraphLine) SetHorizontalStretchFactor(factor float64) {
	l.horizontalStretchFactor = factor
}

func (l *GraphLine) SetVerticalStretchFactor(factor float64) {
	l.verticalStretchFactor = factor
}

func (l *GraphLine) SetXOffset(offset float64) {
	l.xOffset = offset
}

func (l *GraphLine) SetYOffset(offset float64) {
	l.yOffset = offset
}

func (l *GraphLine) SetXAxisZoomFactor(xAxisZoomFactor float64) {
	l.xAxisZoomFactor = xAxisZoomFactor
}

func (l *GraphLine) SetXAxisShift(xAxisShift float64) {
	l.xAxisShift = xAxisShift
}

func (l *GraphLine) SetYAxisZoomFactor(yAxisZoomFactor float64) {
	l.yAxisZoomFactor = yAxisZoomFactor
}

func (l *GraphLine) SetYAxisShift(yAxisShift float64) {
	l.yAxisShift = yAxisShift
}

func (l *GraphLine) GetX(i int) float64 {
	x := l.MapItoX(i)
	if math.IsNaN(x) {
		return math.NaN()
	}
	return x
}

func (l *GraphLine) GetF(x float64) float64 {
	return l.f(x)
}

func (l *GraphLine) GetY(x float64) float64 {
	targetXVal := (x + l.xOffset) / l.horizontalStretchFactor
	fVal := l.GetF(targetXVal)
	return (fVal + l.yOffset) * l.verticalStretchFactor
}

func (l *GraphLine) GetName() string {
	return l.name
}

func (l *GraphLine) GetHorizontalStretchFactor() float64 {
	return l.horizontalStretchFactor
}

func (l *GraphLine) GetVerticalStretchFactor() float64 {
	return l.verticalStretchFactor
}

func (l *GraphLine) GetXOffset() float64 {
	return l.xOffset
}

func (l *GraphLine) GetYOffset() float64 {
	return l.yOffset
}

func (l *GraphLine) GetXAxisZoomFactor() float64 {
	return l.xAxisZoomFactor
}

func (l *GraphLine) GetXAxisShift() float64 {
	return l.xAxisShift
}

func (l *GraphLine) GetYAxisZoomFactor() float64 {
	return l.yAxisZoomFactor
}

func (l *GraphLine) GetYAxisShift() float64 {
	return l.yAxisShift
}

func (l *GraphLine) GetXFunc() func(int) float64 {
	// TODO should this be static or a parameter?
	return func(i int) float64 {
		return l.MapItoX(i)
	}
}

func (l *GraphLine) GetFFunc() func(float64) float64 {
	return l.f
}

func (l *GraphLine) GetXLabel(i int) string {
	xVal := l.GetX(i)
	if math.IsNaN(xVal) {
		return ""
	}

	scaledX := xVal*l.xAxisZoomFactor + l.xAxisShift

	xMax := l.GetXMax()
	if xMax != nil && scaledX > *xMax {
		return ""
	}

	label := l.xLabelFunc(i, scaledX)
	return label
}

func (l *GraphLine) GetXMax() *float64 {
	return l.xMax
}

func (l *GraphLine) SetXRange(xMin, xMax float64) {
	l.xAxisShift = -1 * xMin
	l.xMax = &xMax
}

func (l *GraphLine) ResetXRange() {
	l.xAxisShift = 0
	l.xMax = nil
}

func (l *GraphLine) GetXRange() (float64, *float64) {
	return l.xAxisShift, l.xMax
}

func (l *GraphLine) GetColor() tcell.Color {
	return l.color
}

func (l *GraphLine) SetColor(color tcell.Color) {
	l.color = color
}

func (l *GraphLine) MapItoX(i int) float64 {
	return float64(i)*l.xAxisZoomFactor + l.xAxisShift
}

func (l *GraphLine) MapXtoI(x float64) int {
	return int((x - l.xAxisShift) / l.xAxisZoomFactor)
}
