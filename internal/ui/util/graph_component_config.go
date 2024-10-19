package util

import (
	"github.com/gdamore/tcell/v2"
	"github.com/navidys/tvxwidgets"
)

type GraphComponentConfig struct {
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
}

// NewGraphComponentConfig creates a new GraphComponentConfig with default values
func NewGraphComponentConfig() *GraphComponentConfig {
	return &GraphComponentConfig{
		PlotType:           tvxwidgets.PlotTypeLineChart,
		MarkerType:         tvxwidgets.PlotMarkerBraille,
		DrawXAxisLabel:     true,
		DrawYAxisLabel:     true,
		YAxisAutoScaleMin:  false,
		YAxisAutoScaleMax:  true,
		YAxisLabelDataType: tvxwidgets.PlotYAxisLabelDataFloat,
		PlotColors:         make([]tcell.Color, 0),
		Reversed:           false,
	}
}

// WithReversedOrder sets the graph to be drawn from left to right instead of right to left
func (c *GraphComponentConfig) WithReversedOrder() *GraphComponentConfig {
	c.Reversed = true
	return c
}

// WithPlotColors sets the colors to use for the plot lines
func (c *GraphComponentConfig) WithPlotColors(colors ...tcell.Color) *GraphComponentConfig {
	c.PlotColors = colors
	return c
}

// WithPlotColorList sets the list of colors to use for the plot lines
func (c *GraphComponentConfig) WithPlotColorList(colors []tcell.Color) *GraphComponentConfig {
	c.PlotColors = colors
	return c
}

// WithXMax sets the maximum value for the x-axis
func (c *GraphComponentConfig) WithXMax(xMax int) *GraphComponentConfig {
	c.XMax = xMax
	return c
}

func (c *GraphComponentConfig) WithDrawXAxisLabel(draw bool) *GraphComponentConfig {
	c.DrawXAxisLabel = draw
	return c
}

func (c *GraphComponentConfig) WithDrawYAxisLabel(draw bool) *GraphComponentConfig {
	c.DrawYAxisLabel = draw
	return c
}

func (c *GraphComponentConfig) WithYAxisAutoScaleMin(autoScale bool) *GraphComponentConfig {
	c.YAxisAutoScaleMin = autoScale
	return c
}

func (c *GraphComponentConfig) WithYAxisAutoScaleMax(autoScale bool) *GraphComponentConfig {
	c.YAxisAutoScaleMax = autoScale
	return c
}

func (c *GraphComponentConfig) WithYAxisLabelDataType(dataType tvxwidgets.PlotYAxisLabelDataType) *GraphComponentConfig {
	c.YAxisLabelDataType = dataType
	return c
}
