package graph

import (
	"github.com/gdamore/tcell/v2"
)

type AxisLabelTextFunc func(ctx OverlayRenderContext) string

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
