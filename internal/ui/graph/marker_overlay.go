package graph

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

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
		rune:  rune(0x2836), // ⠶
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
