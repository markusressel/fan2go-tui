package graph

import (
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

type GraphComponentConfig[T any] struct {
	// PlotType defines the plot type to use
	PlotType tvxwidgets.PlotType

	// MarkerType defines the character style to use for plot markers (points)
	MarkerType tvxwidgets.Marker

	// DrawXAxisLabel determines if the x-axis label should be drawn
	DrawXAxisLabel bool

	// DrawXAxisLabel determines if the y-axis label should be drawn
	DrawYAxisLabel bool

	YAxisAutoScaleMin bool
	YAxisAutoScaleMax bool

	YAxisLabelDataType tvxwidgets.PlotYAxisLabelDataType

	// PlotColors is a list of colors to use for the plot lines
	PlotColors []tcell.Color

	// Reversed determines if the graph should be drawn from left tp right instead of right to left
	// Default: false
	Reversed bool

	// XMax is the maximum value for the x-axis
	XMax int

	// Overlays are rendered in declaration order, later overlays are drawn on top.
	Overlays []GraphComponentOverlay[T]
}

// NewGraphComponentConfig creates a new GraphComponentConfig with default values
func NewGraphComponentConfig[T any]() *GraphComponentConfig[T] {
	return &GraphComponentConfig[T]{
		PlotType:           tvxwidgets.PlotTypeLineChart,
		MarkerType:         tvxwidgets.PlotMarkerBraille,
		DrawXAxisLabel:     true,
		DrawYAxisLabel:     true,
		YAxisAutoScaleMin:  false,
		YAxisAutoScaleMax:  true,
		YAxisLabelDataType: tvxwidgets.PlotYAxisLabelDataFloat,
		PlotColors:         []tcell.Color{tcell.ColorWhite, tcell.ColorWhite, tcell.ColorWhite, tcell.ColorWhite, tcell.ColorWhite},
		Reversed:           false,
	}
}

// NewGraphComponentConfigFor infers the generic type from a sample pointer.
func NewGraphComponentConfigFor[T any](_ *T) *GraphComponentConfig[T] {
	return NewGraphComponentConfig[T]()
}

// WithReversedOrder sets the graph to be drawn from left to right instead of right to left
func (c *GraphComponentConfig[T]) WithReversedOrder() *GraphComponentConfig[T] {
	c.Reversed = true
	return c
}

// WithPlotColors sets the colors to use for the plot lines
func (c *GraphComponentConfig[T]) WithPlotColors(colors ...tcell.Color) *GraphComponentConfig[T] {
	c.PlotColors = colors
	return c
}

// WithPlotColorList sets the list of colors to use for the plot lines
func (c *GraphComponentConfig[T]) WithPlotColorList(colors []tcell.Color) *GraphComponentConfig[T] {
	c.PlotColors = colors
	return c
}

// WithXMax sets the maximum value for the x-axis
func (c *GraphComponentConfig[T]) WithXMax(xMax int) *GraphComponentConfig[T] {
	c.XMax = xMax
	return c
}

func (c *GraphComponentConfig[T]) WithDrawXAxisLabel(draw bool) *GraphComponentConfig[T] {
	c.DrawXAxisLabel = draw
	return c
}

func (c *GraphComponentConfig[T]) WithDrawYAxisLabel(draw bool) *GraphComponentConfig[T] {
	c.DrawYAxisLabel = draw
	return c
}

func (c *GraphComponentConfig[T]) WithYAxisAutoScaleMin(autoScale bool) *GraphComponentConfig[T] {
	c.YAxisAutoScaleMin = autoScale
	return c
}

func (c *GraphComponentConfig[T]) WithYAxisAutoScaleMax(autoScale bool) *GraphComponentConfig[T] {
	c.YAxisAutoScaleMax = autoScale
	return c
}

func (c *GraphComponentConfig[T]) WithYAxisLabelDataType(dataType tvxwidgets.PlotYAxisLabelDataType) *GraphComponentConfig[T] {
	c.YAxisLabelDataType = dataType
	return c
}

func (c *GraphComponentConfig[T]) WithOverlays(overlays ...GraphComponentOverlay[T]) *GraphComponentConfig[T] {
	c.Overlays = append(c.Overlays, overlays...)
	return c
}

// WithOverlay is a backward-compatible alias for WithOverlays.
func (c *GraphComponentConfig[T]) WithOverlay(overlays ...GraphComponentOverlay[T]) *GraphComponentConfig[T] {
	return c.WithOverlays(overlays...)
}
