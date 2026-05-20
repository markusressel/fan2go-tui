package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/graph"
	"fan2go-tui/internal/ui/theme"
	"fmt"
	"math"
)

func newCurrentSensorYAxisLabelOverlay(getSensor func() *client.Sensor) *graph.YAxisLabelOverlay {
	if getSensor == nil {
		getSensor = func() *client.Sensor { return nil }
	}

	return graph.NewYAxisValueLabelOverlay(
		func() float64 {
			sensor := getSensor()
			if sensor == nil {
				return math.NaN()
			}
			return sensor.MovingAvg / 1000
		},
		func(v float64) string {
			return fmt.Sprintf("%.2f", v)
		},
	).
		WithTextColor(theme.Colors.Graph.YAxisValueLabelText).
		WithBackgroundColor(theme.Colors.Graph.YAxisValueLabelBackground)
}
