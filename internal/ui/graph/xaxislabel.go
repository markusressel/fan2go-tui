package graph

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type XAxisLabelOverlay struct {
	x         func() float64
	text      AxisLabelTextFunc
	textColor tcell.Color
	bgColor   *tcell.Color
}

func NewXAxisLabelOverlay(x func() float64, text AxisLabelTextFunc) *XAxisLabelOverlay {
	if x == nil {
		panic("x-axis label overlay requires an x function")
	}
	if text == nil {
		panic("x-axis label overlay requires a text function")
	}

	bgColor := tcell.ColorBlack
	return &XAxisLabelOverlay{
		x:         x,
		text:      text,
		textColor: tcell.ColorWhite,
		bgColor:   &bgColor,
	}
}

// XLabel is a concise alias for NewXAxisLabelOverlay.
func XLabel(x func() float64, text AxisLabelTextFunc) *XAxisLabelOverlay {
	return NewXAxisLabelOverlay(x, text)
}

func (o *XAxisLabelOverlay) WithX(x func() float64) *XAxisLabelOverlay {
	if x == nil {
		panic("x-axis label overlay requires an x function")
	}

	o.x = x
	return o
}

func (o *XAxisLabelOverlay) WithText(text string) *XAxisLabelOverlay {
	o.text = func(OverlayRenderContext) string {
		return text
	}
	return o
}

func (o *XAxisLabelOverlay) WithTextFunc(text func() string) *XAxisLabelOverlay {
	if text == nil {
		panic("x-axis label overlay requires a text function")
	}

	o.text = func(OverlayRenderContext) string {
		return text()
	}
	return o
}

func (o *XAxisLabelOverlay) WithTextFuncCtx(text AxisLabelTextFunc) *XAxisLabelOverlay {
	if text == nil {
		panic("x-axis label overlay requires a text function")
	}

	o.text = text
	return o
}

func (o *XAxisLabelOverlay) WithTextColor(color tcell.Color) *XAxisLabelOverlay {
	o.textColor = color
	return o
}

func (o *XAxisLabelOverlay) WithBackgroundColor(color tcell.Color) *XAxisLabelOverlay {
	o.bgColor = &color
	return o
}

func (o *XAxisLabelOverlay) draw(screen tcell.Screen, ctx OverlayRenderContext) {
	if o.x == nil || o.text == nil || ctx.Plot == nil || ctx.XValueToIndex == nil {
		return
	}

	label := o.text(ctx)
	if label == "" {
		return
	}

	xValue := o.x()
	if math.IsNaN(xValue) || math.IsInf(xValue, 0) {
		return
	}

	plotX, _, plotWidth, _ := ctx.Plot.GetPlotRect()
	innerX, innerY, innerWidth, innerHeight := ctx.Plot.GetInnerRect()
	if plotWidth <= 0 || innerWidth <= 0 || innerHeight <= 0 {
		return
	}

	xIndex := ctx.XValueToIndex(xValue)
	if xIndex < 0 || xIndex >= plotWidth {
		return
	}

	textWidth := len([]rune(label))
	if textWidth == 0 {
		return
	}

	screenX := plotX + xIndex
	startX := screenX - (textWidth / 2)
	screenY := innerY + innerHeight - 1
	if screenY < innerY || screenY >= innerY+innerHeight {
		return
	}

	bgColor := ctx.Background
	if o.bgColor != nil {
		bgColor = *o.bgColor
	}

	drawOverlayText(screen, label, startX, screenY, innerX, innerX+innerWidth, o.textColor, bgColor)
}
