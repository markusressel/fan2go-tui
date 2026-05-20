package graph

import (
	"math"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

func TestHorizontalLineDrawsAcrossPlot(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	overlay := NewHorizontalLine(func() float64 { return 50 }).WithColor(tcell.ColorAqua)
	ctx := OverlayRenderContext{
		Plot:       plot,
		YMin:       0,
		YMax:       100,
		Background: tcell.ColorBlack,
	}

	overlay.draw(screen, ctx)

	plotX, plotY, plotWidth, plotHeight := plot.GetPlotRect()
	pointHeight := int(((50.0 - 0.0) / (100.0 - 0.0)) * float64(plotHeight-1))
	expectedY := plotY + plotHeight - 1 - pointHeight

	mainc, _, style, _ := screen.GetContent(plotX+plotWidth/2, expectedY)
	if mainc == ' ' {
		t.Fatalf("expected horizontal line glyph at y=%d", expectedY)
	}
	fg, _, _ := style.Decompose()
	if fg != tcell.ColorAqua {
		t.Fatalf("expected foreground %v, got %v", tcell.ColorAqua, fg)
	}
}

func TestHorizontalLineSkipsInvalidY(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	overlay := NewHorizontalLine(func() float64 { return math.NaN() })
	ctx := OverlayRenderContext{Plot: plot, YMin: 0, YMax: 100, Background: tcell.ColorBlack}
	overlay.draw(screen, ctx)

	plotX, plotY, plotWidth, plotHeight := plot.GetPlotRect()
	mainc, _, _, _ := screen.GetContent(plotX+plotWidth/2, plotY+plotHeight/2)
	if mainc != ' ' {
		t.Fatalf("expected no draw for invalid y, got %q", mainc)
	}
}
