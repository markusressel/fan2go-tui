package graph

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type HorizontalLine struct {
	y     func() float64
	color tcell.Color
	rune  rune
}

func NewHorizontalLine(y func() float64) *HorizontalLine {
	if y == nil {
		panic("horizontal line overlay requires a y function")
	}

	return &HorizontalLine{
		y:     y,
		color: tcell.ColorYellow,
		rune:  rune(0x2812),
	}
}

// HLine is a concise alias for NewHorizontalLine.
func HLine(y func() float64) *HorizontalLine {
	return NewHorizontalLine(y)
}

func (o *HorizontalLine) WithY(y func() float64) *HorizontalLine {
	if y == nil {
		panic("horizontal line overlay requires a y function")
	}

	o.y = y
	return o
}

func (o *HorizontalLine) WithColor(color tcell.Color) *HorizontalLine {
	o.color = color
	return o
}

func (o *HorizontalLine) WithRune(r rune) *HorizontalLine {
	o.rune = r
	return o
}

func (o *HorizontalLine) draw(screen tcell.Screen, ctx OverlayRenderContext) {
	if o.y == nil || ctx.Plot == nil || ctx.YMax <= ctx.YMin {
		return
	}

	yValue := o.y()
	if math.IsNaN(yValue) || math.IsInf(yValue, 0) {
		return
	}

	x, y, width, height := ctx.Plot.GetPlotRect()
	if width <= 0 || height <= 0 {
		return
	}

	pointHeightFloat := ((yValue - ctx.YMin) / (ctx.YMax - ctx.YMin)) * float64(height-1)
	pointHeight := int(pointHeightFloat)
	if pointHeight < 0 || pointHeight >= height {
		return
	}

	screenY := y + height - 1 - pointHeight
	lineStyle := tcell.StyleDefault.Background(ctx.Background).Foreground(o.color)
	lineRune := o.rune
	if lineRune == rune(0x2812) {
		fraction := pointHeightFloat - float64(pointHeight)
		switch {
		case fraction >= 0.75:
			lineRune = rune(0x2809)
		case fraction >= 0.5:
			lineRune = rune(0x2812)
		case fraction >= 0.25:
			lineRune = rune(0x2824)
		default:
			lineRune = rune(0x28C0)
		}
	}

	for xPos := x; xPos < x+width; xPos++ {
		_, currentCombc, _, _ := screen.GetContent(xPos, screenY)
		screen.SetContent(xPos, screenY, lineRune, currentCombc, lineStyle)
	}
}
