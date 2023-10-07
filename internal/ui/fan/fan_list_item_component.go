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

	fanInfoComponent  *FanInfoComponent
	fanGraphComponent *FanGraphComponent
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
	uiutil.SetupWindow(fanColumnLayout, c.Fan.Label)
	fanColumnLayout.SetTitleAlign(tview.AlignLeft)
	fanColumnLayout.SetBorder(true)
	rootLayout.AddItem(fanColumnLayout, 0, 1, true)

	f := c.Fan
	fanInfoComponent := NewFanInfoComponent(c.application, f)
	c.fanInfoComponent = fanInfoComponent
	fanInfoComponent.Refresh()
	layout := fanInfoComponent.GetLayout()

	fanColumnLayout.AddItem(layout, 0, 1, true)

	fanGraphComponent := NewFanGraphComponent(c.application, f)
	c.fanGraphComponent = fanGraphComponent
	fanGraphComponent.Refresh()
	layout = fanGraphComponent.GetLayout()
	fanColumnLayout.AddItem(layout, 0, 3, true)

	return rootLayout
}

func (c *FanListItemComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanListItemComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
}

func (c *FanListItemComponent) Refresh() {
	c.fanInfoComponent.SetFan(c.Fan)
	c.fanGraphComponent.InsertValue(c.Fan)
}
