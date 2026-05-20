package graph

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestPickReadableTextColorForDarkBackground(t *testing.T) {
	bg := tcell.NewRGBColor(10, 20, 30)
	if got := pickReadableTextColor(bg); got != tcell.ColorWhite {
		t.Fatalf("expected white text for dark background, got %v", got)
	}
}

func TestPickReadableTextColorForLightBackground(t *testing.T) {
	bg := tcell.NewRGBColor(240, 240, 240)
	if got := pickReadableTextColor(bg); got != tcell.ColorBlack {
		t.Fatalf("expected black text for light background, got %v", got)
	}
}

func TestPickReadableTextColorForSeriesBlueBackground(t *testing.T) {
	bg := tcell.NewRGBColor(0, 120, 255)
	if got := pickReadableTextColor(bg); got != tcell.ColorBlack {
		t.Fatalf("expected black text for bright blue background, got %v", got)
	}
}
