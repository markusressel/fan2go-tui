package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/curve"
	"fan2go-tui/internal/ui/fan"
	"fan2go-tui/internal/ui/sensor"
	"fan2go-tui/internal/ui/status_message"
	"fan2go-tui/internal/ui/util"
	"github.com/elliotchance/orderedmap/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slices"
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

	layout *tview.Flex
	header *ApplicationHeaderComponent

	page                Page
	mainPagePagerLayout *tview.Pages

	pagesMap orderedmap.OrderedMap[Page, util.PagesPage]
}

func NewMainPage(application *tview.Application, client client.Fan2goApiClient) *MainPage {

	mainPage := &MainPage{
		application: application,
		client:      client,
		page:        FansPage,
		pagesMap:    *orderedmap.NewOrderedMap[Page, util.PagesPage](),
	}

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
	fansPage := fan.NewFansPage(mainPage.application, mainPage.client)
	curvesPage := curve.NewCurvesPage(mainPage.application, mainPage.client)
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
	mainPage.Refresh()
}

func (mainPage *MainPage) Refresh() {
	defer mainPage.application.ForceDraw()
	mainPage.UpdateHeader()

	currentPage := mainPage.GetCurrentPage()
	err := currentPage.Refresh()

	for page, pagesPage := range mainPage.pagesMap.Iterator() {
		if page == mainPage.page {
			continue
		}
		err = pagesPage.Refresh()
		if err != nil {
			mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
			return
		}
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
	mainPage.Refresh()

	pagesPages, _ := mainPage.pagesMap.Get(mainPage.page)

	switch pagesPages.(type) {
	case util.CanScrollToItem:
		pagesPages.(util.CanScrollToItem).ScrollToItem()
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

func (mainPage *MainPage) clearStatusMessage() {
	mainPage.header.SetStatus(status_message.NewInfoStatusMessage(""))
}

func (mainPage *MainPage) showStatusMessage(status *status_message.StatusMessage) {
	mainPage.header.SetStatus(status)
}

func (mainPage *MainPage) UpdateHeader() {
	mainPage.header.Refresh()
}
