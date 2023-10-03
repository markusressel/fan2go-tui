package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/fan"
	"fan2go-tui/internal/ui/status_message"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainPage struct {
	application   *tview.Application
	header        *ApplicationHeaderComponent
	layout        *tview.Flex
	client        client.Fan2goApiClient
	fanComponents []*fan.FanComponent
}

func NewMainPage(application *tview.Application, client client.Fan2goApiClient) *MainPage {

	mainPage := &MainPage{
		application: application,
		client:      client,
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

	windowLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	//dialog := createFileBrowserActionDialog()

	fans := mainPage.client.GetFans()
	for _, f := range fans {
		var fanComponents []*fan.FanComponent
		fanComponent := fan.NewFanComponent(mainPage.application, &f)
		fanComponents = append(fanComponents, fanComponent)

		layout := fanComponent.GetLayout()
		windowLayout.AddItem(layout, 0, 1, true)

		mainPage.fanComponents = fanComponents
	}

	infoLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	windowLayout.AddItem(infoLayout, 0, 1, false)

	mainPageLayout.AddItem(windowLayout, 0, 1, true)

	mainPage.header = header

	return mainPageLayout
}

func (mainPage *MainPage) Init() {
	//fanData := mainPage.client.GetFans()
	//curveData := mainPage.client.GetCurves()
	//sensorData := mainPage.client.GetSensors()
	//
	//var data []map[string]interface{}
	//data = append(data, fanData)
	//data = append(data, curveData)
	//data = append(data, sensorData)
	//
	//for _, item := range data {
	//	var keys []string
	//	for pwm := range item {
	//		keys = append(keys, pwm)
	//	}
	//	sort.Strings(keys)
	//
	//	text := fmt.Sprintf("%v", keys)
	//
	//	mainPage.showStatusMessage(status_message.NewWarningStatusMessage(text))
	//}
}

func (mainPage *MainPage) ToggleFocus() {

}

func (mainPage *MainPage) showStatusMessage(status *status_message.StatusMessage) {
	mainPage.header.SetStatus(status)
}
