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
