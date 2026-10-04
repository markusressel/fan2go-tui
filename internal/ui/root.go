package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/dialog"
	"fan2go-tui/internal/ui/shortcut_helper"
	"fan2go-tui/internal/ui/util"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	Main             util.Page = "main"
	ConnectingDialog util.Page = "connecting"
)

var (
	UpdateTicker *time.Ticker

	UpdateIntervalStepSize = 100 * time.Millisecond
)

func CreateUi(fullscreen bool) *tview.Application {
	shortcut_helper.InitShortcutVisibility()
	UpdateTicker = time.NewTicker(configuration.CurrentConfig.Ui.UpdateInterval)

	application := tview.NewApplication()
	application.EnableMouse(true)

	host := configuration.CurrentConfig.Api.Host
	port := configuration.CurrentConfig.Api.Port
	apiClient := client.NewApiClient(host, port)

	store := state.NewStore()

	mainPage := NewMainPage(application, store)

	serverAddress := client.JoinHostPort(host, port)
	connectingDialog := dialog.NewConnectingDialog(serverAddress)

	hasReceivedData := false
	dismissedConnectingDialog := false

	var pagesLayout *tview.Pages

	poller := state.NewPoller(apiClient, store, func(err error) {
		application.QueueUpdateDraw(func() {
			if err != nil {
				if !hasReceivedData && !store.HasData() {
					connectingDialog.SetError(err)
					if !dismissedConnectingDialog {
						pagesLayout.ShowPage(string(ConnectingDialog))
					} else {
						mainPage.SetConnectionStatus(false, err.Error())
					}
				} else {
					mainPage.SetConnectionStatus(false, err.Error())
				}
			} else {
				hasReceivedData = true
				pagesLayout.HidePage(string(ConnectingDialog))
				mainPage.SetConnectionStatus(true, "")
			}
			mainPage.Refresh()
		})
	})

	pagesLayout = tview.NewPages().
		AddPage(string(Main), mainPage.layout, true, true).
		AddPage(string(ConnectingDialog), connectingDialog.GetLayout(), true, true)

	pagesLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlQ {
			application.Stop()
			return nil
		}

		frontPage, _ := pagesLayout.GetFrontPage()
		switch frontPage {
		case string(ConnectingDialog):
			if event.Key() == tcell.KeyEscape {
				dismissedConnectingDialog = true
				pagesLayout.HidePage(string(ConnectingDialog))
				if connectingDialog.LastError() != nil {
					mainPage.SetConnectionStatus(false, connectingDialog.LastError().Error())
				}
				application.SetFocus(mainPage.layout)
				return nil
			}
			return event
		}

		if (event.Key() == tcell.KeyRune && event.Rune() == '?' && !util.IsTextInputActive(application.GetFocus())) ||
			event.Key() == tcell.KeyF1 {
			shortcut_helper.ToggleShortcuts()
			mainPage.header.UpdateShortcutHint()
			return nil
		} else if event.Rune() == '+' {
			slowDownUpdateInterval(mainPage)
			return nil
		} else if event.Rune() == '-' {
			speedUpUpdateInterval(mainPage)
			return nil
		} else if event.Modifiers() == tcell.ModNone && event.Key() == tcell.KeyBacktab {
			mainPage.PreviousPage()
			return nil
		} else if event.Modifiers() == tcell.ModNone && event.Key() == tcell.KeyTab {
			mainPage.NextPage()
			return nil
		}
		return event
	})

	connectingDialog.GetLayout().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			dismissedConnectingDialog = true
			pagesLayout.HidePage(string(ConnectingDialog))
			if connectingDialog.LastError() != nil {
				mainPage.SetConnectionStatus(false, connectingDialog.LastError().Error())
			}
			application.SetFocus(mainPage.layout)
			return nil
		} else if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlQ {
			application.Stop()
			return nil
		}
		return event
	})

	go func() {
		application.QueueUpdateDraw(func() {
			mainPage.Init()
		})
	}()

	go func() {
		// Initial fetch
		poller.FetchAndUpdate()
		for {
			<-UpdateTicker.C
			poller.FetchAndUpdate()
		}
	}()

	return application.SetRoot(pagesLayout, fullscreen)
}

func speedUpUpdateInterval(mainPage *MainPage) {
	if configuration.CurrentConfig.Ui.UpdateInterval <= UpdateIntervalStepSize {
		configuration.CurrentConfig.Ui.UpdateInterval = UpdateIntervalStepSize
	} else {
		configuration.CurrentConfig.Ui.UpdateInterval -= UpdateIntervalStepSize
	}
	UpdateTicker.Reset(configuration.CurrentConfig.Ui.UpdateInterval)
	mainPage.UpdateHeader()
}

func slowDownUpdateInterval(mainPage *MainPage) {
	configuration.CurrentConfig.Ui.UpdateInterval += UpdateIntervalStepSize
	UpdateTicker.Reset(configuration.CurrentConfig.Ui.UpdateInterval)
	mainPage.UpdateHeader()
}
