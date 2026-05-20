package graph

import (
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

type XY struct {
	X float64
	Y float64
}

type OverlayRenderContext struct {
	Plot          *tvxwidgets.Plot
	XValueToIndex func(float64) int
	// XValueToIndexFloat preserves the fractional sub-cell x position.
	XValueToIndexFloat func(float64) float64
	YMin               float64
	YMax               float64
	Background         tcell.Color
	Bars               []*GraphBar
	ValueBufferSize    int
	Reversed           bool
	SeriesData         [][]float64
	SeriesColors       []tcell.Color
	YAxisLabelsAreInts bool
	LegendOverlay      *LegendOverlay
}

type GraphComponentOverlay interface {
	draw(screen tcell.Screen, ctx OverlayRenderContext)
}
