package ui

import (
	"fan2go-tui/cmd/global"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/ui/shortcut_helper"
	"fan2go-tui/internal/ui/status_message"
	"fan2go-tui/internal/ui/theme"
	uiutil "fan2go-tui/internal/ui/util"
	"fmt"
	"strings"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ApplicationHeaderComponent struct {
	application *tview.Application

	name    string
	version string

	layout                *tview.Flex
	statusTextView        *tview.TextView
	pageIndicatorTextView *tview.TextView

	lastStatus             *status_message.StatusMessage
	pagesMap               orderedmap.OrderedMap[Page, uiutil.PagesPage]
	page                   Page
	updateIntervalTextView *tview.TextView

	connectionStatusTextView *tview.TextView
	connected                bool
	shortcutHint             *tview.TextView
}

func NewApplicationHeader(
	application *tview.Application,
	pagesMap orderedmap.OrderedMap[Page, uiutil.PagesPage],
) *ApplicationHeaderComponent {
	versionText := fmt.Sprintf("%s-(#%s)-%s", global.Version, global.Commit, global.Date)

	applicationHeader := &ApplicationHeaderComponent{
		application: application,
		name:        "fan2go-tui",
		version:     versionText,
		pagesMap:    pagesMap,
		page:        pagesMap.Keys()[0],
		connected:   true,
	}

	applicationHeader.createLayout()
	applicationHeader.updateUi()

	return applicationHeader
}

func (applicationHeader *ApplicationHeaderComponent) createLayout() {
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)

	nameTextView := tview.NewTextView()
	nameTextView.SetTextColor(theme.Colors.Header.Name)
	nameTextView.SetBackgroundColor(theme.Colors.Header.NameBackground)
	nameText := fmt.Sprintf(" %s ", applicationHeader.name)
	nameTextView.SetText(nameText)
	nameTextView.SetTextAlign(tview.AlignCenter)

	versionTextView := tview.NewTextView()
	versionTextView.SetBackgroundColor(theme.Colors.Header.VersionBackground)
	versionTextView.SetTextColor(theme.Colors.Header.Version)
	versionText := fmt.Sprintf("  %s  ", applicationHeader.version)
	versionTextView.SetText(versionText)
	versionTextView.SetTextAlign(tview.AlignCenter)

	connectionStatusTextView := tview.NewTextView().
		SetTextAlign(tview.AlignCenter)
	applicationHeader.connectionStatusTextView = connectionStatusTextView

	statusTextView := tview.NewTextView()
	statusTextView.SetBorderPadding(0, 0, 1, 1)
	statusTextView.SetTextColor(tcell.ColorGray)
	statusTextView.SetTextAlign(tview.AlignLeft)

	applicationHeader.updateIntervalTextView = tview.NewTextView()
	applicationHeader.updateIntervalTextView.SetBackgroundColor(theme.Colors.Header.UpdateIntervalBackground)
	applicationHeader.updateIntervalTextView.SetTextColor(theme.Colors.Header.UpdateInterval)
	applicationHeader.updateIntervalTextView.SetTextAlign(tview.AlignCenter)

	page := applicationHeader.page
	pageName := string(applicationHeader.page)
	pageIdx := indexOf(page, applicationHeader.pagesMap.Keys())
	pageCount := len(applicationHeader.pagesMap.Keys())
	pageIndicatorText := fmt.Sprintf("%-7s %d/%d", pageName, pageIdx+1, pageCount)

	pageIndicatorTextView := tview.NewTextView()
	pageIndicatorTextView.SetText(pageIndicatorText).
		SetTextColor(theme.Colors.Header.PageIndicator).
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(theme.Colors.Header.PageIndicatorBackground)

	applicationHeader.pageIndicatorTextView = pageIndicatorTextView

	shortcutHint := uiutil.CreateAttentionTextView(shortcutHintText)

	layout.AddItem(nameTextView, len(nameText), 0, false)
	layout.AddItem(versionTextView, len(versionText), 0, false)
	layout.AddItem(connectionStatusTextView, 0, 0, false)
	layout.AddItem(statusTextView, 0, 1, false)
	layout.AddItem(pageIndicatorTextView, len(pageIndicatorText)+4, 0, false)
	layout.AddItem(applicationHeader.updateIntervalTextView, len(pageIndicatorText)+4, 0, false)
	layout.AddItem(shortcutHint, 0, 0, false)

	applicationHeader.shortcutHint = shortcutHint
	applicationHeader.statusTextView = statusTextView
	applicationHeader.layout = layout

	applicationHeader.updateUi()
	applicationHeader.UpdateShortcutHint()
}

const shortcutHintText = "? shortcuts"

func (applicationHeader *ApplicationHeaderComponent) UpdateShortcutHint() {
	width := 0
	if shortcut_helper.ShortcutsHidden() {
		width = len(shortcutHintText) + 4
	}
	applicationHeader.layout.ResizeItem(applicationHeader.shortcutHint, width, 0)
}

func (applicationHeader *ApplicationHeaderComponent) updateUi() {
	page := applicationHeader.page
	pageName := strings.ToUpper(string(applicationHeader.page))
	pageIdx := indexOf(page, applicationHeader.pagesMap.Keys())
	pageCount := len(applicationHeader.pagesMap.Keys())
	pageIndicatorText := fmt.Sprintf("%-7s %d/%d", pageName, pageIdx+1, pageCount)
	applicationHeader.pageIndicatorTextView.SetText(pageIndicatorText)
}

func (applicationHeader *ApplicationHeaderComponent) SetConnectionStatus(connected bool, message string) {
	applicationHeader.connected = connected
	if connected {
		applicationHeader.connectionStatusTextView.SetText("")
		applicationHeader.layout.ResizeItem(applicationHeader.connectionStatusTextView, 0, 0)
		applicationHeader.ResetStatus()
	} else {
		badgeText := " DISCONNECTED "
		applicationHeader.connectionStatusTextView.
			SetText(badgeText).
			SetTextColor(theme.Colors.Header.DisconnectedBadge).
			SetBackgroundColor(theme.Colors.Header.DisconnectedBadgeBackground)
		applicationHeader.layout.ResizeItem(applicationHeader.connectionStatusTextView, len(badgeText), 0)

		if message != "" {
			applicationHeader.statusTextView.SetText(message).SetTextColor(tcell.ColorRed)
			applicationHeader.lastStatus = status_message.NewErrorStatusMessage(message)
		}
	}
	applicationHeader.application.ForceDraw()
}

func (applicationHeader *ApplicationHeaderComponent) SetStatus(status *status_message.StatusMessage) {
	if !applicationHeader.connected && status.Message == "" {
		return
	}
	applicationHeader.statusTextView.SetText(status.Message).SetTextColor(status.Color)
	applicationHeader.application.ForceDraw()
	if status.Duration > 0 {
		go func() {
			time.Sleep(status.Duration)
			applicationHeader.application.QueueUpdateDraw(func() {
				if applicationHeader.lastStatus != status {
					return
				}
				applicationHeader.ResetStatus()
			})
		}()
	}
	applicationHeader.lastStatus = status
}

func (applicationHeader *ApplicationHeaderComponent) ResetStatus() {
	if !applicationHeader.connected {
		return
	}
	applicationHeader.statusTextView.SetText("").SetTextColor(tcell.ColorWhite)
	applicationHeader.application.ForceDraw()
}

func (applicationHeader *ApplicationHeaderComponent) SetPage(page Page) {
	applicationHeader.page = page
	applicationHeader.updateUi()
}

func (applicationHeader *ApplicationHeaderComponent) Refresh() {
	millis := configuration.CurrentConfig.Ui.UpdateInterval.Milliseconds()
	updateIntervalText := fmt.Sprintf("- %d ms +", millis)
	applicationHeader.updateIntervalTextView.SetText(updateIntervalText)
}

func indexOf[T comparable](word T, data []T) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}
