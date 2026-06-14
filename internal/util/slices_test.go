package util

import (
	"reflect"
	"testing"
)

func TestRotateSliceLeft(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		expected []int
	}{
		{"empty", []int{}, 1, []int{}},
		{"rotate 1", []int{1, 2, 3}, 1, []int{2, 3, 1}},
		{"rotate 2", []int{1, 2, 3}, 2, []int{3, 1, 2}},
		{"rotate 3", []int{1, 2, 3}, 3, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RotateSliceLeft(tt.slice, tt.n); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RotateSliceLeft() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRotateSliceRight(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		expected []int
	}{
		{"empty", []int{}, 1, []int{}},
		{"rotate 1", []int{1, 2, 3}, 1, []int{3, 1, 2}},
		{"rotate 2", []int{1, 2, 3}, 2, []int{2, 3, 1}},
		{"rotate 3", []int{1, 2, 3}, 3, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RotateSliceRight(tt.slice, tt.n); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RotateSliceRight() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestRotateSliceBy(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		n        int
		expected []int
	}{
		{"empty", []int{}, 1, []int{}},
		{"zero", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"positive (right)", []int{1, 2, 3}, 1, []int{3, 1, 2}},
		{"negative (left)", []int{1, 2, 3}, -1, []int{2, 3, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RotateSliceBy(tt.slice, tt.n); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("RotateSliceBy() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDistributeValuesOverRange(t *testing.T) {
	tests := []struct {
		name       string
		values     []float64
		totalRange int
		expected   []float64
	}{
		{"empty values", []float64{}, 5, []float64{0, 0, 0, 0, 0}},
		{"single value", []float64{1.0}, 3, []float64{1.0, 1.0, 1.0}},
		{"two values", []float64{1.0, 2.0}, 4, []float64{1.0, 1.0, 2.0, 2.0}},
		{"three values", []float64{1.0, 2.0, 3.0}, 6, []float64{1.0, 1.0, 2.0, 2.0, 3.0, 3.0}},
		{"more values than range", []float64{1, 2, 3}, 2, []float64{1, 2}},
		{"totalRange is 0", []float64{1, 2, 3}, 0, []float64{}},
		{"valuesCount much larger than totalRange", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 1, []float64{1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DistributeValuesOverRange(tt.values, tt.totalRange); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DistributeValuesOverRange() = %v, want %v", got, tt.expected)
			}
		})
	}
}
