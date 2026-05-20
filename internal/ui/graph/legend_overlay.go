package graph

import (
	"github.com/gdamore/tcell/v2"
)

type LegendCorner int

const (
	LegendCornerTopLeft LegendCorner = iota
	LegendCornerTopRight
	LegendCornerBottomLeft
	LegendCornerBottomRight
)

type LegendOverlayEntry struct {
	Text            string
	TextColor       tcell.Color
	AutoTextColor   bool
	BackgroundColor tcell.Color
	Glyph           rune
}

type LegendOverlay struct {
	Corner  LegendCorner
	Entries []LegendOverlayEntry
}

func (o *LegendOverlay) draw(screen tcell.Screen, ctx OverlayRenderContext) {
	if o == nil || ctx.Plot == nil || len(o.Entries) == 0 {
		return
	}

	plotX, plotY, plotWidth, plotHeight := ctx.Plot.GetPlotRect()
	if plotWidth <= 0 || plotHeight <= 0 {
		return
	}

	maxWidth := 0
	visibleEntries := make([]LegendOverlayEntry, 0, len(o.Entries))
	for _, entry := range o.Entries {
		if entry.Text == "" {
			continue
		}
		visibleEntries = append(visibleEntries, entry)
		if w := 2 + len([]rune(entry.Text)); w > maxWidth {
			maxWidth = w
		}
	}
	if len(visibleEntries) == 0 || maxWidth <= 0 {
		return
	}

	height := len(visibleEntries)
	startX := plotX
	startY := plotY
	switch o.Corner {
	case LegendCornerTopRight:
		startX = plotX + plotWidth - maxWidth
	case LegendCornerBottomLeft:
		startY = plotY + plotHeight - height
	case LegendCornerBottomRight:
		startX = plotX + plotWidth - maxWidth
		startY = plotY + plotHeight - height
	}

	if startY < plotY {
		startY = plotY
	}

	for i, entry := range visibleEntries {
		y := startY + i
		if y < plotY || y >= plotY+plotHeight {
			continue
		}

		textColor := entry.TextColor
		if entry.AutoTextColor {
			textColor = pickReadableTextColor(ctx.Background)
		}

		swatchX := startX
		if swatchX >= plotX && swatchX < plotX+plotWidth {
			swatchStyle := tcell.StyleDefault.Background(ctx.Background).Foreground(entry.BackgroundColor)
			_, swatchCombc, _, _ := screen.GetContent(swatchX, y)
			swatchGlyph := entry.Glyph
			if swatchGlyph == 0 {
				swatchGlyph = rune(0x25CF) // ●
			}
			screen.SetContent(swatchX, y, swatchGlyph, swatchCombc, swatchStyle)
		}

		spacerX := startX + 1
		if spacerX >= plotX && spacerX < plotX+plotWidth {
			spacerStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
			_, spacerCombc, _, _ := screen.GetContent(spacerX, y)
			screen.SetContent(spacerX, y, ' ', spacerCombc, spacerStyle)
		}

		drawOverlayText(
			screen,
			entry.Text,
			startX+2,
			y,
			plotX,
			plotX+plotWidth,
			textColor,
			ctx.Background,
		)
	}
}
