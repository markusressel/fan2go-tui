package util

import (
	"testing"
)

func TestWindowUtils(t *testing.T) {
	size := 3
	window := CreateRollingWindow(size)

	if window == nil {
		t.Fatal("CreateRollingWindow returned nil")
	}

	t.Run("FillAndAvg", func(t *testing.T) {
		FillWindow(window, size, 10.0)
		avg := GetWindowAvg(window)
		if avg != 10.0 {
			t.Errorf("Expected avg 10.0, got %v", avg)
		}
	})

	t.Run("Max", func(t *testing.T) {
		window.Append(20.0) // Window now: [10, 10, 20]
		max := GetWindowMax(window)
		if max != 20.0 {
			t.Errorf("Expected max 20.0, got %v", max)
		}
	})

	t.Run("RollingAvg", func(t *testing.T) {
		window.Append(30.0) // Window now: [10, 20, 30]
		avg := GetWindowAvg(window)
		if avg != 20.0 {
			t.Errorf("Expected avg 20.0, got %v", avg)
		}
	})
}
