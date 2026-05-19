package util

// Coerce returns a value that is at least min and at most max, otherwise value
func Coerce(value float64, min float64, max float64) float64 {
	if value > max {
		return max
	}
	if value < min {
		return min
	}
	return value
}

func Clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
