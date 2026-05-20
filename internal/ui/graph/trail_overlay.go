package graph

import (
	"github.com/gdamore/tcell/v2"
)

type TrailOverlay struct {
	points    func() []XY
	color     tcell.Color
	hasColor  bool
	maxPoints int
	opacity   float64
}

func NewTrailOverlay(points func() []XY) *TrailOverlay {
	if points == nil {
		panic("trail overlay requires a points function")
	}

	return &TrailOverlay{
		points:    points,
		color:     tcell.ColorOrange,
		hasColor:  false,
		maxPoints: 200,
		opacity:   0.6,
	}
}

// Trail is a concise alias for NewTrailOverlay.
func Trail(points func() []XY) *TrailOverlay {
	return NewTrailOverlay(points)
}

func (o *TrailOverlay) WithPoints(points func() []XY) *TrailOverlay {
	if points == nil {
		panic("trail overlay requires a points function")
	}
	o.points = points
	return o
}

func (o *TrailOverlay) WithColor(color tcell.Color) *TrailOverlay {
	o.color = color
	o.hasColor = true
	return o
}

func (o *TrailOverlay) WithMaxPoints(maxPoints int) *TrailOverlay {
	o.maxPoints = maxPoints
	return o
}

// WithOpacity sets the maximum trail opacity in [0,1]. Lower values make the trail subtler.
func (o *TrailOverlay) WithOpacity(opacity float64) *TrailOverlay {
	o.opacity = coreClamp01(opacity)
	return o
}

func (o *TrailOverlay) draw(screen tcell.Screen, ctx OverlayRenderContext) {
	if o.points == nil || !hasValidXYOverlayContext(ctx) {
		return
	}

	points := o.points()
	if len(points) == 0 {
		return
	}
	if o.maxPoints > 0 && len(points) > o.maxPoints {
		points = points[len(points)-o.maxPoints:]
	}

	x, y, width, height := ctx.Plot.GetPlotRect()
	if width <= 0 || height <= 0 {
		return
	}

	intensityByCell := make(map[[2]int]float64, len(points))
	count := float64(len(points))

	for i, point := range points {
		if !isFiniteXY(point) {
			continue
		}

		xIndex := ctx.XValueToIndex(point.X)
		if xIndex < 0 || xIndex >= width {
			continue
		}

		pointHeightFloat := ((point.Y - ctx.YMin) / (ctx.YMax - ctx.YMin)) * float64(height-1)
		pointHeight := int(pointHeightFloat)
		if pointHeight < 0 || pointHeight >= height {
			continue
		}

		screenX := x + xIndex
		screenY := y + height - 1 - pointHeight
		key := [2]int{screenX, screenY}

		intensity := 0.15 + (0.85 * (float64(i+1) / count))
		if intensity > intensityByCell[key] {
			intensityByCell[key] = intensity
		}
	}

	for key, intensity := range intensityByCell {
		trailBaseColor := o.resolvedTrailBaseColor(ctx)
		trailBg := lerpColor(ctx.Background, trailBaseColor, coreClamp01(intensity*o.opacity))
		xPos, yPos := key[0], key[1]
		mainc, currentCombc, currentStyle, _ := screen.GetContent(xPos, yPos)
		currentFg, _, attrs := currentStyle.Decompose()
		trailStyle := tcell.StyleDefault.Foreground(currentFg).Background(trailBg).Attributes(attrs)
		screen.SetContent(xPos, yPos, mainc, currentCombc, trailStyle)
	}
}

func (o *TrailOverlay) resolvedTrailBaseColor(ctx OverlayRenderContext) tcell.Color {
	base := o.color
	if !o.hasColor && len(ctx.SeriesColors) > 0 {
		base = ctx.SeriesColors[0]
	}
	return base
}
