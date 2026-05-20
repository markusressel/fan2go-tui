package graph

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

func TestTrailOverlayRecencyWinsPerCell(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	points := []XY{
		{X: 10, Y: 40}, // older
		{X: 10, Y: 40}, // newer (same cell)
	}
	overlay := NewTrailOverlay(func() []XY { return points }).WithColor(tcell.ColorRed)

	ctx := OverlayRenderContext{
		Plot:          plot,
		XValueToIndex: func(float64) int { return 8 },
		YMin:          0,
		YMax:          100,
		Background:    tcell.ColorBlack,
	}

	overlay.draw(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	pointHeight := int(((40.0 - 0.0) / (100.0 - 0.0)) * float64(plotHeight-1))
	targetX := plotX + 8
	targetY := plotY + plotHeight - 1 - pointHeight

	_, _, style, _ := screen.GetContent(targetX, targetY)
	_, bg, _ := style.Decompose()
	if bg == tcell.ColorBlack {
		t.Fatalf("expected trail background tint at (%d,%d)", targetX, targetY)
	}
}

func TestTrailOverlayInheritsSeriesColorByDefault(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	overlay := NewTrailOverlay(func() []XY {
		return []XY{{X: 10, Y: 40}}
	})

	seriesColor := tcell.NewRGBColor(0, 120, 255)
	ctx := OverlayRenderContext{
		Plot:          plot,
		XValueToIndex: func(float64) int { return 8 },
		YMin:          0,
		YMax:          100,
		Background:    tcell.ColorBlack,
		SeriesColors:  []tcell.Color{seriesColor},
	}

	overlay.draw(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	pointHeight := int(((40.0 - 0.0) / (100.0 - 0.0)) * float64(plotHeight-1))
	targetX := plotX + 8
	targetY := plotY + plotHeight - 1 - pointHeight

	_, _, style, _ := screen.GetContent(targetX, targetY)
	_, bg, _ := style.Decompose()

	expectedBase := lerpColor(seriesColor, tcell.ColorBlack, 0.0)
	expected := lerpColor(tcell.ColorBlack, expectedBase, 0.6)
	if bg == tcell.ColorBlack || bg == tcell.ColorOrange {
		t.Fatalf("expected inherited series-based trail tint, got %v", bg)
	}
	if bg != expected {
		t.Fatalf("expected default inherited trail tint %v, got %v", expected, bg)
	}
}

func TestTrailOverlayWithOpacityControlsTintStrength(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	overlay := NewTrailOverlay(func() []XY {
		return []XY{{X: 10, Y: 40}}
	}).WithOpacity(0.3)

	seriesColor := tcell.NewRGBColor(0, 120, 255)
	ctx := OverlayRenderContext{
		Plot:          plot,
		XValueToIndex: func(float64) int { return 8 },
		YMin:          0,
		YMax:          100,
		Background:    tcell.ColorBlack,
		SeriesColors:  []tcell.Color{seriesColor},
	}

	overlay.draw(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	pointHeight := int(((40.0 - 0.0) / (100.0 - 0.0)) * float64(plotHeight-1))
	targetX := plotX + 8
	targetY := plotY + plotHeight - 1 - pointHeight

	_, _, style, _ := screen.GetContent(targetX, targetY)
	_, bg, _ := style.Decompose()

	expected := lerpColor(tcell.ColorBlack, seriesColor, 0.3)
	if bg != expected {
		t.Fatalf("expected opacity-based inherited trail tint %v, got %v", expected, bg)
	}
}
