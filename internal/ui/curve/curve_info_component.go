package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/txwidgets"
	"strings"

	"github.com/rivo/tview"
)

type CurveInfoComponent struct {
	application *tview.Application

	Curve *client.Curve

	layout *tview.Flex

	configComponent *txwidgets.ConfigInfoComponent
	onOpenSensor    func(sensorID string)
	onOpenCurve     func(curveID string)
}

func NewCurveInfoComponent(application *tview.Application, curve *client.Curve, onOpenSensor func(sensorID string), onOpenCurve func(curveID string)) *CurveInfoComponent {
	c := &CurveInfoComponent{
		application:  application,
		Curve:        curve,
		onOpenSensor: onOpenSensor,
		onOpenCurve:  onOpenCurve,
	}

	c.layout = c.createLayout()

	return c
}

func (c *CurveInfoComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	configComponent := txwidgets.NewConfigInfoComponent()
	configComponent.SetFieldClickablePredicate(func(sectionTitle, label, value string) bool {
		isSensor := strings.EqualFold(sectionTitle, "Curve") && strings.EqualFold(label, "Sensor") && value != ""
		isCurveRef := strings.EqualFold(sectionTitle, "Curve") && (strings.EqualFold(label, "Curves") || label == "") && strings.HasPrefix(value, "- ")
		return isSensor || isCurveRef
	})
	configComponent.SetFieldClickHandler(func(sectionTitle, label, value string) {
		if !strings.EqualFold(sectionTitle, "Curve") {
			return
		}
		if strings.EqualFold(label, "Sensor") {
			if c.onOpenSensor != nil {
				c.onOpenSensor(value)
			}
			return
		}
		if (strings.EqualFold(label, "Curves") || label == "") && strings.HasPrefix(value, "- ") {
			curveID := strings.TrimSpace(strings.TrimPrefix(value, "- "))
			if curveID != "" && c.onOpenCurve != nil {
				c.onOpenCurve(curveID)
			}
		}
	})
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
	c.configComponent.SetSections(txwidgets.CurveConfigSections(config))
}

func (c *CurveInfoComponent) ScrollHorizontal(delta int) {
	if c.configComponent != nil {
		c.configComponent.ScrollHorizontal(delta)
	}
}
