package graph

import (
	coreutil "fan2go-tui/internal/util"
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

func lerpColor(a, b tcell.Color, t float64) tcell.Color {
	t = coreutil.Clamp01(t)
	ar, ag, ab := a.RGB()
	br, bg, bb := b.RGB()
	r := int32(math.Round(float64(ar) + (float64(br-ar) * t)))
	g := int32(math.Round(float64(ag) + (float64(bg-ag) * t)))
	bv := int32(math.Round(float64(ab) + (float64(bb-ab) * t)))
	return tcell.NewRGBColor(r, g, bv)
}

func colorAtY(y float64, stops []GraphBarGradientStop, fallback tcell.Color) tcell.Color {
	if len(stops) == 0 {
		return fallback
	}

	if y <= stops[0].YValue {
		return stops[0].Color
	}
	last := stops[len(stops)-1]
	if y >= last.YValue {
		return last.Color
	}

	for i := 1; i < len(stops); i++ {
		prev := stops[i-1]
		next := stops[i]
		if y <= next.YValue {
			rangeY := next.YValue - prev.YValue
			if rangeY <= 0 {
				return next.Color
			}
			t := (y - prev.YValue) / rangeY
			return lerpColor(prev.Color, next.Color, t)
		}
	}

	return fallback
}

// OverlayPlot extends tvxwidgets.Plot with lightweight overlay rendering.
type OverlayPlot struct {
	*tvxwidgets.Plot
	overlays   []GraphComponentOverlay
	overlayCtx OverlayRenderContext
}

type brailleLineCell struct {
	bits  rune
	color tcell.Color
}

var brailleLineBit = [4][2]rune{
	{rune(0x0001), rune(0x0008)},
	{rune(0x0002), rune(0x0010)},
	{rune(0x0004), rune(0x0020)},
	{rune(0x0040), rune(0x0080)},
}

func NewOverlayPlot() *OverlayPlot {
	return &OverlayPlot{
		Plot: tvxwidgets.NewPlot(),
	}
}

func (p *OverlayPlot) SetOverlays(overlays []GraphComponentOverlay) {
	p.overlays = append([]GraphComponentOverlay{}, overlays...)
}

func (p *OverlayPlot) SetOverlayContext(ctx OverlayRenderContext) {
	p.overlayCtx = ctx
}

func (p *OverlayPlot) Draw(screen tcell.Screen) {
	p.Plot.Draw(screen)
	ctx := p.overlayCtx
	ctx.Plot = p.Plot
	ctx.Background = p.GetBackgroundColor()
	p.drawBars(screen, ctx)
	p.drawLineSeries(screen, ctx)
	for _, overlay := range p.overlays {
		overlay.draw(screen, ctx)
	}
	if ctx.LegendOverlay != nil {
		ctx.LegendOverlay.draw(screen, ctx)
	}
}

func (p *OverlayPlot) drawLineSeries(screen tcell.Screen, ctx OverlayRenderContext) {
	if len(ctx.SeriesData) == 0 || ctx.YMax <= ctx.YMin {
		return
	}

	x, y, width, height := p.Plot.GetPlotRect()
	if width <= 0 || height <= 0 {
		return
	}
	totalSubRows := height * 4
	totalSubCols := width * 2

	for seriesIdx, series := range ctx.SeriesData {
		displaySlots := len(series)
		if ctx.ValueBufferSize > 0 && ctx.ValueBufferSize < displaySlots {
			displaySlots = ctx.ValueBufferSize
		}
		if width < displaySlots {
			displaySlots = width
		}

		availableCount := 0
		for sourceIdx := 0; sourceIdx < displaySlots; sourceIdx++ {
			val := series[sourceIdx]
			if math.IsNaN(val) || math.IsInf(val, 0) {
				continue
			}
			availableCount++
		}

		seriesColor := tcell.ColorWhite
		if seriesIdx < len(ctx.SeriesColors) {
			seriesColor = ctx.SeriesColors[seriesIdx]
		}

		cells := map[[2]int]brailleLineCell{}
		hasPrev := false
		prevSubX := 0
		prevSubY := 0

		for i := 0; i < displaySlots; i++ {
			sourceIndex := i
			if ctx.Reversed {
				sourceIndex = availableCount - 1 - i
			} else {
				sourceIndex = i - (displaySlots - availableCount)
			}

			if sourceIndex < 0 || sourceIndex >= displaySlots {
				hasPrev = false
				continue
			}

			val := series[sourceIndex]
			if math.IsNaN(val) || math.IsInf(val, 0) {
				hasPrev = false
				continue
			}

			pointSubY, ok := mapValueToBrailleSubRow(val, ctx.YMin, ctx.YMax, totalSubRows)
			if !ok {
				hasPrev = false
				continue
			}
			pointSubX := (i * 2) + 1 // Anchor samples at the center of each cell to leverage both braille columns.

			if hasPrev {
				p.drawLineSegment(cells, prevSubX, prevSubY, pointSubX, pointSubY, totalSubCols, totalSubRows, seriesColor)
			} else {
				p.setLinePoint(cells, pointSubX, pointSubY, totalSubCols, totalSubRows, seriesColor)
			}

			hasPrev = true
			prevSubX = pointSubX
			prevSubY = pointSubY
		}

		for point, cell := range cells {
			screenX := x + point[0]
			screenY := y + point[1]
			if screenX < x || screenX >= x+width || screenY < y || screenY >= y+height {
				continue
			}

			_, combc, _, _ := screen.GetContent(screenX, screenY)
			style := tcell.StyleDefault.Background(ctx.Background).Foreground(cell.color)
			screen.SetContent(screenX, screenY, rune(0x2800)+cell.bits, combc, style)
		}
	}
}

func mapValueToBrailleSubRow(val, yMin, yMax float64, totalSubRows int) (int, bool) {
	if totalSubRows <= 0 || yMax <= yMin || math.IsNaN(val) || math.IsInf(val, 0) {
		return 0, false
	}

	normalized := (val - yMin) / (yMax - yMin)
	if math.IsNaN(normalized) || math.IsInf(normalized, 0) {
		return 0, false
	}

	subRowFromBottom := int(math.Round(normalized * float64(totalSubRows-1)))
	if subRowFromBottom < 0 {
		subRowFromBottom = 0
	} else if subRowFromBottom >= totalSubRows {
		subRowFromBottom = totalSubRows - 1
	}

	return (totalSubRows - 1) - subRowFromBottom, true
}

func (p *OverlayPlot) drawLineSegment(cells map[[2]int]brailleLineCell, sx0, sy0, sx1, sy1, totalSubCols, totalSubRows int, color tcell.Color) {
	dx := int(math.Abs(float64(sx1 - sx0)))
	dy := -int(math.Abs(float64(sy1 - sy0)))
	sx := -1
	if sx0 < sx1 {
		sx = 1
	}
	sy := -1
	if sy0 < sy1 {
		sy = 1
	}
	err := dx + dy

	for {
		p.setLinePoint(cells, sx0, sy0, totalSubCols, totalSubRows, color)
		if sx0 == sx1 && sy0 == sy1 {
			break
		}

		e2 := 2 * err
		if e2 >= dy {
			err += dy
			sx0 += sx
		}
		if e2 <= dx {
			err += dx
			sy0 += sy
		}
	}
}

func (p *OverlayPlot) setLinePoint(cells map[[2]int]brailleLineCell, sx, sy, totalSubCols, totalSubRows int, color tcell.Color) {
	if sx < 0 || sx >= totalSubCols || sy < 0 || sy >= totalSubRows {
		return
	}

	cellX := sx / 2
	cellY := sy / 4
	bit := brailleLineBit[sy%4][sx%2]

	key := [2]int{cellX, cellY}
	cell := cells[key]
	cell.bits |= bit
	cell.color = color
	cells[key] = cell
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

func (p *OverlayPlot) drawBars(screen tcell.Screen, ctx OverlayRenderContext) {
	if len(ctx.Bars) == 0 || ctx.ValueBufferSize <= 0 || ctx.YMax <= ctx.YMin || ctx.XValueToIndex == nil {
		return
	}

	x, y, width, height := p.Plot.GetPlotRect()
	totalSubRows := height * 4

	for _, bar := range ctx.Bars {
		stops := bar.GetGradientStops(ctx.YMin, ctx.YMax)

		availableCount := 0
		for sourceIdx := 0; sourceIdx < ctx.ValueBufferSize; sourceIdx++ {
			xVal := bar.GetX(sourceIdx)
			if math.IsNaN(xVal) || math.IsInf(xVal, 0) {
				continue
			}
			yVal := bar.GetY(xVal)
			if !math.IsNaN(yVal) && !math.IsInf(yVal, 0) && yVal > ctx.YMin {
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
				normalizedVertical := 0.0
				if height > 1 {
					normalizedVertical = float64(cellOffset) / float64(height-1)
				}
				yAtCell := ctx.YMin + (normalizedVertical * (ctx.YMax - ctx.YMin))
				gradientColor := colorAtY(yAtCell, stops, bar.GetColor())
				barStyle := tcell.StyleDefault.Background(ctx.Background).Foreground(gradientColor)
				_, combc, _, _ := screen.GetContent(screenX, yPos)
				screen.SetContent(screenX, yPos, r, combc, barStyle)
			}
		}
	}
}
