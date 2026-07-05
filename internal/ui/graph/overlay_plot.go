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
	overlays       []GraphComponentOverlay
	overlayCtx     OverlayRenderContext
	onLayoutChange func()
}

type brailleLineCell struct {
	bits  rune
	color tcell.Color
}

type barBrailleCell struct {
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

func (p *OverlayPlot) SetOnLayoutChange(f func()) {
	p.onLayoutChange = f
}

func (p *OverlayPlot) Draw(screen tcell.Screen) {
	p.Plot.Draw(screen)

	_, _, width, height := p.Plot.GetPlotRect()
	if width <= 0 || height <= 0 {
		return
	}

	ctx := p.overlayCtx

	// If the data context was computed for a different width,
	// update synchronously and redraw to avoid any visual glitches.
	if ctx.ValueBufferSize > 0 && ctx.ValueBufferSize != width {
		if p.onLayoutChange != nil {
			p.onLayoutChange()
			ctx = p.overlayCtx
			p.Plot.Draw(screen) // Redraw axes with the correct data context
		} else {
			return
		}
	}

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
		runHasSegment := false

		for i := 0; i < displaySlots; i++ {
			sourceIndex := i
			if ctx.Reversed {
				sourceIndex = availableCount - 1 - i
			} else {
				sourceIndex = i - (displaySlots - availableCount)
			}

			if sourceIndex < 0 || sourceIndex >= displaySlots {
				if hasPrev && !runHasSegment {
					p.setLinePoint(cells, prevSubX, prevSubY, totalSubCols, totalSubRows, seriesColor)
				}
				hasPrev = false
				runHasSegment = false
				continue
			}

			val := series[sourceIndex]
			if math.IsNaN(val) || math.IsInf(val, 0) {
				if hasPrev && !runHasSegment {
					p.setLinePoint(cells, prevSubX, prevSubY, totalSubCols, totalSubRows, seriesColor)
				}
				hasPrev = false
				runHasSegment = false
				continue
			}

			pointSubY, ok := mapValueToBrailleSubRow(val, ctx.YMin, ctx.YMax, totalSubRows)
			if !ok {
				if hasPrev && !runHasSegment {
					p.setLinePoint(cells, prevSubX, prevSubY, totalSubCols, totalSubRows, seriesColor)
				}
				hasPrev = false
				runHasSegment = false
				continue
			}
			pointSubX := (i * 2) + 1

			if hasPrev {
				startSubX, endSubX := segmentSubXEndpoints(i-1, prevSubY, i, pointSubY)
				p.drawLineSegment(cells, startSubX, prevSubY, endSubX, pointSubY, totalSubCols, totalSubRows, seriesColor)
				runHasSegment = true
			} else {
				runHasSegment = false
			}

			hasPrev = true
			prevSubX = pointSubX
			prevSubY = pointSubY
		}

		if hasPrev && !runHasSegment {
			p.setLinePoint(cells, prevSubX, prevSubY, totalSubCols, totalSubRows, seriesColor)
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

func segmentSubXEndpoints(startCellX, startSubY, endCellX, endSubY int) (int, int) {
	startLeft := startCellX * 2
	startRight := startLeft + 1
	endLeft := endCellX * 2
	endRight := endLeft + 1

	if endSubY > startSubY {
		return startLeft, endRight
	}
	return startRight, endLeft
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

func barFillBits(level, column int) rune {
	if level <= 0 || column < 0 || column > 1 {
		return 0
	}
	if level > 4 {
		level = 4
	}

	bits := rune(0)
	for i := 0; i < level; i++ {
		rowFromTop := 3 - i // fill from bottom to top
		bits |= brailleLineBit[rowFromTop][column]
	}
	return bits
}

func resolveBarSubX(ctx OverlayRenderContext, xVal float64, xIndex int) int {
	if ctx.XValueToIndexFloat != nil {
		xFloat := ctx.XValueToIndexFloat(xVal)
		if !math.IsNaN(xFloat) && !math.IsInf(xFloat, 0) {
			return int(math.Round(xFloat * 2))
		}
	}
	return (xIndex * 2) + 1
}

func (p *OverlayPlot) fillBarAtSubX(
	cells map[[2]int]barBrailleCell,
	subX, filledSubRows, height, totalSubCols int,
	ctx OverlayRenderContext,
	stops []GraphBarGradientStop,
	fallbackColor tcell.Color,
) {
	if subX < 0 || subX >= totalSubCols || height <= 0 || filledSubRows <= 0 {
		return
	}

	totalSubRows := height * 4
	if filledSubRows > totalSubRows {
		filledSubRows = totalSubRows
	}

	cellX := subX / 2
	column := subX % 2

	for cellOffset := 0; cellOffset < height; cellOffset++ {
		subRowsInCell := filledSubRows - cellOffset*4
		if subRowsInCell <= 0 {
			break
		}
		if subRowsInCell > 4 {
			subRowsInCell = 4
		}

		bits := barFillBits(subRowsInCell, column)
		if bits == 0 {
			continue
		}

		normalizedVertical := 0.0
		if height > 1 {
			normalizedVertical = float64(cellOffset) / float64(height-1)
		}
		yAtCell := ctx.YMin + (normalizedVertical * (ctx.YMax - ctx.YMin))
		gradientColor := colorAtY(yAtCell, stops, fallbackColor)

		key := [2]int{cellX, height - 1 - cellOffset}
		cell := cells[key]
		cell.bits |= bits
		cell.color = gradientColor
		cells[key] = cell
	}
}

func (p *OverlayPlot) fillBarSegment(
	cells map[[2]int]barBrailleCell,
	startSubX, startFilledSubRows, endSubX, endFilledSubRows, height, totalSubCols int,
	ctx OverlayRenderContext,
	stops []GraphBarGradientStop,
	fallbackColor tcell.Color,
) {
	if startSubX == endSubX {
		fill := startFilledSubRows
		if endFilledSubRows > fill {
			fill = endFilledSubRows
		}
		p.fillBarAtSubX(cells, startSubX, fill, height, totalSubCols, ctx, stops, fallbackColor)
		return
	}

	deltaX := endSubX - startSubX
	step := 1
	if deltaX < 0 {
		deltaX = -deltaX
		step = -1
	}

	for offset := 0; offset <= deltaX; offset++ {
		t := float64(offset) / float64(deltaX)
		filled := int(math.Round(float64(startFilledSubRows) + (float64(endFilledSubRows-startFilledSubRows) * t)))
		subX := startSubX + (offset * step)
		p.fillBarAtSubX(cells, subX, filled, height, totalSubCols, ctx, stops, fallbackColor)
	}
}

func (p *OverlayPlot) drawBars(screen tcell.Screen, ctx OverlayRenderContext) {
	if len(ctx.Bars) == 0 || ctx.ValueBufferSize <= 0 || ctx.YMax <= ctx.YMin || ctx.XValueToIndex == nil {
		return
	}

	x, y, width, height := p.Plot.GetPlotRect()
	totalSubRows := height * 4
	totalSubCols := width * 2
	cells := map[[2]int]barBrailleCell{}

	for _, bar := range ctx.Bars {
		stops := bar.GetGradientStops(ctx.YMin, ctx.YMax)
		hasPrevPoint := false
		prevSubX := 0
		prevFilledSubRows := 0

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
				hasPrevPoint = false
				continue
			}

			subX := resolveBarSubX(ctx, displayXVal, xIndex)
			if subX < 0 || subX >= totalSubCols {
				hasPrevPoint = false
				continue
			}

			if hasPrevPoint {
				p.fillBarSegment(cells, prevSubX, prevFilledSubRows, subX, filledSubRows, height, totalSubCols, ctx, stops, bar.GetColor())
			} else {
				p.fillBarAtSubX(cells, subX, filledSubRows, height, totalSubCols, ctx, stops, bar.GetColor())
			}

			hasPrevPoint = true
			prevSubX = subX
			prevFilledSubRows = filledSubRows
		}
	}

	for point, cell := range cells {
		screenX := x + point[0]
		screenY := y + point[1]
		if screenX < x || screenX >= x+width || screenY < y || screenY >= y+height {
			continue
		}

		_, combc, _, _ := screen.GetContent(screenX, screenY)
		barStyle := tcell.StyleDefault.Background(ctx.Background).Foreground(cell.color)
		screen.SetContent(screenX, screenY, rune(0x2800)+cell.bits, combc, barStyle)
	}
}
