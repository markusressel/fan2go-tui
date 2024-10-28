package util

func RotateSliceLeft[T any](slice []T, n int) []T {
	if len(slice) == 0 {
		return slice
	}
	for i := 0; i < n; i++ {
		slice = append(slice[1:], slice[0])
	}
	return slice
}

func RotateSliceRight[T any](slice []T, n int) []T {
	if len(slice) == 0 {
		return slice
	}
	for i := 0; i < n; i++ {
		slice = append([]T{slice[len(slice)-1]}, slice[:len(slice)-1]...)
	}

	return slice
}

func RotateSliceBy[T any](slice []T, n int) []T {
	if len(slice) == 0 {
		return slice
	}
	if n < 0 {
		return RotateSliceLeft(slice, -1*n)
	}
	if n > 0 {
		return RotateSliceRight(slice, n)
	}
	return slice
}

func DistributeValuesOverRange(values []float64, totalRange int) []float64 {
	valuesCount := len(values)
	slotsPerValue := float64(totalRange) / float64(valuesCount)

	result := make([]float64, totalRange)
	for i := 0; i < totalRange; i++ {
		slotIndex := int(float64(i) / slotsPerValue)
		if slotIndex >= valuesCount {
			// use the last known value for the remaining slots
			slotIndex = valuesCount - 1
		}
		result[i] = values[slotIndex]
	}
	return result
}
