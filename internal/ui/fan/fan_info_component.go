package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/txwidgets"
	"strings"

	"github.com/rivo/tview"
)

type FanInfoComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout *tview.Flex

	warningsComponent *txwidgets.WarningsComponent
	configComponent   *txwidgets.ConfigInfoComponent
	onOpenCurve       func(curveID string)
}

func NewFanInfoComponent(application *tview.Application, fan *client.Fan, onOpenCurve func(curveID string)) *FanInfoComponent {
	c := &FanInfoComponent{
		application: application,
		Fan:         fan,
		onOpenCurve: onOpenCurve,
	}

	c.layout = c.createLayout()

	return c
}

func (c *FanInfoComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	// hidden until there are warnings, see refresh
	c.warningsComponent = txwidgets.NewWarningsComponent()
	layout.AddItem(c.warningsComponent.GetPrimitive(), 0, 0, false)

	configComponent := txwidgets.NewConfigInfoComponent()
	configComponent.SetFieldClickablePredicate(func(sectionTitle, label, value string) bool {
		return strings.EqualFold(sectionTitle, "General") && strings.EqualFold(label, "Curve") && value != ""
	})
	configComponent.SetFieldClickHandler(func(sectionTitle, label, value string) {
		if c.onOpenCurve != nil && strings.EqualFold(sectionTitle, "General") && strings.EqualFold(label, "Curve") {
			c.onOpenCurve(value)
		}
	})
	layout.AddItem(configComponent.GetPrimitive(), 0, 1, false)
	c.configComponent = configComponent

	return layout
}

func (c *FanInfoComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanInfoComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
	c.refresh()
}

func (c *FanInfoComponent) refresh() {
	if c.Fan == nil {
		c.setWarnings(nil)
		c.configComponent.SetSections(nil)
		return
	}

	c.setWarnings(c.Fan.Warnings)

	config := c.Fan.Config
	c.configComponent.SetSections(txwidgets.FanConfigSections(config))
}

func (c *FanInfoComponent) ScrollHorizontal(delta int) {
	if c.configComponent != nil {
		c.configComponent.ScrollHorizontal(delta)
	}
}

// setWarnings shows the warnings of the fan above its configuration, and hides them if there are none.
func (c *FanInfoComponent) setWarnings(warnings []client.Warning) {
	c.warningsComponent.SetWarnings(warnings)
	proportion := 0
	if c.warningsComponent.HasWarnings() {
		proportion = 1
	}
	c.layout.ResizeItem(c.warningsComponent.GetPrimitive(), 0, proportion)
}
