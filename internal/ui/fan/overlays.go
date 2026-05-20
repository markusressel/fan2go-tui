package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"
	"math"
	"strconv"
)

func newCurrentRpmYAxisLabelOverlay(getFan func() *client.Fan) *graph.YAxisLabelOverlay {
	if getFan == nil {
		getFan = func() *client.Fan { return nil }
	}

	return graph.NewYAxisValueLabelOverlay(
		func() float64 {
			fan := getFan()
			if fan == nil {
				return math.NaN()
			}
			return float64(fan.Rpm)
		},
		func(v float64) string {
			return strconv.Itoa(int(v))
		},
	).WithTextColor(theme.Colors.Graph.YAxisValueLabelText).
		WithBackgroundColor(theme.Colors.Graph.YAxisValueLabelBackground)
}
