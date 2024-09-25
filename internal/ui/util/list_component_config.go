package util

type ListComponentConfig struct {
	MaxVisibleItems   int
	MinHeightPerEntry int
}

func NewListComponentConfig() *ListComponentConfig {
	return &ListComponentConfig{
		MaxVisibleItems:   0,
		MinHeightPerEntry: 18,
	}
}

func (c *ListComponentConfig) WithMaxVisibleItems(count int) *ListComponentConfig {
	c.MaxVisibleItems = count
	return c
}

func (c *ListComponentConfig) WithMinHeightPerEntry(height int) *ListComponentConfig {
	if height < 1 {
		height = 1
	}
	c.MinHeightPerEntry = height
	return c
}
