package util

type GraphComponentConfig struct {
	// Reversed determines if the graph should be drawn from left tp right instead of right to left
	// Default: false
	Reversed bool
}

func NewGraphComponentConfig() *GraphComponentConfig {
	return &GraphComponentConfig{
		Reversed: false,
	}
}

func (c *GraphComponentConfig) WithReversedOrder() *GraphComponentConfig {
	c.Reversed = true
	return c
}
