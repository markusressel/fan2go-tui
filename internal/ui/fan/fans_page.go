package fan

import (
	"fan2go-tui/internal/client"
	uiutil "fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
	"sort"
	"strings"
)

type FansPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	Fans map[string]*client.Fan

	layout *tview.Flex

	fanInfoComponents  []*FanInfoComponent
	fanGraphComponents []*FanGraphComponent
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

	fansPageLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	fanRowLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	fansPageLayout.AddItem(fanRowLayout, 0, 1, true)

	fans, fanIds, err := c.fetchFans()
	if err != nil {
		// TODO: handle error
		//c.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return fansPageLayout
	}
	c.Fans = *fans

	for _, fId := range fanIds {
		fanColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
		uiutil.SetupWindow(fanColumnLayout, fId)
		fanColumnLayout.SetTitleAlign(tview.AlignLeft)
		fanColumnLayout.SetBorder(true)
		fanRowLayout.AddItem(fanColumnLayout, 0, 1, true)

		f := (*fans)[fId]
		fanInfoComponent := NewFanInfoComponent(c.application, f)
		c.fanInfoComponents = append(c.fanInfoComponents, fanInfoComponent)
		fanInfoComponent.Refresh()
		layout := fanInfoComponent.GetLayout()

		fanColumnLayout.AddItem(layout, 0, 1, true)
		//fanInfoLayout.AddItem(layout, 0, 1, true)

		fanGraphComponent := NewFanGraphComponent(c.application, f)
		c.fanGraphComponents = append(c.fanGraphComponents, fanGraphComponent)
		//fanGraphComponent.SetTitle(f.Config.Id)
		fanGraphComponent.Refresh()
		layout = fanGraphComponent.GetLayout()
		fanColumnLayout.AddItem(layout, 0, 3, true)
	}

	return fansPageLayout
}

func (c *FansPage) fetchFans() (*map[string]*client.Fan, []string, error) {
	result, err := c.client.GetFans()
	if err != nil {
		return nil, nil, err
	}

	var fanIds []string
	for _, f := range *result {
		fanIds = append(fanIds, f.Config.Id)
	}

	sort.SliceStable(fanIds, func(i, j int) bool {
		a := fanIds[i]
		b := fanIds[j]

		result := strings.Compare(strings.ToLower(a), strings.ToLower(b))

		if result <= 0 {
			return true
		} else {
			return false
		}
	})

	return result, fanIds, nil
}

func (c *FansPage) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FansPage) Refresh() {
	fans, _, err := c.fetchFans()
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
