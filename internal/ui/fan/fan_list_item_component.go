package fan

import (
	"fan2go-tui/internal/client"
	uiutil "fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
)

type FanListItemComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout *tview.Flex

	fanInfoComponent     *FanInfoComponent
	fanGraphComponent    *FanGraphComponent
	fanRpmCurveComponent *FanRpmCurveComponent
}

func NewFanListItemComponent(application *tview.Application, fan *client.Fan) *FanListItemComponent {
	c := &FanListItemComponent{
		application: application,
		Fan:         fan,
	}

	c.layout = c.createLayout()

	return c
}

func (c *FanListItemComponent) createLayout() *tview.Flex {
	rootLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	fanColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	uiutil.SetupWindow(fanColumnLayout, c.Fan.Config.ID)
	fanColumnLayout.SetTitleAlign(tview.AlignLeft)
	fanColumnLayout.SetBorder(true)
	rootLayout.AddItem(fanColumnLayout, 0, 1, true)

	f := c.Fan
	fanInfoComponent := NewFanInfoComponent(c.application, f)
	c.fanInfoComponent = fanInfoComponent
	fanInfoComponent.SetFan(f)
	layout := fanInfoComponent.GetLayout()
	fanColumnLayout.AddItem(layout, 0, 1, true)

	fanGraphsRowLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	fanColumnLayout.AddItem(fanGraphsRowLayout, 0, 3, true)

	{
		fanGraphComponent := NewFanGraphComponent(c.application, f)
		c.fanGraphComponent = fanGraphComponent
		fanGraphComponent.SetFan(f)
		layout = fanGraphComponent.GetLayout()
		fanGraphsRowLayout.AddItem(layout, 0, 1, true)

		if f.FanCurveData != nil {
			fanRpmCurveComponent := NewFanRpmCurveComponent(c.application, f)
			c.fanRpmCurveComponent = fanRpmCurveComponent
			fanRpmCurveComponent.SetFan(f)
			layout = fanRpmCurveComponent.GetLayout()
			fanGraphsRowLayout.AddItem(layout, 0, 1, true)
		}
	}

	return rootLayout
}

func (c *FanListItemComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanListItemComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
	c.refresh()
}

func (c *FanListItemComponent) refresh() {
	c.fanInfoComponent.SetFan(c.Fan)
	c.fanGraphComponent.SetFan(c.Fan)
	c.fanRpmCurveComponent.SetFan(c.Fan)
}
