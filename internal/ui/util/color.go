package util

import (
	coreutil "fan2go-tui/internal/util"
	"math"

	"github.com/gdamore/tcell/v2"
)

func LerpColor(a, b [3]int32, t float64) tcell.Color {
	t = coreutil.Clamp01(t)
	r := int32(math.Round(float64(a[0]) + (float64(b[0]-a[0]) * t)))
	g := int32(math.Round(float64(a[1]) + (float64(b[1]-a[1]) * t)))
	bv := int32(math.Round(float64(a[2]) + (float64(b[2]-a[2]) * t)))
	return tcell.NewRGBColor(r, g, bv)
}

func GradientColorAt(t float64, blue, green, yellow, red [3]int32) tcell.Color {
	t = coreutil.Clamp01(t)

	switch {
	case t <= 1.0/3.0:
		return LerpColor(blue, green, t*3.0)
	case t <= 2.0/3.0:
		return LerpColor(green, yellow, (t-1.0/3.0)*3.0)
	default:
		return LerpColor(yellow, red, (t-2.0/3.0)*3.0)
	}
}
