package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/curve"
	"fan2go-tui/internal/ui/fan"
	"fan2go-tui/internal/ui/sensor"
	"fan2go-tui/internal/ui/shortcut_helper"
	"fan2go-tui/internal/ui/status_message"
	"fan2go-tui/internal/ui/util"
	"slices"
	"sync"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page string

const (
	GraphTestPage Page = "graph_test"
	FansPage      Page = "fans"
	CurvesPage    Page = "curves"
	SensorsPage   Page = "sensors"
)

type MainPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	layout      *tview.Flex
	header      *ApplicationHeaderComponent
	shortcutMap *shortcut_helper.ShortcutMapComponent

	page                Page
	mainPagePagerLayout *tview.Pages
	refreshMutex        sync.Mutex

	pagesMap orderedmap.OrderedMap[Page, util.PagesPage]
}

func NewMainPage(application *tview.Application, client client.Fan2goApiClient) *MainPage {

	mainPage := &MainPage{
		application: application,
		client:      client,
		page:        FansPage,
		pagesMap:    *orderedmap.NewOrderedMap[Page, util.PagesPage](),
	}

	mainPage.shortcutMap = shortcut_helper.NewShortcutMap(application)

	mainPage.layout = mainPage.createLayout()
	mainPage.layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		rune := event.Rune()
		if key == tcell.KeyTab || key == tcell.KeyBacktab {
			// TODO: Switch between tabs
		} else if key == tcell.KeyCtrlR {

		} else if rune == int32(49) {
			page := mainPage.GetPageAtIndex(0)
			mainPage.SetPage(page)
		} else if rune == int32(50) {
			page := mainPage.GetPageAtIndex(1)
			mainPage.SetPage(page)
		} else if rune == int32(51) {
			page := mainPage.GetPageAtIndex(2)
			mainPage.SetPage(page)
		}
		return event
	})

	return mainPage
}

func (mainPage *MainPage) createLayout() *tview.Flex {
	mainPageLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	mainPagePagerLayout := tview.NewPages()
	mainPage.mainPagePagerLayout = mainPagePagerLayout

	//graphTestPage := graph_test.NewGraphTestPage(mainPage.application)
	fansPage := fan.NewFansPage(mainPage.application, mainPage.client, mainPage.OpenCurveByID)
	curvesPage := curve.NewCurvesPage(mainPage.application, mainPage.client, mainPage.OpenSensorByID, mainPage.OpenCurveByID)
	sensorsPage := sensor.NewSensorsPage(mainPage.application, mainPage.client)

	//mainPage.AddPage(GraphTestPage, &graphTestPage, true)
	mainPage.AddPage(FansPage, &fansPage)
	mainPage.AddPage(CurvesPage, &curvesPage)
	mainPage.AddPage(SensorsPage, &sensorsPage)

	header := NewApplicationHeader(
		mainPage.application,
		mainPage.pagesMap,
	)
	mainPage.header = header

	mainPageLayout.AddItem(header.layout, 1, 0, false)
	mainPageLayout.AddItem(mainPagePagerLayout, 0, 1, true)
	mainPageLayout.AddItem(mainPage.shortcutMap.GetLayout(), 1, 0, false)

	for page, pagesPage := range mainPage.pagesMap.Iterator() {
		mainPage.mainPagePagerLayout.AddPage(
			string(page),
			pagesPage.GetLayout(),
			true,
			page == mainPage.page,
		)
	}

	return mainPageLayout
}

func (mainPage *MainPage) AddPage(s Page, pagesPage util.PagesPage) {
	mainPage.pagesMap.Set(s, pagesPage)
}

func (mainPage *MainPage) Init() {
	mainPage.application.ForceDraw()
	mainPage.Refresh()
	mainPage.scrollCurrentPageToItem()
	mainPage.updateShortcutMap()
}

func (mainPage *MainPage) Refresh() {
	mainPage.refreshMutex.Lock()
	defer mainPage.refreshMutex.Unlock()

	defer mainPage.application.ForceDraw()
	mainPage.UpdateHeader()

	currentPage := mainPage.GetCurrentPage()
	err := currentPage.Refresh()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return
	}

	mainPage.clearStatusMessage()
}

func (mainPage *MainPage) GetCurrentPage() util.PagesPage {
	val, _ := mainPage.pagesMap.Get(mainPage.page)
	return val
}

func (mainPage *MainPage) GetPageAtIndex(i int) Page {
	page := mainPage.pagesMap.Keys()[i]
	return page
}

func (mainPage *MainPage) SetPage(page Page) {
	mainPage.page = page
	mainPage.header.SetPage(page)

	mainPage.mainPagePagerLayout.SwitchToPage(string(page))
	mainPage.application.ForceDraw()
	mainPage.Refresh()

	mainPage.scrollCurrentPageToItem()
	mainPage.updateShortcutMap()
}

func (mainPage *MainPage) scrollCurrentPageToItem() {
	pagesPage, _ := mainPage.pagesMap.Get(mainPage.page)
	scrollablePage, ok := pagesPage.(util.CanScrollToItem)
	if ok {
		scrollablePage.ScrollToItem()
	}
}

func (mainPage *MainPage) PreviousPage() {
	keys := mainPage.pagesMap.Keys()
	currentIndex := slices.Index(keys, mainPage.page)
	nextIndex := (len(keys) + currentIndex - 1) % len(keys)
	mainPage.SetPage(keys[nextIndex])
}

func (mainPage *MainPage) NextPage() {
	keys := mainPage.pagesMap.Keys()
	currentIndex := slices.Index(keys, mainPage.page)
	nextIndex := (currentIndex + 1) % len(keys)
	mainPage.SetPage(keys[nextIndex])
}

func (mainPage *MainPage) OpenCurveByID(curveID string) {
	mainPage.SetPage(CurvesPage)

	pagesPage, _ := mainPage.pagesMap.Get(CurvesPage)
	curveSelectablePage, ok := pagesPage.(interface{ SelectCurveByID(curveID string) bool })
	if ok {
		curveSelectablePage.SelectCurveByID(curveID)
	}
}

func (mainPage *MainPage) OpenSensorByID(sensorID string) {
	mainPage.SetPage(SensorsPage)

	pagesPage, _ := mainPage.pagesMap.Get(SensorsPage)
	sensorSelectablePage, ok := pagesPage.(interface{ SelectSensorByID(sensorID string) bool })
	if ok {
		sensorSelectablePage.SelectSensorByID(sensorID)
	}
}

func (mainPage *MainPage) clearStatusMessage() {
	mainPage.header.SetStatus(status_message.NewInfoStatusMessage(""))
}

func (mainPage *MainPage) showStatusMessage(status *status_message.StatusMessage) {
	mainPage.header.SetStatus(status)
}

func (mainPage *MainPage) UpdateHeader() {
	mainPage.header.Refresh()
}

func (mainPage *MainPage) setShortcutMap(shortcutEntries []shortcut_helper.ShortcutEntry) {
	mainPage.shortcutMap.SetEntries(shortcutEntries)
}

func (mainPage *MainPage) clearShortcutMap() {
	mainPage.shortcutMap.Clear()
}

func (mainPage *MainPage) updateShortcutMap() {
	component := mainPage.GetCurrentPage()
	if c, ok := component.(shortcut_helper.ShortcutMapProvider); ok {
		shortcutMap := c.GetShortcutMap()

		globalShortcutMapEntries := []shortcut_helper.ShortcutEntry{
			{KeyCombo: []string{"?"}, Name: "Help"},
			{KeyCombo: []string{"Tab"}, Name: "Next"},
			{KeyCombo: []string{"1-3"}, Name: "Switch"},
			{KeyCombo: []string{"Ctrl+Q"}, Name: "Quit"},
		}

		shortcutMap = append(shortcutMap, globalShortcutMapEntries...)
		mainPage.setShortcutMap(shortcutMap)
	} else {
		mainPage.clearShortcutMap()
	}
}
