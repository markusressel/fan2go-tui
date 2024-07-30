package util

type ListComponentConfig struct {
	MaxVisibleItems int
}

func NewListComponentConfig() *ListComponentConfig {
	return &ListComponentConfig{
		MaxVisibleItems: 3,
	}
}

func (c *ListComponentConfig) WithMaxVisibleItems(count int) *ListComponentConfig {
	c.MaxVisibleItems = count
	return c
}
