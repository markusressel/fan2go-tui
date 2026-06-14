package util

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestLerpColor(t *testing.T) {
	black := [3]int32{0, 0, 0}
	white := [3]int32{255, 255, 255}

	tests := []struct {
		name     string
		a, b     [3]int32
		t        float64
		expected tcell.Color
	}{
		{"t=0", black, white, 0.0, tcell.NewRGBColor(0, 0, 0)},
		{"t=1", black, white, 1.0, tcell.NewRGBColor(255, 255, 255)},
		{"t=0.5", black, white, 0.5, tcell.NewRGBColor(128, 128, 128)},
		{"t=-0.1 (clamped)", black, white, -0.1, tcell.NewRGBColor(0, 0, 0)},
		{"t=1.1 (clamped)", black, white, 1.1, tcell.NewRGBColor(255, 255, 255)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LerpColor(tt.a, tt.b, tt.t)
			if got != tt.expected {
				t.Errorf("LerpColor() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestGradientColorAt(t *testing.T) {
	blue := [3]int32{0, 0, 255}
	green := [3]int32{0, 255, 0}
	yellow := [3]int32{255, 255, 0}
	red := [3]int32{255, 0, 0}

	tests := []struct {
		name     string
		t        float64
		expected tcell.Color
	}{
		{"t=0 (blue)", 0.0, tcell.NewRGBColor(0, 0, 255)},
		{"t=0.33 (green-ish)", 1.0 / 3.0, tcell.NewRGBColor(0, 255, 0)},
		{"t=0.5 (between green and yellow)", 0.5, tcell.NewRGBColor(128, 255, 0)},
		{"t=0.66 (yellow-ish)", 2.0 / 3.0, tcell.NewRGBColor(255, 255, 0)},
		{"t=1.0 (red)", 1.0, tcell.NewRGBColor(255, 0, 0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GradientColorAt(tt.t, blue, green, yellow, red)
			if got != tt.expected {
				t.Errorf("GradientColorAt(%v) = %v, want %v", tt.t, got, tt.expected)
			}
		})
	}
}
