package ui

import (
	"context"
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/dialog"
	"fan2go-tui/internal/ui/shortcut_helper"
	"fan2go-tui/internal/ui/status_message"
	"fan2go-tui/internal/ui/util"
	"fan2go-tui/internal/uistate"
	"fmt"
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

// CreateUi creates the application, restoring the last session with the configured fan2go daemon (see
// configuration.UiConfig.RestoreSession), and remembering the session in uistate.Current. The fan2go daemon is polled
// until ctx is done; cancel it when the application stopped.
func CreateUi(ctx context.Context, fullscreen bool) *tview.Application {
	application, _ := createUi(ctx, fullscreen)
	return application
}

// createUi creates the application, see CreateUi, and returns its main page as well (for tests).
func createUi(ctx context.Context, fullscreen bool) (*tview.Application, *MainPage) {
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
	mainPage.PersistSession(uistate.Current, serverAddress, configuration.CurrentConfig.Ui.RestoreSession)
	if version, newer := uistate.Current.NewerVersion(); newer {
		mainPage.header.SetDefaultStatus(newerStateVersionMessage(version))
	}
	connectingDialog := dialog.NewConnectingDialog(serverAddress)

	hasReceivedData := false
	dismissedConnectingDialog := false

	var pagesLayout *tview.Pages

	poller := state.NewPoller(apiClient, store, func(err error) {
		queueUpdateDraw(ctx, application, func() {
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
				if frontPage, _ := pagesLayout.GetFrontPage(); frontPage == string(ConnectingDialog) {
					pagesLayout.HidePage(string(ConnectingDialog))
				}
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

	// the ticker of this UI, UpdateTicker is replaced by the next call (in tests)
	ticker := UpdateTicker
	go func() {
		defer ticker.Stop()
		// Initial fetch
		poller.FetchAndUpdate()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				poller.FetchAndUpdate()
			}
		}
	}()

	return application.SetRoot(pagesLayout, fullscreen), mainPage
}

// queueUpdateDraw runs f on the UI thread and waits for it, like application.QueueUpdateDraw, but only until ctx is
// done: once the application stopped, it runs no queued updates anymore, and QueueUpdateDraw would wait forever. The
// call itself is left waiting then (one goroutine, until the process exits).
func queueUpdateDraw(ctx context.Context, application *tview.Application, f func()) {
	if ctx.Err() != nil {
		return
	}
	done := make(chan struct{})
	go func() {
		application.QueueUpdateDraw(f)
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
	}
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

// newerStateVersionMessage is the warning shown if the remembered state was saved by a newer version of the
// application, see uistate.Store.NewerVersion.
func newerStateVersionMessage(version int) *status_message.StatusMessage {
	return status_message.NewWarningStatusMessage(fmt.Sprintf(
		"Remembered settings are from a newer fan2go-tui (format %d) and are not saved. Please update fan2go-tui.",
		version))
}
