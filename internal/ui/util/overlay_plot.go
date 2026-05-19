package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

// OverlayPlot extends tvxwidgets.Plot with lightweight overlay rendering.
type OverlayPlot struct {
	*tvxwidgets.Plot

	verticalLineIndex *int
	verticalLineColor tcell.Color

	overlayPointXIndex *int
	overlayPointYValue *float64
	overlayPointYMin   float64
	overlayPointYMax   float64
	overlayPointColor  tcell.Color
}

func NewOverlayPlot() *OverlayPlot {
	return &OverlayPlot{
		Plot:              tvxwidgets.NewPlot(),
		verticalLineColor: tcell.ColorYellow,
		overlayPointColor: tcell.ColorYellow,
	}
}

func (p *OverlayPlot) SetVerticalLineIndex(index *int) {
	if index == nil {
		p.verticalLineIndex = nil
		return
	}

	idx := *index
	p.verticalLineIndex = &idx
}

func (p *OverlayPlot) SetVerticalLineColor(color tcell.Color) {
	p.verticalLineColor = color
}

func (p *OverlayPlot) SetOverlayPoint(index *int, value *float64) {
	if index == nil || value == nil {
		p.overlayPointXIndex = nil
		p.overlayPointYValue = nil
		return
	}

	idx := *index
	val := *value
	p.overlayPointXIndex = &idx
	p.overlayPointYValue = &val
}

func (p *OverlayPlot) SetOverlayPointColor(color tcell.Color) {
	p.overlayPointColor = color
}

func (p *OverlayPlot) SetOverlayPointYRange(min, max float64) {
	p.overlayPointYMin = min
	p.overlayPointYMax = max
}

func (p *OverlayPlot) Draw(screen tcell.Screen) {
	p.Plot.Draw(screen)

	if p.verticalLineIndex == nil {
		// continue below to still allow point-only overlays
	} else {
		x, y, width, height := p.Plot.GetPlotRect()
		idx := *p.verticalLineIndex
		if idx >= 0 && idx < width {
			lineStyle := tcell.StyleDefault.Background(p.GetBackgroundColor()).Foreground(p.verticalLineColor)
			screenX := x + idx
			for yPos := y; yPos < y+height; yPos++ {
				screen.SetContent(screenX, yPos, '|', nil, lineStyle)
			}
		}
	}

	x, y, width, height := p.Plot.GetPlotRect()
	if p.overlayPointXIndex == nil || p.overlayPointYValue == nil {
		return
	}

	if p.overlayPointYMax <= p.overlayPointYMin {
		return
	}

	pointXIndex := *p.overlayPointXIndex
	if pointXIndex < 0 || pointXIndex >= width {
		return
	}

	pointY := *p.overlayPointYValue
	pointHeight := int(((pointY - p.overlayPointYMin) / (p.overlayPointYMax - p.overlayPointYMin)) * float64(height-1))
	if pointHeight < 0 || pointHeight >= height {
		return
	}

	pointStyle := tcell.StyleDefault.Background(p.GetBackgroundColor()).Foreground(p.overlayPointColor)
	screenX := x + pointXIndex
	screenY := y + height - 1 - pointHeight
	if screenY < y || screenY >= y+height {
		return
	}

	screen.SetContent(screenX, screenY, 'o', nil, pointStyle)
}
