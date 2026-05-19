package graph

import (
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

// OverlayPlot extends tvxwidgets.Plot with lightweight overlay rendering.
type OverlayPlot[T any] struct {
	*tvxwidgets.Plot
	overlays   []GraphComponentOverlay[T]
	overlayCtx OverlayRenderContext[T]
}

func NewOverlayPlot[T any]() *OverlayPlot[T] {
	return &OverlayPlot[T]{
		Plot: tvxwidgets.NewPlot(),
	}
}

func (p *OverlayPlot[T]) SetOverlays(overlays []GraphComponentOverlay[T]) {
	p.overlays = append([]GraphComponentOverlay[T]{}, overlays...)
}

func (p *OverlayPlot[T]) SetOverlayContext(ctx OverlayRenderContext[T]) {
	p.overlayCtx = ctx
}

func (p *OverlayPlot[T]) Draw(screen tcell.Screen) {
	p.Plot.Draw(screen)
	ctx := p.overlayCtx
	ctx.Plot = p.Plot
	ctx.Background = p.GetBackgroundColor()
	p.drawBars(screen, ctx)
	for _, overlay := range p.overlays {
		overlay.draw(screen, ctx)
	}
}

func barFillRune(level int) rune {
	switch {
	case level <= 0:
		return rune(0x2800)
	case level == 1:
		return rune(0x28C0) // bottom row, both columns
	case level == 2:
		return rune(0x28E4) // bottom two rows, both columns
	case level == 3:
		return rune(0x28F6) // bottom three rows, both columns
	default:
		return rune(0x28FF) // full cell
	}
}

func (p *OverlayPlot[T]) drawBars(screen tcell.Screen, ctx OverlayRenderContext[T]) {
	if len(ctx.Bars) == 0 || ctx.ValueBufferSize <= 0 || ctx.YMax <= ctx.YMin || ctx.XValueToIndex == nil {
		return
	}

	x, y, width, height := p.Plot.GetPlotRect()
	totalSubRows := height * 4

	for _, bar := range ctx.Bars {
		barStyle := tcell.StyleDefault.Background(ctx.Background).Foreground(bar.GetColor())

		availableCount := 0
		for sourceIdx := 0; sourceIdx < ctx.ValueBufferSize; sourceIdx++ {
			xVal := bar.GetX(sourceIdx)
			if math.IsNaN(xVal) || math.IsInf(xVal, 0) {
				continue
			}
			yVal := bar.GetY(xVal)
			if !math.IsNaN(yVal) && !math.IsInf(yVal, 0) {
				availableCount++
			}
		}

		for i := 0; i < ctx.ValueBufferSize; i++ {
			displayXVal := bar.GetX(i)
			if math.IsNaN(displayXVal) || math.IsInf(displayXVal, 0) {
				continue
			}

			xIndex := ctx.XValueToIndex(displayXVal)
			if xIndex < 0 || xIndex >= width {
				continue
			}

			sourceIndex := i
			if ctx.Reversed {
				sourceIndex = availableCount - 1 - i
			} else {
				sourceIndex = i - (ctx.ValueBufferSize - availableCount)
			}

			if sourceIndex < 0 || sourceIndex >= ctx.ValueBufferSize {
				continue
			}

			sourceXVal := bar.GetX(sourceIndex)
			if math.IsNaN(sourceXVal) || math.IsInf(sourceXVal, 0) {
				continue
			}

			yVal := bar.GetY(sourceXVal)
			if math.IsNaN(yVal) || math.IsInf(yVal, 0) || yVal <= ctx.YMin {
				continue
			}
			if yVal > ctx.YMax {
				yVal = ctx.YMax
			}

			filledSubRows := int(math.Round(((yVal - ctx.YMin) / (ctx.YMax - ctx.YMin)) * float64(totalSubRows)))
			if filledSubRows <= 0 {
				continue
			}

			screenX := x + xIndex
			for cellOffset := 0; cellOffset < height; cellOffset++ {
				subRowsInCell := filledSubRows - cellOffset*4
				if subRowsInCell <= 0 {
					break
				}

				if subRowsInCell > 4 {
					subRowsInCell = 4
				}

				r := barFillRune(subRowsInCell)
				yPos := y + height - 1 - cellOffset
				_, combc, _, _ := screen.GetContent(screenX, yPos)
				screen.SetContent(screenX, yPos, r, combc, barStyle)
			}
		}
	}
}
