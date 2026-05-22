package graph

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestBarFillBitsUsesSingleBrailleColumn(t *testing.T) {
	left := barFillBits(4, 0)
	right := barFillBits(4, 1)

	if left == 0 || right == 0 {
		t.Fatalf("expected non-empty braille bits for full-height bar")
	}
	if left == right {
		t.Fatalf("expected left and right column masks to differ")
	}
	if left&right != 0 {
		t.Fatalf("expected left and right masks to not overlap, got intersection %U", left&right)
	}
}

func TestDrawBarsRendersSingleColumnByDefault(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := NewOverlayPlot()
	plot.SetRect(0, 0, 40, 10)

	bar := NewGraphBar("bar", func(i int) float64 { return float64(i) }, func(float64) float64 { return 1.0 }, func(int, float64) string { return "" })
	bar.SetColor(tcell.ColorGreen)

	ctx := OverlayRenderContext{
		YMin:            0,
		YMax:            1,
		Background:      tcell.ColorBlack,
		Bars:            []*GraphBar{bar},
		ValueBufferSize: 1,
		XValueToIndex:   func(float64) int { return 0 },
	}

	plot.drawBars(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	mainc, _, _, _ := screen.GetContent(plotX, plotY+plotHeight-1)
	if mainc == ' ' {
		t.Fatalf("expected bottom bar cell to be rendered")
	}

	bits := mainc - rune(0x2800)
	if bits&barFillBits(1, 1) == 0 {
		t.Fatalf("expected default bar rendering to use right braille column bit")
	}
	if bits&barFillBits(1, 0) != 0 {
		t.Fatalf("expected default bar rendering to avoid left braille column bit")
	}
}

func TestDrawBarsUsesFractionalSubColumnFromXMapping(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := NewOverlayPlot()
	plot.SetRect(0, 0, 40, 10)

	bar := NewGraphBar("bar", func(i int) float64 { return float64(i) }, func(float64) float64 { return 1.0 }, func(int, float64) string { return "" })
	bar.SetColor(tcell.ColorGreen)

	ctx := OverlayRenderContext{
		YMin:            0,
		YMax:            1,
		Background:      tcell.ColorBlack,
		Bars:            []*GraphBar{bar},
		ValueBufferSize: 1,
		XValueToIndex:   func(float64) int { return 0 },
		XValueToIndexFloat: func(float64) float64 {
			return 0.0 // maps to left sub-column
		},
	}

	plot.drawBars(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	mainc, _, _, _ := screen.GetContent(plotX, plotY+plotHeight-1)
	if mainc == ' ' {
		t.Fatalf("expected bottom bar cell to be rendered")
	}

	bits := mainc - rune(0x2800)
	if bits&barFillBits(1, 0) == 0 {
		t.Fatalf("expected fractional x mapping to select left braille column")
	}
	if bits&barFillBits(1, 1) != 0 {
		t.Fatalf("expected fractional x mapping to avoid right braille column")
	}
}

func TestDrawBarsInterpolatesBetweenAdjacentSamples(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := NewOverlayPlot()
	plot.SetRect(0, 0, 40, 10)

	bar := NewGraphBar(
		"bar",
		func(i int) float64 { return float64(i) },
		func(float64) float64 { return 1.0 },
		func(int, float64) string { return "" },
	)
	bar.SetColor(tcell.ColorGreen)

	ctx := OverlayRenderContext{
		YMin:            0,
		YMax:            1,
		Background:      tcell.ColorBlack,
		Bars:            []*GraphBar{bar},
		ValueBufferSize: 2,
		XValueToIndex:   func(x float64) int { return int(x) },
	}

	plot.drawBars(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	rightCellRune, _, _, _ := screen.GetContent(plotX+1, plotY+plotHeight-1)
	if rightCellRune == ' ' {
		t.Fatalf("expected interpolated bar rendering in second cell")
	}

	bits := rightCellRune - rune(0x2800)
	if bits&barFillBits(1, 0) == 0 || bits&barFillBits(1, 1) == 0 {
		t.Fatalf("expected interpolation to fill both sub-columns in the second cell, got bits %U", bits)
	}
}
