package ui

import (
	"fan2go-tui/cmd/global"
	"fan2go-tui/internal/ui/status_message"
	uiutil "fan2go-tui/internal/ui/util"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"time"
)

type ApplicationHeaderComponent struct {
	application *tview.Application

	name    string
	version string

	layout                *tview.Flex
	statusTextView        *tview.TextView
	pageIndicatorTextView *tview.TextView

	lastStatus *status_message.StatusMessage
	page       Page
}

func NewApplicationHeader(application *tview.Application) *ApplicationHeaderComponent {
	versionText := fmt.Sprintf("%s-(#%s)-%s", global.Version, global.Commit, global.Date)

	applicationHeader := &ApplicationHeaderComponent{
		application: application,
		name:        "fan2go-tui",
		version:     versionText,
		page:        FansPage,
	}

	applicationHeader.createLayout()
	applicationHeader.updateUi()

	return applicationHeader
}

func (applicationHeader *ApplicationHeaderComponent) createLayout() {
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)
	// TODO: check colors
	layout.SetBackgroundColor(tcell.ColorRed)
	layout.SetTitleColor(tcell.ColorRed)
	layout.SetBorderColor(tcell.ColorGreen)

	nameTextView := tview.NewTextView()
	nameTextView.SetTextColor(tcell.ColorWhite)
	nameTextView.SetBackgroundColor(tcell.ColorDodgerBlue)
	nameText := fmt.Sprintf(" %s ", applicationHeader.name)
	nameTextView.SetText(nameText)
	nameTextView.SetTextAlign(tview.AlignCenter)

	versionTextView := tview.NewTextView()
	versionTextView.SetBackgroundColor(tcell.ColorGreenYellow)
	versionTextView.SetTextColor(tcell.ColorBlack)
	versionText := fmt.Sprintf("  %s  ", applicationHeader.version)
	versionTextView.SetText(versionText)
	versionTextView.SetTextAlign(tview.AlignCenter)

	statusTextView := tview.NewTextView()
	statusTextView.SetBorderPadding(0, 0, 1, 1)
	statusTextView.SetTextColor(tcell.ColorGray)
	statusTextView.SetTextAlign(tview.AlignLeft)

	page := applicationHeader.page
	pageName := string(applicationHeader.page)
	pageIdx := indexOf(page, Pages)
	pageCount := len(Pages)
	pageIndicatorText := fmt.Sprintf("Page: %s %d/%d", pageName, pageIdx+1, pageCount)

	pageIndicatorTextView := tview.NewTextView()
	pageIndicatorTextView.SetText(pageIndicatorText).
		SetTextColor(tcell.ColorWhite).
		SetTextAlign(tview.AlignCenter).
		SetBackgroundColor(tcell.ColorBlue)

	applicationHeader.pageIndicatorTextView = pageIndicatorTextView

	helpText := "Press '?' for help"
	helpTextView := uiutil.CreateAttentionTextView(helpText)

	layout.AddItem(nameTextView, len(nameText), 0, false)
	layout.AddItem(versionTextView, len(versionText), 0, false)
	layout.AddItem(statusTextView, 0, 1, false)
	layout.AddItem(pageIndicatorTextView, len(pageIndicatorText)+6, 0, false)
	layout.AddItem(helpTextView, len(helpText)+4, 0, false)

	applicationHeader.statusTextView = statusTextView
	applicationHeader.layout = layout

	applicationHeader.updateUi()
}

func (applicationHeader *ApplicationHeaderComponent) updateUi() {
	page := applicationHeader.page
	pageName := string(applicationHeader.page)
	pageIdx := indexOf(page, Pages)
	pageCount := len(Pages)
	pageIndicatorText := fmt.Sprintf("Page: %s %d/%d", pageName, pageIdx+1, pageCount)

	applicationHeader.pageIndicatorTextView.SetText(pageIndicatorText)
	applicationHeader.layout.GetItem(3).SetRect(0, 0, len(pageIndicatorText)+8, 1)

	//applicationHeader.pageIndicatorTextView.SetRect(0, 0, len(pageIndicatorText)+4, 1)
	//applicationHeader.pageIndicatorTextView.SetLabelWidth(len(pageIndicatorText) + 4)
	//applicationHeader.pageIndicatorTextView.SetLabel(pageIndicatorText)

}

func (applicationHeader *ApplicationHeaderComponent) SetStatus(status *status_message.StatusMessage) {
	applicationHeader.statusTextView.SetText(status.Message).SetTextColor(status.Color)
	applicationHeader.application.ForceDraw()
	if status.Duration > 0 {
		go func() {
			time.Sleep(status.Duration)
			if applicationHeader.lastStatus != status {
				return
			}
			applicationHeader.ResetStatus()
		}()
	}
	applicationHeader.lastStatus = status
}

func (applicationHeader *ApplicationHeaderComponent) ResetStatus() {
	applicationHeader.statusTextView.SetText("").SetTextColor(tcell.ColorWhite)
	applicationHeader.application.ForceDraw()
}

func (applicationHeader *ApplicationHeaderComponent) SetPage(page Page) {
	applicationHeader.page = page
	applicationHeader.updateUi()
}

func indexOf[T comparable](word T, data []T) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}
