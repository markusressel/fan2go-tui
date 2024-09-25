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

	// PlotColors is a list of colors to use for the plot lines
	PlotColors []tcell.Color

	// Reversed determines if the graph should be drawn from left tp right instead of right to left
	// Default: false
	Reversed bool
}

// NewGraphComponentConfig creates a new GraphComponentConfig with default values
func NewGraphComponentConfig() *GraphComponentConfig {
	return &GraphComponentConfig{
		PlotType:   tvxwidgets.PlotTypeLineChart,
		MarkerType: tvxwidgets.PlotMarkerBraille,
		PlotColors: make([]tcell.Color, 0),
		Reversed:   false,
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
