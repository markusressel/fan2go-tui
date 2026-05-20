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

func newCurrentPwmXAxisLabelOverlay(getFan func() *client.Fan) *graph.XAxisLabelOverlay {
	if getFan == nil {
		getFan = func() *client.Fan { return nil }
	}

	return graph.XLabel(
		func() float64 {
			fan := getFan()
			if fan == nil {
				return math.NaN()
			}
			return float64(fan.Pwm)
		},
		func(_ graph.OverlayRenderContext) string {
			fan := getFan()
			if fan == nil {
				return ""
			}
			return strconv.Itoa(fan.Pwm)
		},
	).WithTextColor(theme.Colors.Graph.XAxisValueLabelText).
		WithBackgroundColor(theme.Colors.Graph.XAxisValueLabelBackground)
}
