package graph

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

func TestHeatmapOverlayDrawsHistoryPointsInsidePlot(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	overlay := NewHeatmapOverlay(func() []XY {
		return []XY{{X: 12, Y: 50}}
	}).WithColor(tcell.ColorRed)

	ctx := OverlayRenderContext{
		Plot:          plot,
		XValueToIndex: func(float64) int { return 8 },
		YMin:          0,
		YMax:          100,
		Background:    tcell.ColorBlack,
	}

	overlay.draw(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	pointHeight := int(((50.0 - 0.0) / (100.0 - 0.0)) * float64(plotHeight-1))
	expectedX := plotX + 8
	expectedY := plotY + plotHeight - 1 - pointHeight

	_, _, style, _ := screen.GetContent(expectedX, expectedY)
	fg, bg, _ := style.Decompose()
	if fg == tcell.ColorBlack && bg == tcell.ColorBlack {
		t.Fatalf("expected visible heatmap styling at (%d,%d)", expectedX, expectedY)
	}

	_, _, style, _ = screen.GetContent(expectedX, expectedY)
	_, bg, _ = style.Decompose()
	if bg == tcell.ColorBlack {
		t.Fatalf("expected heatmap background tint at (%d,%d)", expectedX, expectedY)
	}
}
