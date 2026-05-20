package graph

import (
	"github.com/gdamore/tcell/v2"
)

type HeatmapOverlay struct {
	points    func() []XY
	color     tcell.Color
	maxPoints int
}

func NewHeatmapOverlay(points func() []XY) *HeatmapOverlay {
	if points == nil {
		panic("heatmap overlay requires a points function")
	}

	return &HeatmapOverlay{
		points:    points,
		color:     tcell.ColorOrange,
		maxPoints: 200,
	}
}

// Heatmap is a concise alias for NewHeatmapOverlay.
func Heatmap(points func() []XY) *HeatmapOverlay {
	return NewHeatmapOverlay(points)
}

func (o *HeatmapOverlay) WithPoints(points func() []XY) *HeatmapOverlay {
	if points == nil {
		panic("heatmap overlay requires a points function")
	}
	o.points = points
	return o
}

func (o *HeatmapOverlay) WithColor(color tcell.Color) *HeatmapOverlay {
	o.color = color
	return o
}

func (o *HeatmapOverlay) WithMaxPoints(maxPoints int) *HeatmapOverlay {
	o.maxPoints = maxPoints
	return o
}

func (o *HeatmapOverlay) draw(screen tcell.Screen, ctx OverlayRenderContext) {
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

	weights := make(map[[2]int]float64, len(points))
	maxWeight := 0.0
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

		recencyWeight := 0.2 + (0.8 * (float64(i+1) / count))
		weights[key] += recencyWeight
		if weights[key] > maxWeight {
			maxWeight = weights[key]
		}
	}

	if maxWeight <= 0 {
		return
	}

	for key, weight := range weights {
		intensity := coreClamp01(weight / maxWeight)
		heatBg := lerpColor(ctx.Background, o.color, intensity)

		xPos, yPos := key[0], key[1]
		mainc, currentCombc, currentStyle, _ := screen.GetContent(xPos, yPos)
		currentFg, _, attrs := currentStyle.Decompose()
		heatStyle := tcell.StyleDefault.Foreground(currentFg).Background(heatBg).Attributes(attrs)
		screen.SetContent(xPos, yPos, mainc, currentCombc, heatStyle)
	}
}

func coreClamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
