package ui

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/ui/dialog"
	"fan2go-tui/internal/ui/util"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

const (
	Main       util.Page = "main"
	HelpDialog util.Page = "help"
)

var (
	UpdateTicker *time.Ticker
)

func CreateUi(fullscreen bool) *tview.Application {
	UpdateTicker = time.NewTicker(configuration.CurrentConfig.Ui.UpdateInterval)

	application := tview.NewApplication()
	application.EnableMouse(true)

	baseUrl := configuration.CurrentConfig.Api.Host
	port := configuration.CurrentConfig.Api.Port
	apiClient := client.NewApiClient(baseUrl, port)

	mainPage := NewMainPage(application, apiClient)
	helpPage := dialog.NewHelpPage()

	pagesLayout := tview.NewPages().
		AddPage(string(Main), mainPage.layout, true, true).
		AddPage(string(HelpDialog), helpPage.GetLayout(), true, false)

	pagesLayout.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// ignore events, if some other page is open
		//name, _ := pagesLayout.GetFrontPage()
		//fileBrowserPage, _ := mainPage.fileBrowser.GetLayout().GetFrontPage()
		//if name != string(Main) || fileBrowserPage != string(file_browser.FileBrowserPage) {
		//	return event
		//}

		if event.Key() == tcell.KeyCtrlC || event.Key() == tcell.KeyCtrlQ {
			application.Stop()
			return nil
		} else if event.Rune() == '?' || event.Key() == tcell.KeyF1 {
			pagesLayout.ShowPage(string(HelpDialog))
			return nil
		} else if event.Rune() == '+' {
			configuration.CurrentConfig.Ui.UpdateInterval += 100 * time.Millisecond
			UpdateTicker.Reset(configuration.CurrentConfig.Ui.UpdateInterval)
			mainPage.UpdateHeader()
			return nil
		} else if event.Rune() == '-' {
			stepSize := 100 * time.Millisecond
			if configuration.CurrentConfig.Ui.UpdateInterval <= stepSize {
				configuration.CurrentConfig.Ui.UpdateInterval = stepSize
			} else {
				configuration.CurrentConfig.Ui.UpdateInterval -= stepSize
			}
			UpdateTicker.Reset(configuration.CurrentConfig.Ui.UpdateInterval)
			mainPage.UpdateHeader()
			return nil
		}
		return event
	})

	helpPage.GetLayout().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape {
			pagesLayout.HidePage(string(HelpDialog))
			return nil
		}
		return event
	})

	mainPage.Init()

	go func() {
		for {
			<-UpdateTicker.C
			mainPage.Refresh()
			application.Draw()
		}
	}()

	return application.SetRoot(pagesLayout, fullscreen)
}
