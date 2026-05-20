package graph

import (
	"math"
	"slices"
	"strconv"

	"golang.org/x/exp/maps"
)

// SeriesValueProvider defines how a graph series resolves x values, y values and x labels.
type SeriesValueProvider interface {
	// X returns the x value for a given index i. It should return NaN for out of range or invalid indices.
	X(i int) float64
	// F returns the y value for a given x value. It should return NaN for x values that are out of range or invalid.
	F(x float64) float64
	// XLabel returns the label for the x value at index i. The x value is also provided for convenience.
	XLabel(i int, x float64) string
}

// NewGraphLineFromSeriesValueProvider adapts a provider for a graph line.
func NewGraphLineFromSeriesValueProvider(name string, provider SeriesValueProvider) *GraphLine {
	return NewGraphLine(name, provider.X, provider.F, provider.XLabel)
}

// NewGraphBarFromSeriesValueProvider adapts a provider for a graph bar.
func NewGraphBarFromSeriesValueProvider(name string, provider SeriesValueProvider) *GraphBar {
	return NewGraphBar(name, provider.X, provider.F, provider.XLabel)
}

// DiscreteIntSeriesValueProvider maps sorted integer x keys to float64 y values.
type DiscreteIntSeriesValueProvider struct {
	keys   []int
	values map[int]float64
}

func NewDiscreteIntSeriesValueProvider(values map[int]float64) *DiscreteIntSeriesValueProvider {
	keys := maps.Keys(values)
	slices.Sort(keys)

	return &DiscreteIntSeriesValueProvider{
		keys:   keys,
		values: values,
	}
}

func (p *DiscreteIntSeriesValueProvider) X(i int) float64 {
	if i < 0 || i >= len(p.keys) {
		return math.NaN()
	}
	return float64(p.keys[i])
}

func (p *DiscreteIntSeriesValueProvider) F(x float64) float64 {
	val, ok := p.values[int(math.Floor(x))]
	if !ok {
		return math.NaN()
	}
	return val
}

func (p *DiscreteIntSeriesValueProvider) XLabel(_ int, x float64) string {
	return strconv.Itoa(int(math.Round(x)))
}

// RoundedSliceSeriesValueProvider resolves values from a slice by rounded x index.
type RoundedSliceSeriesValueProvider struct {
	values *[]float64
}

func NewRoundedSliceSeriesValueProvider(values *[]float64) *RoundedSliceSeriesValueProvider {
	return &RoundedSliceSeriesValueProvider{values: values}
}

func (p *RoundedSliceSeriesValueProvider) X(i int) float64 {
	return float64(i)
}

func (p *RoundedSliceSeriesValueProvider) F(x float64) float64 {
	idx := int(math.Round(x))
	if p.values == nil || idx < 0 || idx >= len(*p.values) {
		return math.NaN()
	}
	return (*p.values)[idx]
}

func (p *RoundedSliceSeriesValueProvider) XLabel(_ int, x float64) string {
	return strconv.Itoa(int(math.Round(x)))
}
