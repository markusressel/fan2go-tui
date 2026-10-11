package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/shortcut_helper"
	"fan2go-tui/internal/ui/util"
	"fan2go-tui/internal/uistate"
	"sort"
	"strings"

	"github.com/rivo/tview"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

const (
	// selectionNamespace is the namespace of the selected fan of each fan2go daemon, see PersistSelection
	selectionNamespace = "fans.selection"
	// maxRememberedServers is the number of fan2go daemons the selection is remembered for
	maxRememberedServers = 20
)

type FansPage struct {
	application *tview.Application

	store *state.Store

	layout *tview.Flex

	fanList *util.ListComponent[FanListItemComponent]
	// selection remembers the selected fan, see PersistSelection
	selection *util.ListSelection[FanListItemComponent]

	fanListItemComponents map[string]*FanListItemComponent
	onOpenCurve           func(curveID string)
}

func NewFansPage(application *tview.Application, store *state.Store, onOpenCurve func(curveID string)) FansPage {

	fansPage := FansPage{
		application:           application,
		store:                 store,
		fanListItemComponents: map[string]*FanListItemComponent{},
		onOpenCurve:           onOpenCurve,
	}

	fansPage.layout = fansPage.createLayout()

	return fansPage
}

func (c *FansPage) createLayout() *tview.Flex {
	fansPageLayout := tview.NewFlex()

	listConfig := util.NewListComponentConfig()
	fanListComponent := util.NewListComponent[FanListItemComponent](
		c.application,
		listConfig,
		func(entry *FanListItemComponent) (layout *tview.Flex) {
			return entry.GetLayout()
		},
		//func(a, b *FanListItemComponent) bool {
		//	x := a.Fan.Config.ID
		//	y := b.Fan.Config.ID
		//	return strings.Compare(strings.ToLower(x), strings.ToLower(y)) <= 0
		//},
		func(entries []*FanListItemComponent, inverted bool) []*FanListItemComponent {
			sort.SliceStable(entries, func(i, j int) bool {
				x := entries[i].FanState.Fan.Config.ID
				y := entries[j].FanState.Fan.Config.ID

				result := strings.Compare(strings.ToLower(x), strings.ToLower(y))

				if result <= 0 {
					return true
				} else {
					return false
				}
			})

			if inverted {
				slices.Reverse(entries)
			}

			return entries
		},
	)
	c.fanList = fanListComponent
	c.selection = util.NewListSelection(c.fanList, func(entry *FanListItemComponent) string {
		return entry.FanState.Fan.Config.ID
	})
	fansPageLayout.AddItem(c.fanList.GetLayout(), 0, 1, true)

	return fansPageLayout
}

func (c *FansPage) getFanIds(fans map[string]*client.Fan) []string {
	var fanIds []string
	for _, f := range fans {
		fanIds = append(fanIds, f.Config.ID)
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

	return fanIds
}

func (c *FansPage) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FansPage) Refresh() error {
	fans := c.store.GetFans()
	fanIds := c.getFanIds(fans)

	var fanListItemsComponents []*FanListItemComponent

	oldFIds := maps.Keys(c.fanListItemComponents)
	// remove now nonexisting entries
	for _, oldFId := range oldFIds {
		_, ok := fans[oldFId]
		if !ok {
			delete(c.fanListItemComponents, oldFId)
		}
	}

	// add new entries / update existing entries
	for _, fId := range fanIds {
		fanState := c.store.GetFanState(fId)
		fanListItemComponent, ok := c.fanListItemComponents[fId]
		if ok {
			fanListItemComponent.SetFan(fanState)
			fanListItemsComponents = append(fanListItemsComponents, fanListItemComponent)
		} else {
			fanListItemComponent = NewFanListItemComponent(c.application, fanState, c.onOpenCurve)
			c.fanListItemComponents[fId] = fanListItemComponent
			fanListItemComponent.SetFan(fanState)
			fanListItemsComponents = append(fanListItemsComponents, fanListItemComponent)
		}
	}

	c.fanList.SetData(fanListItemsComponents)
	c.selection.Update()

	return nil
}

func (c *FansPage) ScrollToItem() {
	c.fanList.SelectEntry(c.fanList.GetSelectedItem())
}

func (c *FansPage) GetSelectedItem() *FanListItemComponent {
	return c.fanList.GetSelectedItem()
}

func (c *FansPage) GetSelectedIndex() int {
	return c.fanList.GetSelectedIndex()
}

func (c *FansPage) GetShortcutMap() []shortcut_helper.ShortcutEntry {
	return []shortcut_helper.ShortcutEntry{
		{KeyCombo: []string{"↑", "↓"}, Name: "Select fan", Group: shortcut_helper.GroupNavigation},
		{KeyCombo: []string{shortcut_helper.KeyPgUp, shortcut_helper.KeyPgDn}, Name: "Scroll list", Group: shortcut_helper.GroupNavigation},
	}
}

// PersistSelection remembers the selected fan of the fan2go daemon at the address server in store, and selects the
// one remembered in the last run if restore is true. Must be called before the first Refresh.
func (c *FansPage) PersistSelection(store *uistate.Store, server string, restore bool) {
	c.selection.Persist(store, selectionNamespace, server, maxRememberedServers, restore)
}
