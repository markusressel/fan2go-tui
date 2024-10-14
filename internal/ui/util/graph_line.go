package util

import "math"

type GraphLine struct {
	HorizontalStretchFactor float64
	VerticalStretchFactor   float64
	XOffset                 float64
	YOffset                 float64

	xAxisZoomFactor float64
	xAxisShift      float64
	yAxisZoomFactor float64
	yAxisShift      float64

	Name string
	X    func(i int) float64
	F    func(float64) float64
}

func NewGraphLine(name string, xFunc func(i int) float64, fFunc func(float64) float64) *GraphLine {
	return &GraphLine{
		Name: name,

		HorizontalStretchFactor: 1.0,
		VerticalStretchFactor:   1.0,
		XOffset:                 0.0,
		YOffset:                 0.0,

		// TODO: set based on available width if data set is finite
		xAxisZoomFactor: 1.0,
		xAxisShift:      0.0,
		yAxisZoomFactor: 1.0,
		yAxisShift:      0.0,

		X: xFunc,
		F: fFunc,
	}
}

func (l *GraphLine) SetHorizontalStretchFactor(factor float64) {
	l.HorizontalStretchFactor = factor
}

func (l *GraphLine) SetVerticalStretchFactor(factor float64) {
	l.VerticalStretchFactor = factor
}

func (l *GraphLine) SetXOffset(offset float64) {
	l.XOffset = offset
}

func (l *GraphLine) SetYOffset(offset float64) {
	l.YOffset = offset
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
	return (l.GetF((x+l.XOffset)/l.HorizontalStretchFactor) + l.YOffset) * l.VerticalStretchFactor
}

func (l *GraphLine) GetX(i int) float64 {
	x := l.X(i)
	if math.IsNaN(x) {
		return math.NaN()
	}
	return (float64(i) / l.xAxisZoomFactor) + l.xAxisShift
}

func (l *GraphLine) GetF(x float64) float64 {
	return l.F(x)
}

func (l *GraphLine) GetName() string {
	return l.Name
}

func (l *GraphLine) GetHorizontalStretchFactor() float64 {
	return l.HorizontalStretchFactor
}

func (l *GraphLine) GetVerticalStretchFactor() float64 {
	return l.VerticalStretchFactor
}

func (l *GraphLine) GetXOffset() float64 {
	return l.XOffset
}

func (l *GraphLine) GetYOffset() float64 {
	return l.YOffset
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
	return l.X
}

func (l *GraphLine) GetFFunc() func(float64) float64 {
	return l.F
}
