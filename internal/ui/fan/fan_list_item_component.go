package fan

import (
	"fan2go-tui/internal/client"
	uiutil "fan2go-tui/internal/ui/util"

	"github.com/rivo/tview"
)

type FanListItemComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout          *tview.Flex
	graphHostLayout *tview.Flex

	activeGraphVariant fanGraphVariant

	fanInfoComponent     *FanInfoComponent
	fanGraphComponent    *FanGraphComponent
	fanRpmCurveComponent *FanRpmCurveComponent
}

type fanGraphVariant int

const (
	fanGraphVariantNone fanGraphVariant = iota
	fanGraphVariantHistory
	fanGraphVariantCurve
)

func hasFanCurveData(fan *client.Fan) bool {
	if fan == nil || fan.FanCurveData == nil {
		return false
	}
	return len(*fan.FanCurveData) > 0
}

func NewFanListItemComponent(application *tview.Application, fan *client.Fan, onOpenCurve func(curveID string)) *FanListItemComponent {
	c := &FanListItemComponent{
		application: application,
		Fan:         fan,
	}

	c.layout = c.createLayout(onOpenCurve)

	return c
}

func (c *FanListItemComponent) createLayout(onOpenCurve func(curveID string)) *tview.Flex {
	rootLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	fanColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	uiutil.SetupWindow(fanColumnLayout, c.Fan.Config.ID)
	fanColumnLayout.SetTitleAlign(tview.AlignLeft)
	fanColumnLayout.SetBorder(true)
	rootLayout.AddItem(fanColumnLayout, 0, 1, true)

	f := c.Fan
	fanInfoComponent := NewFanInfoComponent(c.application, f, onOpenCurve)
	c.fanInfoComponent = fanInfoComponent
	fanInfoComponent.SetFan(f)
	layout := fanInfoComponent.GetLayout()
	fanColumnLayout.AddItem(layout, 0, 1, true)
	fanColumnLayout.AddItem(tview.NewBox(), 1, 0, false)

	c.graphHostLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	fanColumnLayout.AddItem(c.graphHostLayout, 0, 3, true)

	c.ensureGraphVariant()
	c.refresh()

	return rootLayout
}

func (c *FanListItemComponent) desiredGraphVariant() fanGraphVariant {
	if hasFanCurveData(c.Fan) {
		return fanGraphVariantCurve
	}
	return fanGraphVariantHistory
}

func (c *FanListItemComponent) ensureGraphVariant() {
	if c == nil || c.graphHostLayout == nil {
		return
	}

	desiredVariant := c.desiredGraphVariant()
	if desiredVariant == c.activeGraphVariant {
		return
	}

	if c.fanGraphComponent != nil {
		c.graphHostLayout.RemoveItem(c.fanGraphComponent.GetLayout())
		c.fanGraphComponent = nil
	}
	if c.fanRpmCurveComponent != nil {
		c.graphHostLayout.RemoveItem(c.fanRpmCurveComponent.GetLayout())
		c.fanRpmCurveComponent = nil
	}

	if desiredVariant == fanGraphVariantCurve {
		c.fanRpmCurveComponent = NewFanRpmCurveComponent(c.application, c.Fan)
		c.graphHostLayout.AddItem(c.fanRpmCurveComponent.GetLayout(), 0, 1, true)
	} else {
		c.fanGraphComponent = NewFanGraphComponent(c.application, c.Fan)
		c.graphHostLayout.AddItem(c.fanGraphComponent.GetLayout(), 0, 1, true)
	}

	c.activeGraphVariant = desiredVariant
}

func (c *FanListItemComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanListItemComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
	c.ensureGraphVariant()
	c.refresh()
}

func (c *FanListItemComponent) refresh() {
	if c.fanInfoComponent != nil {
		c.fanInfoComponent.SetFan(c.Fan)
	}
	if c.fanGraphComponent != nil {
		c.fanGraphComponent.SetFan(c.Fan)
	}
	if c.fanRpmCurveComponent != nil {
		c.fanRpmCurveComponent.SetFan(c.Fan)
	}
}

func (c *FanListItemComponent) ScrollHorizontal(delta int) {
	if c.fanInfoComponent != nil {
		c.fanInfoComponent.ScrollHorizontal(delta)
	}
}
