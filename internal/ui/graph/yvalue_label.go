package graph

import (
	"math"
	"strconv"
)

// NewYAxisValueLabelOverlay builds a Y-axis label overlay from a single value function.
// Invalid values (NaN/Inf) are skipped for both y-positioning and text rendering.
func NewYAxisValueLabelOverlay(value func() float64, format func(float64) string) *YAxisLabelOverlay {
	if value == nil {
		panic("y-axis value label overlay requires a value function")
	}
	if format == nil {
		format = func(v float64) string {
			return strconv.FormatFloat(v, 'f', -1, 64)
		}
	}

	return YLabel(
		func() float64 {
			v := value()
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return math.NaN()
			}
			return v
		},
		func(_ OverlayRenderContext) string {
			v := value()
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return ""
			}
			return format(v)
		},
	)
}
