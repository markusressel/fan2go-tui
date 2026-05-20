package graph

import (
	"math"

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

type MarkerOverlay struct {
	coord func() XY
	color tcell.Color
	rune  rune
}

func NewMarkerOverlay(coord func() XY) *MarkerOverlay {
	if coord == nil {
		panic("marker overlay requires a coord function")
	}

	return &MarkerOverlay{
		coord: coord,
		color: tcell.ColorYellow,
		rune:  rune(0x2836),
	}
}

// Marker is a concise alias for NewMarkerOverlay.
func Marker(coord func() XY) *MarkerOverlay {
	return NewMarkerOverlay(coord)
}

func (o *MarkerOverlay) WithCoord(coord func() XY) *MarkerOverlay {
	if coord == nil {
		panic("marker overlay requires a coord function")
	}

	o.coord = coord
	return o
}

func (o *MarkerOverlay) WithColor(color tcell.Color) *MarkerOverlay {
	o.color = color
	return o
}

func (o *MarkerOverlay) WithRune(r rune) *MarkerOverlay {
	o.rune = r
	return o
}

func (o *MarkerOverlay) draw(screen tcell.Screen, ctx OverlayRenderContext) {
	if o.coord == nil || ctx.Plot == nil || ctx.XValueToIndex == nil || ctx.YMax <= ctx.YMin {
		return
	}

	coord := o.coord()
	if math.IsNaN(coord.X) || math.IsNaN(coord.Y) || math.IsInf(coord.X, 0) || math.IsInf(coord.Y, 0) {
		return
	}

	x, y, width, height := ctx.Plot.GetPlotRect()
	xIndex := ctx.XValueToIndex(coord.X)
	if xIndex < 0 || xIndex >= width {
		return
	}

	pointHeightFloat := ((coord.Y - ctx.YMin) / (ctx.YMax - ctx.YMin)) * float64(height-1)
	pointHeight := int(pointHeightFloat)
	if pointHeight < 0 || pointHeight >= height {
		return
	}

	pointStyle := tcell.StyleDefault.Background(ctx.Background).Foreground(o.color)
	screenX := x + xIndex
	screenY := y + height - 1 - pointHeight
	_, currentCombc, _, _ := screen.GetContent(screenX, screenY)

	markerRune := o.rune
	if markerRune == rune(0x2836) {
		// Shift a 4-dot block top/center/bottom depending on sub-cell Y position.
		fraction := pointHeightFloat - float64(pointHeight)
		switch {
		case fraction >= 2.0/3.0:
			markerRune = rune(0x281B) // rows 0+1
		case fraction <= 1.0/3.0:
			markerRune = rune(0x28E4) // rows 2+3
		default:
			markerRune = rune(0x2836) // rows 1+2
		}
	}

	screen.SetContent(screenX, screenY, markerRune, currentCombc, pointStyle)
}

type VerticalLine struct {
	x     func() float64
	color tcell.Color
	rune  rune
}

func NewVerticalLine(x func() float64) *VerticalLine {
	if x == nil {
		panic("vertical line overlay requires an x function")
	}

	return &VerticalLine{
		x:     x,
		color: tcell.ColorYellow,
		rune:  rune(0x2830),
	}
}

// VLine is a concise alias for NewVerticalLine.
func VLine(x func() float64) *VerticalLine {
	return NewVerticalLine(x)
}

func (o *VerticalLine) WithX(x func() float64) *VerticalLine {
	if x == nil {
		panic("vertical line overlay requires an x function")
	}

	o.x = x
	return o
}

func (o *VerticalLine) WithColor(color tcell.Color) *VerticalLine {
	o.color = color
	return o
}

func (o *VerticalLine) WithRune(r rune) *VerticalLine {
	o.rune = r
	return o
}

func (o *VerticalLine) draw(screen tcell.Screen, ctx OverlayRenderContext) {
	if o.x == nil || ctx.Plot == nil || ctx.XValueToIndex == nil {
		return
	}

	xValue := o.x()
	if math.IsNaN(xValue) || math.IsInf(xValue, 0) {
		return
	}

	x, y, width, height := ctx.Plot.GetPlotRect()
	xIndex := ctx.XValueToIndex(xValue)
	if xIndex < 0 || xIndex >= width {
		return
	}

	lineStyle := tcell.StyleDefault.Background(ctx.Background).Foreground(o.color)
	screenX := x + xIndex
	lineRune := o.rune
	if lineRune == rune(0x2830) {
		// Adapt default dotted vertical line glyph to the closest sub-cell x position.
		if ctx.XValueToIndexFloat != nil {
			xIndexFloat := ctx.XValueToIndexFloat(xValue)
			if !math.IsNaN(xIndexFloat) && !math.IsInf(xIndexFloat, 0) {
				fraction := xIndexFloat - math.Floor(xIndexFloat)
				if fraction < 0.5 {
					lineRune = rune(0x2806) // left column, middle two dots
				} else {
					lineRune = rune(0x2830) // right column, middle two dots
				}
			}
		}
	}

	for yPos := y; yPos < y+height; yPos++ {
		_, currentCombc, _, _ := screen.GetContent(screenX, yPos)
		screen.SetContent(screenX, yPos, lineRune, currentCombc, lineStyle)
	}
}
