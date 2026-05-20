package graph

import (
	"math"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

func expectedYLabelStartX(plot *tvxwidgets.Plot, yValue float64, text string, ctx OverlayRenderContext, manualOffset int) int {
	plotX, _, _, _ := plot.GetPlotRect()

	countDigits := func(v float64) int {
		value := math.Abs(math.Trunc(v))
		if value < 1 {
			return 1
		}
		digits := 0
		for value >= 1 {
			value /= 10
			digits++
		}
		return digits
	}

	maxDigits := countDigits(ctx.YMax)
	if minDigits := countDigits(ctx.YMin); minDigits > maxDigits {
		maxDigits = minDigits
	}
	currentDigits := countDigits(yValue)
	dynamicOffset := maxDigits - currentDigits
	if dynamicOffset < 0 {
		dynamicOffset = 0
	}

	constantOffset := 1
	if ctx.YAxisLabelsAreInts {
		constantOffset = 4
	}

	maxXExclusive := plotX - (constantOffset + dynamicOffset) + manualOffset
	return maxXExclusive - len([]rune(text))
}

func TestXAxisLabelOverlayDrawsOnAxisRowOutsidePlotRect(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	overlay := NewXAxisLabelOverlay(
		func() float64 { return 12 },
		func(OverlayRenderContext) string { return "PWM" },
	).WithTextColor(tcell.ColorRed).WithBackgroundColor(tcell.ColorBlue)

	ctx := OverlayRenderContext{
		Plot:          plot,
		XValueToIndex: func(float64) int { return 8 },
		YMin:          0,
		YMax:          100,
		Background:    tcell.ColorBlack,
	}

	overlay.draw(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	innerX, innerY, _, innerHeight := plot.GetInnerRect()
	expectedY := innerY + innerHeight - 1
	if expectedY < plotY+plotHeight {
		t.Fatalf("expected x-axis label row to be outside plot area: expectedY=%d plotBottom=%d", expectedY, plotY+plotHeight)
	}

	expectedStartX := plotX + 8 - 1 // center of 3-rune text
	mainc, _, style, _ := screen.GetContent(expectedStartX+1, expectedY)
	if mainc != 'W' {
		t.Fatalf("expected center rune W at (%d,%d), got %q", expectedStartX+1, expectedY, mainc)
	}

	fg, bg, _ := style.Decompose()
	if fg != tcell.ColorRed {
		t.Fatalf("expected foreground %v, got %v", tcell.ColorRed, fg)
	}
	if bg != tcell.ColorBlue {
		t.Fatalf("expected background %v, got %v", tcell.ColorBlue, bg)
	}

	if expectedStartX < innerX {
		t.Fatalf("unexpected test setup: text starts before inner rect")
	}
}

func TestYAxisLabelOverlayDrawsLeftOfPlotRect(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	overlay := NewYAxisLabelOverlay(
		func() float64 { return 50 },
		func(OverlayRenderContext) string { return "RPM" },
	).WithTextColor(tcell.ColorGreen).WithBackgroundColor(tcell.ColorMaroon)

	ctx := OverlayRenderContext{
		Plot:       plot,
		YMin:       0,
		YMax:       100,
		Background: tcell.ColorBlack,
	}

	overlay.draw(screen, ctx)

	plotX, plotY, _, plotHeight := plot.GetPlotRect()
	innerX, _, _, _ := plot.GetInnerRect()
	pointHeight := int(((50.0 - 0.0) / (100.0 - 0.0)) * float64(plotHeight-1))
	expectedY := plotY + plotHeight - 1 - pointHeight
	expectedStartX := expectedYLabelStartX(plot, 50, "RPM", ctx, 0)
	if expectedStartX+len([]rune("RPM")) > plotX-2 {
		t.Fatalf("expected y-axis label to be left of plot area")
	}
	if expectedStartX < innerX {
		t.Fatalf("unexpected test setup: text starts before inner rect")
	}

	mainc, _, style, _ := screen.GetContent(expectedStartX+1, expectedY)
	if mainc != 'P' {
		t.Fatalf("expected center rune P at (%d,%d), got %q", expectedStartX+1, expectedY, mainc)
	}

	fg, bg, _ := style.Decompose()
	if fg != tcell.ColorGreen {
		t.Fatalf("expected foreground %v, got %v", tcell.ColorGreen, fg)
	}
	if bg != tcell.ColorMaroon {
		t.Fatalf("expected background %v, got %v", tcell.ColorMaroon, bg)
	}
}

func TestAxisLabelCallbacksCanUseRenderContext(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 100)

	xLabel := NewXAxisLabelOverlay(
		func() float64 { return 12 },
		func(ctx OverlayRenderContext) string {
			if ctx.YMax == 100 {
				return "OK"
			}
			return "NO"
		},
	)

	yLabel := NewYAxisLabelOverlay(
		func() float64 { return 50 },
		func(ctx OverlayRenderContext) string {
			if ctx.YMin == 0 {
				return "YY"
			}
			return "NN"
		},
	)

	ctx := OverlayRenderContext{
		Plot:          plot,
		XValueToIndex: func(float64) int { return 8 },
		YMin:          0,
		YMax:          100,
		Background:    tcell.ColorBlack,
	}

	xLabel.draw(screen, ctx)
	yLabel.draw(screen, ctx)

	plotX, _, _, plotHeight := plot.GetPlotRect()
	innerY, innerHeight := func() (int, int) {
		_, y, _, h := plot.GetInnerRect()
		return y, h
	}()

	xMain, _, _, _ := screen.GetContent(plotX+8, innerY+innerHeight-1)
	if xMain != 'K' {
		t.Fatalf("expected context-based x label content, got %q", xMain)
	}

	pointHeight := int(((50.0 - 0.0) / (100.0 - 0.0)) * float64(plotHeight-1))
	yScreen := (func() int { _, py, _, _ := plot.GetPlotRect(); return py + plotHeight - 1 - pointHeight })()
	yStart := expectedYLabelStartX(plot, 50, "YY", ctx, 0)
	yMain, _, _, _ := screen.GetContent(yStart+1, yScreen)
	if yMain != 'Y' {
		t.Fatalf("expected context-based y label content, got %q", yMain)
	}
}

func TestYAxisLabelOverlayAutoAdjustsForIntegerTicks(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("init simulation screen: %v", err)
	}
	defer screen.Fini()
	screen.SetSize(80, 24)

	plot := tvxwidgets.NewPlot()
	plot.SetRect(0, 0, 40, 10)
	plot.SetYRange(0, 255)

	overlay := NewYAxisLabelOverlay(
		func() float64 { return 128 },
		func(OverlayRenderContext) string { return "128" },
	)

	ctx := OverlayRenderContext{
		Plot:               plot,
		YMin:               0,
		YMax:               255,
		Background:         tcell.ColorBlack,
		YAxisLabelsAreInts: true,
	}

	overlay.draw(screen, ctx)

	_, plotY, _, plotHeight := plot.GetPlotRect()
	pointHeight := int(((128.0 - 0.0) / (255.0 - 0.0)) * float64(plotHeight-1))
	expectedY := plotY + plotHeight - 1 - pointHeight
	expectedStartX := expectedYLabelStartX(plot, 128, "128", ctx, 0)

	mainc, _, _, _ := screen.GetContent(expectedStartX+1, expectedY)
	if mainc != '2' {
		t.Fatalf("expected auto-aligned integer y label at (%d,%d), got %q", expectedStartX+1, expectedY, mainc)
	}

	floatCtx := ctx
	floatCtx.YAxisLabelsAreInts = false
	floatStartX := expectedYLabelStartX(plot, 128, "128", floatCtx, 0)
	if expectedStartX != floatStartX-3 {
		t.Fatalf("expected integer labels to shift three columns left vs float labels, got floatStart=%d intStart=%d", floatStartX, expectedStartX)
	}
}
