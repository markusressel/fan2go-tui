package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/fan"
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

type MainPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	layout               *tview.Flex
	header               *ApplicationHeaderComponent
	fanComponents        []*fan.FanComponent
	curveComponents      []*fan.CurveComponent
	sensorComponents     []*fan.SensorComponent
	fanOverviewComponent *fan.FanOverviewComponent
	page                 Page
}

func NewMainPage(application *tview.Application, client client.Fan2goApiClient) *MainPage {

	mainPage := &MainPage{
		application: application,
		client:      client,
		page:        FansPage,
	}

	//fileBrowser.SetStatusCallback(func(message *status_message.StatusMessage) {
	//	mainPage.showStatusMessage(message)
	//})

	mainPage.layout = mainPage.createLayout()
	mainPage.layout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key == tcell.KeyTab || key == tcell.KeyBacktab {
			mainPage.ToggleFocus()
		} else if key == tcell.KeyCtrlR {

		}
		return event
	})

	return mainPage
}

func (mainPage *MainPage) createLayout() *tview.Flex {
	mainPageLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	header := NewApplicationHeader(mainPage.application)
	mainPageLayout.AddItem(header.layout, 1, 0, false)

	splitLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	infosList := tview.NewFlex().SetDirection(tview.FlexRow)
	splitLayout.AddItem(infosList, 0, 1, true)

	fanOverviewComponent := fan.NewFanOverviewComponent(mainPage.application)
	splitLayout.AddItem(fanOverviewComponent.GetLayout(), 0, 3, true)
	mainPage.fanOverviewComponent = fanOverviewComponent

	// fans
	if mainPage.page == FansPage {
		fans := mainPage.client.GetFans()
		var fanComponents []*fan.FanComponent
		for _, f := range fans {
			fanComponent := fan.NewFanComponent(mainPage.application, f)
			fanComponents = append(fanComponents, fanComponent)
			fanComponent.Refresh()
			layout := fanComponent.GetLayout()
			infosList.AddItem(layout, 0, 1, true)
		}
		mainPage.fanComponents = fanComponents
	}

	// curves
	if mainPage.page == CurvesPage {
		curves := mainPage.client.GetCurves()
		var curveComponents []*fan.CurveComponent
		for _, c := range curves {
			curveComponent := fan.NewCurveComponent(mainPage.application, c)
			curveComponents = append(curveComponents, curveComponent)
			curveComponent.SetCurve(c)
			curveComponent.Refresh()
			layout := curveComponent.GetLayout()
			infosList.AddItem(layout, 0, 1, true)
		}
		mainPage.curveComponents = curveComponents
	}

	// sensors
	if mainPage.page == SensorsPage {
		sensors := mainPage.client.GetSensors()
		var sensorComponents []*fan.SensorComponent
		for _, s := range sensors {
			sensorComponent := fan.NewSensorComponent(mainPage.application, s)
			sensorComponents = append(sensorComponents, sensorComponent)
			sensorComponent.SetSensor(s)
			sensorComponent.Refresh()
			layout := sensorComponent.GetLayout()
			infosList.AddItem(layout, 0, 1, true)
		}
		mainPage.sensorComponents = sensorComponents
	}

	mainPageLayout.AddItem(splitLayout, 0, 1, true)

	mainPage.header = header

	return mainPageLayout
}

func (mainPage *MainPage) Init() {
	fans := mainPage.client.GetFans()

	// update overview
	fanList := []*client.Fan{}
	for _, f := range fans {
		fanList = append(fanList, f)
	}
	mainPage.fanOverviewComponent.SetFans(fanList)

	// update details
	mainPage.Refresh()
}

func (mainPage *MainPage) Refresh() {
	fans := mainPage.client.GetFans()
	// update overview
	fanList := []*client.Fan{}
	for _, f := range fans {
		fanList = append(fanList, f)
	}
	mainPage.fanOverviewComponent.SetFans(fanList)

	for _, component := range mainPage.fanComponents {
		f := mainPage.client.GetFan(component.Fan.Config.Id)
		component.SetFan(f)
		component.Refresh()
	}

	for _, component := range mainPage.curveComponents {
		curve := mainPage.client.GetCurve(component.Curve.Config.ID)
		component.SetCurve(curve)
		component.Refresh()
	}

	for _, component := range mainPage.sensorComponents {
		sensor := mainPage.client.GetSensor(component.Sensor.Config.ID)
		component.SetSensor(sensor)
		component.Refresh()
	}
}

func (mainPage *MainPage) ToggleFocus() {

}

func (mainPage *MainPage) SetPage(page Page) {
	mainPage.page = page
}

func (mainPage *MainPage) showStatusMessage(status *status_message.StatusMessage) {
	mainPage.header.SetStatus(status)
}
