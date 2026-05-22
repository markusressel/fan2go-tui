package graph

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestMapValueToBrailleSubRowClampsToBounds(t *testing.T) {
	totalSubRows := 16

	top, ok := mapValueToBrailleSubRow(100, 0, 10, totalSubRows)
	if !ok {
		t.Fatalf("expected mapping to succeed")
	}
	if top != 0 {
		t.Fatalf("expected values above max to clamp to top sub-row, got %d", top)
	}

	bottom, ok := mapValueToBrailleSubRow(-100, 0, 10, totalSubRows)
	if !ok {
		t.Fatalf("expected mapping to succeed")
	}
	if bottom != totalSubRows-1 {
		t.Fatalf("expected values below min to clamp to bottom sub-row, got %d", bottom)
	}
}

func TestDrawLineSeriesUsesBrailleSubRowsForSinglePoint(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := NewOverlayPlot()
	plot.SetRect(0, 0, 40, 10)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	if plotHeight <= 0 {
		t.Fatalf("expected positive plot height")
	}

	totalSubRows := plotHeight * 4
	value := 0.0
	expectedSubY := -1
	for i := 1; i < 1000; i++ {
		candidate := float64(i) / 1000.0
		subY, ok := mapValueToBrailleSubRow(candidate, 0, 1, totalSubRows)
		if ok && subY%4 != 0 {
			value = candidate
			expectedSubY = subY
			break
		}
	}
	if expectedSubY < 0 {
		t.Fatalf("failed to find a value that maps to a non-cell-aligned braille sub-row")
	}

	ctx := OverlayRenderContext{
		YMin:            0,
		YMax:            1,
		Background:      tcell.ColorBlack,
		SeriesData:      [][]float64{{value}},
		SeriesColors:    []tcell.Color{tcell.ColorGreen},
		ValueBufferSize: 1,
	}

	plot.drawLineSeries(screen, ctx)

	screenY := plotY + (expectedSubY / 4)
	mainc, _, _, _ := screen.GetContent(plotX, screenY)
	if mainc == ' ' {
		t.Fatalf("expected braille point at (%d,%d)", plotX, screenY)
	}

	bits := mainc - rune(0x2800)
	expectedBit := brailleLineBit[expectedSubY%4][0]
	if bits&expectedBit == 0 {
		t.Fatalf("expected braille rune %U to include bit %U for sub-row %d", mainc, expectedBit, expectedSubY%4)
	}
}
