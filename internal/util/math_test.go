package util

import "testing"

func TestCoerce(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected float64
	}{
		{"within range", 5.0, 0.0, 10.0, 5.0},
		{"below min", -1.0, 0.0, 10.0, 0.0},
		{"above max", 11.0, 0.0, 10.0, 10.0},
		{"min equal max within", 5.0, 5.0, 5.0, 5.0},
		{"min equal max below", 4.0, 5.0, 5.0, 5.0},
		{"min equal max above", 6.0, 5.0, 5.0, 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Coerce(tt.value, tt.min, tt.max); got != tt.expected {
				t.Errorf("Coerce() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestClamp01(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{"within range", 0.5, 0.5},
		{"below 0", -0.1, 0.0},
		{"above 1", 1.1, 1.0},
		{"zero", 0.0, 0.0},
		{"one", 1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clamp01(tt.value); got != tt.expected {
				t.Errorf("Clamp01() = %v, want %v", got, tt.expected)
			}
		})
	}
}
