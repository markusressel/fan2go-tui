package graph

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

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
