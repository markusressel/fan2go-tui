package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/fan"
	"fan2go-tui/internal/ui/status_message"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainPage struct {
	application *tview.Application

	client client.Fan2goApiClient

	layout               *tview.Flex
	header               *ApplicationHeaderComponent
	fanComponents        []*fan.FanComponent
	fanOverviewComponent *fan.FanOverviewComponent
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

	fanOverviewComponent := fan.NewFanOverviewComponent(mainPage.application)
	windowLayout.AddItem(fanOverviewComponent.GetLayout(), 0, 3, true)
	mainPage.fanOverviewComponent = fanOverviewComponent

	fans := mainPage.client.GetFans()
	var fanComponents []*fan.FanComponent
	for _, f := range fans {
		fanComponent := fan.NewFanComponent(mainPage.application, f)
		fanComponents = append(fanComponents, fanComponent)
		fanComponent.Refresh()
		layout := fanComponent.GetLayout()
		windowLayout.AddItem(layout, 0, 1, true)
	}
	mainPage.fanComponents = fanComponents

	infoLayout := tview.NewFlex().SetDirection(tview.FlexRow)
	windowLayout.AddItem(infoLayout, 0, 1, false)

	mainPageLayout.AddItem(windowLayout, 0, 1, true)

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
		fan := mainPage.client.GetFan(component.Fan.Label)
		component.SetFan(fan)
		component.Refresh()
	}
}

func (mainPage *MainPage) ToggleFocus() {

}

func (mainPage *MainPage) showStatusMessage(status *status_message.StatusMessage) {
	mainPage.header.SetStatus(status)
}
