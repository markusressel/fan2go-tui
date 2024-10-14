package util

import (
	"math"
)

type GraphLine struct {
	horizontalStretchFactor float64
	verticalStretchFactor   float64
	xOffset                 float64
	yOffset                 float64

	xAxisZoomFactor float64
	xAxisShift      float64
	yAxisZoomFactor float64
	yAxisShift      float64

	name       string
	x          func(i int) float64
	f          func(float64) float64
	xLabelFunc func(i int, x float64) string
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
		yAxisZoomFactor: 1.0,
		yAxisShift:      0.0,

		x:          xFunc,
		f:          fFunc,
		xLabelFunc: xLabelFunc,
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

func (l *GraphLine) GetY(x float64) float64 {
	return (l.GetF((x+l.xOffset)/l.horizontalStretchFactor) + l.yOffset) * l.verticalStretchFactor
}

func (l *GraphLine) GetX(i int) float64 {
	x := l.x(i)
	if math.IsNaN(x) {
		return math.NaN()
	}
	return (float64(i) / l.xAxisZoomFactor) + l.xAxisShift
}

func (l *GraphLine) GetF(x float64) float64 {
	return l.f(x)
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
	return l.x
}

func (l *GraphLine) GetFFunc() func(float64) float64 {
	return l.f
}

func (l *GraphLine) GetXLabel(i int) string {
	xVal := l.GetX(i)
	if math.IsNaN(xVal) {
		return ""
	}
	label := l.xLabelFunc(i, xVal)
	return label
}
