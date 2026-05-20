package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"
	"math"
	"strconv"
)

func newCurrentCurveYAxisLabelOverlay(getCurve func() *client.Curve) *graph.YAxisLabelOverlay {
	if getCurve == nil {
		getCurve = func() *client.Curve { return nil }
	}

	return graph.NewYAxisValueLabelOverlay(
		func() float64 {
			curve := getCurve()
			if curve == nil {
				return math.NaN()
			}
			return curve.Value
		},
		func(v float64) string {
			return strconv.Itoa(int(v))
		},
	).
		WithTextColor(theme.Colors.Graph.YAxisValueLabelText).
		WithBackgroundColor(theme.Colors.Graph.YAxisValueLabelBackground)
}
