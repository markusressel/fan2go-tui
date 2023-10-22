package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/util"
	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"sort"
	"strings"
)

type FansPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	layout *tview.Flex

	fanList *util.ListComponent[FanListItemComponent]

	fanListItemComponents map[string]*FanListItemComponent
}

func NewFansPage(application *tview.Application, c client.Fan2goApiClient) FansPage {

	fansPage := FansPage{
		application:           application,
		client:                c,
		fanListItemComponents: map[string]*FanListItemComponent{},
	}

	fansPage.layout = fansPage.createLayout()

	return fansPage
}

func (c *FansPage) createLayout() *tview.Flex {
	fansPageLayout := tview.NewFlex()
	fansPageLayout.SetDirection(tview.FlexRow)

	fanListComponent := util.NewListComponent[FanListItemComponent](
		c.application,
		func(entry *FanListItemComponent) (layout *tview.Flex) {
			return entry.GetLayout()
		},
		func(a, b *FanListItemComponent) bool {
			return strings.Compare(a.Fan.Config.Id, b.Fan.Config.Id) <= 0
		},
	)
	c.fanList = fanListComponent
	fansPageLayout.AddItem(c.fanList.GetLayout(), 0, 1, true)

	fansPageLayout.Focus(func(p tview.Primitive) {
		layout := c.fanList.GetLayout()
		c.application.SetFocus(layout)
	})

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

func (c *FansPage) Refresh() error {
	fans, fanIds, err := c.fetchFans()
	if err != nil || fans == nil {
		fans = &map[string]*client.Fan{}
	}

	var fanListItemsComponents []*FanListItemComponent

	oldFIds := maps.Keys(c.fanListItemComponents)
	// remove now nonexisting entries
	for _, oldFId := range oldFIds {
		_, ok := (*fans)[oldFId]
		if !ok {
			delete(c.fanListItemComponents, oldFId)
		}
	}

	// add new entries / update existing entries
	for _, fId := range fanIds {
		fan := (*fans)[fId]
		fanListItemComponent, ok := c.fanListItemComponents[fId]
		if ok {
			fanListItemComponent.SetFan(fan)
			fanListItemsComponents = append(fanListItemsComponents, fanListItemComponent)
		} else {
			fanListItemComponent = NewFanListItemComponent(c.application, fan)
			c.fanListItemComponents[fId] = fanListItemComponent
			fanListItemComponent.SetFan(fan)
			fanListItemsComponents = append(fanListItemsComponents, fanListItemComponent)
		}
	}

	c.fanList.SetData(fanListItemsComponents)

	return err
}
