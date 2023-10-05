package fan

import (
	"fan2go-tui/internal/client"
	"github.com/rivo/tview"
)

type FansPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Fans map[string]*client.Fan

	layout *tview.Flex

	fanInfoComponents    []*FanInfoComponent
	fanOverviewComponent *FanGraphsComponent
	fanGraphComponents   []*FanGraphComponent
}

func NewFansPage(application *tview.Application, client client.Fan2goApiClient) FansPage {

	fansPage := FansPage{
		application: application,
		client:      client,
	}

	fansPage.layout = fansPage.createLayout()

	return fansPage
}

func (c *FansPage) createLayout() *tview.Flex {

	fansPageLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	fanInfoLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	fansPageLayout.AddItem(fanInfoLayout, 0, 1, true)
	fanGraphsLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	fansPageLayout.AddItem(fanGraphsLayout, 0, 3, true)

	fans, err := c.fetchFans()
	if err != nil {
		// TODO: handle error
		//c.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return fansPageLayout
	}
	c.Fans = *fans

	for _, f := range c.Fans {
		fanInfoComponent := NewFanInfoComponent(c.application, f)
		c.fanInfoComponents = append(c.fanInfoComponents, fanInfoComponent)
		fanInfoComponent.Refresh()
		layout := fanInfoComponent.GetLayout()
		fanInfoLayout.AddItem(layout, 0, 1, true)

		fanGraphComponent := NewFanGraphComponent(c.application, f)
		c.fanGraphComponents = append(c.fanGraphComponents, fanGraphComponent)
		fanGraphComponent.SetTitle(f.Config.Id)
		fanGraphComponent.Refresh()
		layout = fanGraphComponent.GetLayout()
		fanGraphsLayout.AddItem(layout, 0, 1, false)
	}

	return fansPageLayout
}

func (c *FansPage) fetchFans() (*map[string]*client.Fan, error) {
	return c.client.GetFans()
}

func (c *FansPage) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FansPage) Refresh() {
	fans, err := c.fetchFans()
	if err != nil {
		return
	}
	if fans == nil {
		return
	}

	c.Fans = *fans
	for _, component := range c.fanInfoComponents {
		fan, ok := (*fans)[component.Fan.Config.Id]
		if !ok {
			continue
		}
		component.SetFan(fan)
		component.Refresh()
	}
	for _, component := range c.fanGraphComponents {
		if component.Fan == nil {
			continue
		}
		fan, ok := (*fans)[component.Fan.Config.Id]
		if !ok || fan == nil {
			continue
		}
		component.SetFan(fan)
		component.Refresh()
	}
}

func (c *FansPage) SetFans(fans *map[string]*client.Fan) {
	c.Fans = *fans
	c.Refresh()
}
