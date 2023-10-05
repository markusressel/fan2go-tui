package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/curve"
	"fan2go-tui/internal/ui/fan"
	"fan2go-tui/internal/ui/sensor"
	"fan2go-tui/internal/ui/status_message"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
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
			mainPage.ToggleFocus()
		} else if key == tcell.KeyCtrlR {

		} else if rune == int32(49) {
			page := Pages[0]
			mainPage.header.SetPage(page)
			mainPage.SetPage(page)
		} else if rune == int32(50) {
			page := Pages[1]
			mainPage.header.SetPage(page)
			mainPage.SetPage(page)
		} else if rune == int32(51) {
			page := Pages[2]
			mainPage.header.SetPage(page)
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
	mainPage.fansPage.Refresh()
	mainPage.curvesPage.Refresh()
	mainPage.sensorsPage.Refresh()

	//curves, err := mainPage.client.GetCurves()
	//if err != nil {
	//	mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
	//	return
	//}
	//// update overview
	//curveList := []*client.Curve{}
	//for _, f := range *curves {
	//	curveList = append(curveList, f)
	//}
	//mainPage.curveGraphsComponent.SetCurves(curveList)

	//for _, component := range mainPage.curveComponents {
	//	c, err := mainPage.client.GetCurve(component.Curve.Config.ID)
	//	if err != nil {
	//		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
	//		continue
	//	}
	//	component.SetCurve(c)
	//	component.Refresh()
	//}

	//sensors, err := mainPage.client.GetSensors()
	//if err != nil {
	//	mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
	//	return
	//}
	//// update overview
	//sensorList := []*client.Sensor{}
	//for _, f := range *sensors {
	//	sensorList = append(sensorList, f)
	//}
	//mainPage.sensorGraphsComponent.SetSensors(sensorList)

	//for _, component := range mainPage.sensorComponents {
	//	s, err := mainPage.client.GetSensor(component.Sensor.Config.ID)
	//	if err != nil {
	//		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
	//		continue
	//	}
	//	component.SetSensor(s)
	//	component.Refresh()
	//}

	mainPage.application.ForceDraw()
}

func (mainPage *MainPage) ToggleFocus() {

}

func (mainPage *MainPage) SetPage(page Page) {
	mainPage.page = page
	mainPage.mainPagePagerLayout.SwitchToPage(string(page))
	mainPage.Refresh()
}

func (mainPage *MainPage) showStatusMessage(status *status_message.StatusMessage) {
	mainPage.header.SetStatus(status)
}
