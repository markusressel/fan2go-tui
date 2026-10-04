package graph

import (
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
	"github.com/rivo/tview"
)

func TestSeriesValueProvider(t *testing.T) {
	values := map[int]float64{
		10: 100.0,
		20: 200.0,
		30: 300.0,
	}

	p := NewDiscreteIntSeriesValueProvider(values)
	if p.X(0) != 10.0 {
		t.Errorf("expected X(0)=10, got %v", p.X(0))
	}
	if p.F(10.0) != 100.0 {
		t.Errorf("expected F(10)=100, got %v", p.F(10.0))
	}
	if p.XLabel(1, 20.0) != "20" {
		t.Errorf("expected XLabel(1, 20)=20, got %q", p.XLabel(1, 20.0))
	}

	// Update values
	p.SetValues(map[int]float64{5: 50.0})
	if p.X(0) != 5.0 {
		t.Errorf("expected updated X(0)=5, got %v", p.X(0))
	}

	// Nil / empty coverage
	var nilProvider *DiscreteIntSeriesValueProvider
	nilProvider.SetValues(nil)

	// Rounded slice provider
	sliceData := []float64{10.4, 20.6, 30.0}
	sliceP := NewRoundedSliceSeriesValueProvider(&sliceData)
	if sliceP.X(1) != 1.0 {
		t.Errorf("expected slice X(1)=1, got %v", sliceP.X(1))
	}
	if sliceP.F(1.0) != 20.6 {
		t.Errorf("expected slice F(1)=20.6, got %v", sliceP.F(1.0))
	}
	if sliceP.XLabel(0, 0.0) != "0" {
		t.Errorf("expected slice XLabel(0, 0)=0, got %q", sliceP.XLabel(0, 0.0))
	}

	// Line from provider
	line := NewGraphLineFromSeriesValueProvider("TestLine", p)
	if line.GetName() != "TestLine" {
		t.Errorf("expected title TestLine, got %s", line.GetName())
	}
	line.SetColor(tcell.ColorRed)
	if line.GetColor() != tcell.ColorRed {
		t.Errorf("expected color red, got %v", line.GetColor())
	}

	// Bar from provider
	bar := NewGraphBarFromSeriesValueProvider("TestBar", p)
	if bar == nil {
		t.Fatalf("expected non-nil bar")
	}
	bar.SetColor(tcell.ColorBlue)
	if bar.GetColor() != tcell.ColorBlue {
		t.Errorf("expected color blue, got %v", bar.GetColor())
	}
}

func TestGraphSeriesLegend(t *testing.T) {
	legend := NewGraphSeriesLegend("RPM").
		WithUnit("RPM").
		WithTextColor(tcell.ColorWhite).
		WithBackgroundColor(tcell.ColorBlack).
		WithGlyph('~')

	display := legend.displayText()
	if display != "RPM (RPM)" {
		t.Errorf("expected 'RPM (RPM)', got %q", display)
	}

	legendNoUnit := NewGraphSeriesLegend("FanSpeed")
	if legendNoUnit.displayText() != "FanSpeed" {
		t.Errorf("expected 'FanSpeed', got %q", legendNoUnit.displayText())
	}

	var nilLegend *GraphSeriesLegend
	if nilLegend.displayText() != "" {
		t.Errorf("expected empty string for nil legend, got %q", nilLegend.displayText())
	}
}

func TestVLineAndMarkerOverlay(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	defer screen.Fini()
	screen.SetSize(60, 20)

	tvxPlot := tvxwidgets.NewPlot()
	tvxPlot.SetRect(0, 0, 40, 15)

	ctx := OverlayRenderContext{
		Plot:               tvxPlot,
		XValueToIndex:      func(x float64) int { return int(x) },
		XValueToIndexFloat: func(x float64) float64 { return x },
		YMin:               0,
		YMax:               100,
		Background:         tcell.ColorBlack,
	}

	// VLine
	vline := NewVerticalLine(func() float64 { return 10.0 }).
		WithX(func() float64 { return 10.0 }).
		WithColor(tcell.ColorGreen).
		WithRune('|')
	vline.draw(screen, ctx)

	// VLine helper
	vline2 := VLine(func() float64 { return 5.0 })
	vline2.draw(screen, ctx)

	// Marker
	marker := NewMarkerOverlay(func() XY { return XY{X: 10.0, Y: 50.0} }).
		WithCoord(func() XY { return XY{X: 10.0, Y: 50.0} }).
		WithColor(tcell.ColorYellow).
		WithRune('*')
	marker.draw(screen, ctx)

	// Marker helper
	marker2 := Marker(func() XY { return XY{X: 5.0, Y: 25.0} })
	marker2.draw(screen, ctx)

	// YAxisValueLabelOverlay
	yValOverlay := NewYAxisValueLabelOverlay(func() float64 { return 75.0 }, nil)
	if yValOverlay == nil {
		t.Fatalf("expected non-nil yAxisValueLabelOverlay")
	}
}

func TestGraphLine_RangesAndShifts(t *testing.T) {
	line := NewGraphLine("Line1", nil, nil, nil)
	line.SetXRange(0, 100)
	if line.GetXMax() == nil || *line.GetXMax() != 100 {
		t.Errorf("unexpected XMax")
	}
	line.ResetXRange()

	line.SetYRange(10, 50)
	if line.GetYMax() == nil || *line.GetYMax() != 50 {
		t.Errorf("unexpected YMax")
	}
	line.ResetYRange()

	if line.GetYAxisZoomFactor() != 1.0 {
		t.Errorf("expected zoom factor 1.0")
	}
	if line.GetYAxisShift() != 0.0 {
		t.Errorf("expected shift 0.0")
	}
}

func TestOverlayPlot_LayoutAndOverlays(t *testing.T) {
	plot := NewOverlayPlot()
	plot.SetRect(0, 0, 40, 20)
	plot.SetOverlays([]GraphComponentOverlay{
		VLine(func() float64 { return 10 }),
	})
	plot.SetOverlayContext(OverlayRenderContext{
		XValueToIndex: func(x float64) int { return int(x) },
		YMin:          0,
		YMax:          100,
	})
	layoutCalled := false
	plot.SetOnLayoutChange(func() { layoutCalled = true })

	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	defer screen.Fini()
	screen.SetSize(40, 20)
	plot.Draw(screen)

	_ = layoutCalled
}

func TestGraphComponent(t *testing.T) {
	app := tview.NewApplication()
	cfg := NewGraphComponentConfig().
		WithReversedOrder()

	comp := NewGraphComponent(app, cfg)
	comp.SetTitle("System Graph")

	values := map[int]float64{0: 10, 1: 20, 2: 30}
	p := NewDiscreteIntSeriesValueProvider(values)
	line := NewGraphLineFromSeriesValueProvider("Series 1", p)
	legend := NewGraphSeriesLegend("Series 1").WithUnit("%")

	comp.AddSeriesWithLegend(line, legend)

	bar := NewGraphBarFromSeriesValueProvider("Series 2", p)
	comp.AddSeries(bar)

	// Set ranges
	minY := 0.0
	maxY := 100.0
	comp.SetYMinValue(&minY)
	comp.SetYMaxValue(&maxY)
	comp.SetYRange(0, 100)

	// Set Layout rect so it's visible
	comp.GetLayout().SetRect(0, 0, 80, 24)
	comp.ZoomToRangeX(0, 50)
	comp.Refresh()

	if comp.GetLayout() == nil {
		t.Errorf("expected non-nil layout")
	}

	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	defer screen.Fini()
	screen.SetSize(80, 24)
	comp.GetLayout().Draw(screen)
}
