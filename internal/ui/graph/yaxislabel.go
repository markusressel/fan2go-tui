package graph

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type YAxisLabelOverlay struct {
	y         func() float64
	text      AxisLabelTextFunc
	textColor tcell.Color
	bgColor   *tcell.Color
}

func NewYAxisLabelOverlay(y func() float64, text AxisLabelTextFunc) *YAxisLabelOverlay {
	if y == nil {
		panic("y-axis label overlay requires a y function")
	}
	if text == nil {
		panic("y-axis label overlay requires a text function")
	}

	bgColor := tcell.ColorBlack
	return &YAxisLabelOverlay{
		y:         y,
		text:      text,
		textColor: tcell.ColorWhite,
		bgColor:   &bgColor,
	}
}

// YLabel is a concise alias for NewYAxisLabelOverlay.
func YLabel(y func() float64, text AxisLabelTextFunc) *YAxisLabelOverlay {
	return NewYAxisLabelOverlay(y, text)
}

func (o *YAxisLabelOverlay) WithY(y func() float64) *YAxisLabelOverlay {
	if y == nil {
		panic("y-axis label overlay requires a y function")
	}

	o.y = y
	return o
}

func (o *YAxisLabelOverlay) WithText(text string) *YAxisLabelOverlay {
	o.text = func(OverlayRenderContext) string {
		return text
	}
	return o
}

func (o *YAxisLabelOverlay) WithTextFunc(text func() string) *YAxisLabelOverlay {
	if text == nil {
		panic("y-axis label overlay requires a text function")
	}

	o.text = func(OverlayRenderContext) string {
		return text()
	}
	return o
}

func (o *YAxisLabelOverlay) WithTextFuncCtx(text AxisLabelTextFunc) *YAxisLabelOverlay {
	if text == nil {
		panic("y-axis label overlay requires a text function")
	}

	o.text = text
	return o
}

func (o *YAxisLabelOverlay) WithTextColor(color tcell.Color) *YAxisLabelOverlay {
	o.textColor = color
	return o
}

func (o *YAxisLabelOverlay) WithBackgroundColor(color tcell.Color) *YAxisLabelOverlay {
	o.bgColor = &color
	return o
}

func (o *YAxisLabelOverlay) draw(screen tcell.Screen, ctx OverlayRenderContext) {
	if o.y == nil || o.text == nil || ctx.Plot == nil || ctx.YMax <= ctx.YMin {
		return
	}

	label := o.text(ctx)
	if label == "" {
		return
	}

	yValue := o.y()
	if math.IsNaN(yValue) || math.IsInf(yValue, 0) {
		return
	}

	plotX, plotY, _, plotHeight := ctx.Plot.GetPlotRect()
	innerX, innerY, innerWidth, innerHeight := ctx.Plot.GetInnerRect()
	if plotHeight <= 0 || innerWidth <= 0 || innerHeight <= 0 {
		return
	}

	pointHeightFloat := ((yValue - ctx.YMin) / (ctx.YMax - ctx.YMin)) * float64(plotHeight-1)
	pointHeight := int(pointHeightFloat)
	if pointHeight < 0 || pointHeight >= plotHeight {
		return
	}

	screenY := plotY + plotHeight - 1 - pointHeight
	if screenY < innerY || screenY >= innerY+innerHeight {
		return
	}

	textWidth := len([]rune(label))
	if textWidth == 0 {
		return
	}

	maxXExclusive := plotX - 2
	if maxXExclusive <= innerX {
		return
	}
	startX := maxXExclusive - textWidth

	bgColor := ctx.Background
	if o.bgColor != nil {
		bgColor = *o.bgColor
	}

	drawOverlayText(screen, label, startX, screenY, innerX, maxXExclusive, o.textColor, bgColor)
}
