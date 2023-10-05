package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/curve"
	"fan2go-tui/internal/ui/fan"
	"fan2go-tui/internal/ui/sensor"
	"fan2go-tui/internal/ui/status_message"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sort"
	"strings"
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

	fanComponents    []*fan.FanComponent
	curveComponents  []*curve.CurveComponent
	sensorComponents []*sensor.SensorComponent

	fanOverviewComponent  *fan.FanOverviewComponent
	curveGraphsComponent  *curve.CurveGraphsComponent
	sensorGraphsComponent *sensor.SensorGraphsComponent

	page                Page
	mainPagePagerLayout *tview.Pages
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

	fansPageLayout := createFansPageLayout(mainPage)
	curvesPageLayout := createCurvesPageLayout(mainPage)
	sensorsPageLayout := createSensorsPageLayout(mainPage)

	mainPagePagerLayout.AddPage(string(FansPage), fansPageLayout, true, true)
	mainPagePagerLayout.AddPage(string(CurvesPage), curvesPageLayout, true, false)
	mainPagePagerLayout.AddPage(string(SensorsPage), sensorsPageLayout, true, false)

	return mainPageLayout
}

func createSensorsPageLayout(mainPage *MainPage) *tview.Flex {
	sensorsPageLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	sensorInfoLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	sensorsPageLayout.AddItem(sensorInfoLayout, 0, 1, true)

	sensors, err := mainPage.client.GetSensors()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return sensorsPageLayout
	}
	var sensorComponents []*sensor.SensorComponent
	for _, s := range *sensors {
		sensorComponent := sensor.NewSensorComponent(mainPage.application, s)
		sensorComponents = append(sensorComponents, sensorComponent)
		sensorComponent.SetSensor(s)
		sensorComponent.Refresh()
		layout := sensorComponent.GetLayout()
		sensorInfoLayout.AddItem(layout, 0, 1, true)
	}
	mainPage.sensorComponents = sensorComponents

	sensorGraphsComponent := sensor.NewSensorGraphsComponent(mainPage.application)
	mainPage.sensorGraphsComponent = sensorGraphsComponent
	// sensorComponents = append(sensorComponents, sensorGaphsComponent)

	// update overview
	sensorList := []*client.Sensor{}
	for _, f := range *sensors {
		sensorList = append(sensorList, f)
	}

	sensorGraphsComponent.SetSensors(sensorList)
	sensorGraphsComponent.Refresh()
	layout := sensorGraphsComponent.GetLayout()
	sensorsPageLayout.AddItem(layout, 0, 1, true)

	return sensorsPageLayout
}

func createCurvesPageLayout(mainPage *MainPage) *tview.Flex {
	curvesPageLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	curveInfoLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	curvesPageLayout.AddItem(curveInfoLayout, 0, 1, true)

	curves, err := mainPage.client.GetCurves()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return curvesPageLayout
	}
	var curveComponents []*curve.CurveComponent
	curvesIds := []string{}
	for _, c := range *curves {
		curvesIds = append(curvesIds, c.Config.ID)
	}

	sort.SliceStable(curvesIds, func(i, j int) bool {
		a := curvesIds[i]
		b := curvesIds[j]

		result := strings.Compare(strings.ToLower(a), strings.ToLower(b))

		if result <= 0 {
			return true
		} else {
			return false
		}
	})

	for _, id := range curvesIds {
		c := (*curves)[id]

		curveComponent := curve.NewCurveComponent(mainPage.application, c)
		curveComponents = append(curveComponents, curveComponent)
		curveComponent.SetCurve(c)
		curveComponent.Refresh()
		layout := curveComponent.GetLayout()
		curveInfoLayout.AddItem(layout, 0, 1, true)
	}
	mainPage.curveComponents = curveComponents

	curveGraphsComponent := curve.NewCurveGraphsComponent(mainPage.application)
	mainPage.curveGraphsComponent = curveGraphsComponent
	//curveComponents = append(curveComponents, curveGaphsComponent)

	// update overview
	curveList := []*client.Curve{}
	for _, f := range *curves {
		curveList = append(curveList, f)
	}

	curveGraphsComponent.SetCurves(curveList)
	curveGraphsComponent.Refresh()
	layout := curveGraphsComponent.GetLayout()
	curvesPageLayout.AddItem(layout, 0, 1, true)

	return curvesPageLayout
}

func createFansPageLayout(mainPage *MainPage) *tview.Flex {
	splitLayout := tview.NewFlex().SetDirection(tview.FlexColumn)

	infosList := tview.NewFlex().SetDirection(tview.FlexRow)
	splitLayout.AddItem(infosList, 0, 1, true)

	fanOverviewComponent := fan.NewFanOverviewComponent(mainPage.application)
	splitLayout.AddItem(fanOverviewComponent.GetLayout(), 0, 3, true)
	mainPage.fanOverviewComponent = fanOverviewComponent

	fans, err := mainPage.client.GetFans()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return splitLayout
	}
	var fanComponents []*fan.FanComponent
	for _, f := range *fans {
		fanComponent := fan.NewFanComponent(mainPage.application, f)
		fanComponents = append(fanComponents, fanComponent)
		fanComponent.Refresh()
		layout := fanComponent.GetLayout()
		infosList.AddItem(layout, 0, 1, true)
	}
	mainPage.fanComponents = fanComponents

	return splitLayout
}

func (mainPage *MainPage) Init() {
	fans, err := mainPage.client.GetFans()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return
	}

	// update overview
	fanList := []*client.Fan{}
	for _, f := range *fans {
		fanList = append(fanList, f)
	}
	mainPage.fanOverviewComponent.SetFans(fanList)

	// update details
	mainPage.Refresh()
}

func (mainPage *MainPage) Refresh() {
	// always update fans, to get the latest values while on other pages
	fans, err := mainPage.client.GetFans()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return
	}
	// update overview
	fanList := []*client.Fan{}
	for _, f := range *fans {
		fanList = append(fanList, f)
	}
	mainPage.fanOverviewComponent.SetFans(fanList)

	for _, component := range mainPage.fanComponents {
		f, err := mainPage.client.GetFan(component.Fan.Config.Id)
		if err != nil {
			mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
			continue
		}
		component.SetFan(f)
		component.Refresh()
	}

	curves, err := mainPage.client.GetCurves()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return
	}
	// update overview
	curveList := []*client.Curve{}
	for _, f := range *curves {
		curveList = append(curveList, f)
	}
	mainPage.curveGraphsComponent.SetCurves(curveList)

	for _, component := range mainPage.curveComponents {
		c, err := mainPage.client.GetCurve(component.Curve.Config.ID)
		if err != nil {
			mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
			continue
		}
		component.SetCurve(c)
		component.Refresh()
	}

	sensors, err := mainPage.client.GetSensors()
	if err != nil {
		mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
		return
	}
	// update overview
	sensorList := []*client.Sensor{}
	for _, f := range *sensors {
		sensorList = append(sensorList, f)
	}
	mainPage.sensorGraphsComponent.SetSensors(sensorList)

	for _, component := range mainPage.sensorComponents {
		s, err := mainPage.client.GetSensor(component.Sensor.Config.ID)
		if err != nil {
			mainPage.showStatusMessage(status_message.NewErrorStatusMessage(err.Error()))
			continue
		}
		component.SetSensor(s)
		component.Refresh()
	}

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
