package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/curve"
	"fan2go-tui/internal/ui/fan"
	"fan2go-tui/internal/ui/sensor"
	"fan2go-tui/internal/ui/status_message"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"slices"
)

type Page string

const (
	FansPage    Page = "fans"
	CurvesPage  Page = "curves"
	SensorsPage Page = "sensors"
)

var (
	Pages = []Page{
		FansPage,
		CurvesPage,
		SensorsPage,
	}
)

type MainPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	layout *tview.Flex
	header *ApplicationHeaderComponent

	page                Page
	mainPagePagerLayout *tview.Pages

	fansPage    *fan.FansPage
	curvesPage  *curve.CurvesPage
	sensorsPage *sensor.SensorsPage
}

func NewMainPage(application *tview.Application, client client.Fan2goApiClient) *MainPage {

	mainPage := &MainPage{
		application: application,
		client:      client,
		page:        FansPage,
	}

	mainPage.layout = mainPage.createLayout()
	mainPage.layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		rune := event.Rune()
		if key == tcell.KeyTab || key == tcell.KeyBacktab {
			// TODO: Switch between tabs
		} else if key == tcell.KeyCtrlR {

		} else if rune == int32(49) {
			page := Pages[0]
			mainPage.SetPage(page)
		} else if rune == int32(50) {
			page := Pages[1]
			mainPage.SetPage(page)
		} else if rune == int32(51) {
			page := Pages[2]
			mainPage.SetPage(page)
		}
		return event
	})

	return mainPage
}

func (mainPage *MainPage) createLayout() *tview.Flex {
	mainPageLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	header := NewApplicationHeader(mainPage.application)
	mainPage.header = header
	mainPageLayout.AddItem(header.layout, 1, 0, false)

	mainPagePagerLayout := tview.NewPages()
	mainPage.mainPagePagerLayout = mainPagePagerLayout
	mainPageLayout.AddItem(mainPagePagerLayout, 0, 1, true)

	fansPage := fan.NewFansPage(mainPage.application, mainPage.client)
	mainPage.fansPage = &fansPage

	curvesPage := curve.NewCurvesPage(mainPage.application, mainPage.client)
	mainPage.curvesPage = &curvesPage

	sensorsPage := sensor.NewSensorsPage(mainPage.application, mainPage.client)
	mainPage.sensorsPage = &sensorsPage

	fansPageLayout := mainPage.fansPage.GetLayout()
	curvesPageLayout := mainPage.curvesPage.GetLayout()
	sensorsPageLayout := mainPage.sensorsPage.GetLayout()

	mainPagePagerLayout.AddPage(string(FansPage), fansPageLayout, true, true)
	mainPagePagerLayout.AddPage(string(CurvesPage), curvesPageLayout, true, false)
	mainPagePagerLayout.AddPage(string(SensorsPage), sensorsPageLayout, true, false)

	return mainPageLayout
}

func (mainPage *MainPage) Init() {
	mainPage.Refresh()
}

func (mainPage *MainPage) Refresh() {
	defer mainPage.application.ForceDraw()
	mainPage.UpdateHeader()

	err := mainPage.sensorsPage.Refresh()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return
	}
	err = mainPage.curvesPage.Refresh()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return
	}
	err = mainPage.fansPage.Refresh()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return
	}

	mainPage.clearStatusMessage()
}

func (mainPage *MainPage) SetPage(page Page) {
	mainPage.page = page
	mainPage.header.SetPage(page)

	mainPage.mainPagePagerLayout.SwitchToPage(string(page))
	mainPage.Refresh()

	switch mainPage.page {
	case FansPage:
		mainPage.fansPage.ScrollToItem()
	case CurvesPage:
		mainPage.curvesPage.ScrollToItem()
	case SensorsPage:
		mainPage.sensorsPage.ScrollToItem()
	}
}

func (mainPage *MainPage) PreviousPage() {
	currentIndex := slices.Index(Pages, mainPage.page)
	nextIndex := (len(Pages) + currentIndex - 1) % len(Pages)
	mainPage.SetPage(Pages[nextIndex])
}

func (mainPage *MainPage) NextPage() {
	currentIndex := slices.Index(Pages, mainPage.page)
	nextIndex := (currentIndex + 1) % len(Pages)
	mainPage.SetPage(Pages[nextIndex])
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
