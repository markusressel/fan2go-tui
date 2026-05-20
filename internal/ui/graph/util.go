package graph

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type AxisLabelTextFunc func(ctx OverlayRenderContext) string

func hasValidXYOverlayContext(ctx OverlayRenderContext) bool {
	return ctx.Plot != nil && ctx.XValueToIndex != nil && ctx.YMax > ctx.YMin
}

func isFiniteXY(point XY) bool {
	return !math.IsNaN(point.X) && !math.IsNaN(point.Y) && !math.IsInf(point.X, 0) && !math.IsInf(point.Y, 0)
}

func drawOverlayText(
	screen tcell.Screen,
	text string,
	startX int,
	y int,
	minX int,
	maxXExclusive int,
	fg tcell.Color,
	bg tcell.Color,
) {
	if text == "" || minX >= maxXExclusive {
		return
	}

	style := tcell.StyleDefault.Background(bg).Foreground(fg)
	runes := []rune(text)
	for i, r := range runes {
		x := startX + i
		if x < minX || x >= maxXExclusive {
			continue
		}

		_, combc, _, _ := screen.GetContent(x, y)
		screen.SetContent(x, y, r, combc, style)
	}
}
