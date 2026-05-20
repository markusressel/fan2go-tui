package graph

import "github.com/gdamore/tcell/v2"

type GraphSeriesLegend struct {
	text            string
	unit            string
	textColor       *tcell.Color
	backgroundColor *tcell.Color
	glyph           *rune
}

func NewGraphSeriesLegend(text string) *GraphSeriesLegend {
	return &GraphSeriesLegend{
		text: text,
	}
}

func (l *GraphSeriesLegend) WithUnit(unit string) *GraphSeriesLegend {
	l.unit = unit
	return l
}

func (l *GraphSeriesLegend) WithTextColor(color tcell.Color) *GraphSeriesLegend {
	l.textColor = &color
	return l
}

func (l *GraphSeriesLegend) WithBackgroundColor(color tcell.Color) *GraphSeriesLegend {
	l.backgroundColor = &color
	return l
}

func (l *GraphSeriesLegend) WithGlyph(glyph rune) *GraphSeriesLegend {
	l.glyph = &glyph
	return l
}

func (l *GraphSeriesLegend) displayText() string {
	if l == nil {
		return ""
	}
	if l.unit == "" {
		return l.text
	}
	return l.text + " (" + l.unit + ")"
}
