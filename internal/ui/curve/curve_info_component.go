package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/txwidget"

	"github.com/rivo/tview"
)

type CurveInfoComponent struct {
	application *tview.Application

	Curve *client.Curve

	layout *tview.Flex

	configComponent *txwidget.ConfigInfoComponent
}

func NewCurveInfoComponent(application *tview.Application, curve *client.Curve) *CurveInfoComponent {
	c := &CurveInfoComponent{
		application: application,
		Curve:       curve,
	}

	c.layout = c.createLayout()

	return c
}

func (c *CurveInfoComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	configComponent := txwidget.NewConfigInfoComponent()
	layout.AddItem(configComponent.GetPrimitive(), 0, 1, false)
	c.configComponent = configComponent

	return layout
}

func (c *CurveInfoComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurveInfoComponent) SetCurve(curve *client.Curve) {
	c.Curve = curve
	c.refresh()
}

func (c *CurveInfoComponent) refresh() {
	if c.Curve == nil {
		c.configComponent.SetSections(nil)
		return
	}

	config := c.Curve.Config
	c.configComponent.SetSections(txwidget.CurveConfigSections(config))
}
