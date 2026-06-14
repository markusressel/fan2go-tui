package shortcut_helper

import (
	"fan2go-tui/internal/ui/theme"
	"fan2go-tui/internal/ui/txwidgets"
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type ShortcutEntry struct {
	KeyCombo []string
	Name     string
}

type ShortcutMapComponent struct {
	application *tview.Application

	layout                  *tview.Flex
	shortcutEntriesTextView *tview.TextView

	ShortCutEntries []ShortcutEntry
}

func NewShortcutMap(application *tview.Application) *ShortcutMapComponent {
	shortcutMap := &ShortcutMapComponent{
		application: application,
	}

	shortcutMap.createLayout()

	return shortcutMap
}

func (sm *ShortcutMapComponent) createLayout() {
	layout := tview.NewFlex().SetDirection(tview.FlexColumn)

	shortcutEntriesTextView := tview.NewTextView().
		SetDynamicColors(true)
	shortcutEntriesTextView.SetBorderPadding(0, 0, 1, 1)
	shortcutEntriesTextView.SetTextAlign(tview.AlignLeft)

	layout.AddItem(shortcutEntriesTextView, 0, 1, false)

	sm.shortcutEntriesTextView = shortcutEntriesTextView
	sm.layout = layout
}

func (sm *ShortcutMapComponent) SetEntries(entries []ShortcutEntry) {
	sm.ShortCutEntries = entries
	var statusText string
	for _, entry := range entries {
		// comma separated list:
		shortCutsText := strings.Join(entry.KeyCombo, "|")
		shortcuts := txwidgets.Span(theme.Colors.ShortcutMap.KeyCombo, "[%s]", shortCutsText)
		name := txwidgets.Span(theme.Colors.ShortcutMap.Name, "%s", entry.Name)
		statusText += fmt.Sprintf("%s: %s  ", shortcuts, name)
	}
	sm.shortcutEntriesTextView.SetText(statusText)
	sm.application.ForceDraw()
}

func (sm *ShortcutMapComponent) Clear() {
	sm.shortcutEntriesTextView.SetText("")
	sm.application.ForceDraw()
}

func (sm *ShortcutMapComponent) GetLayout() *tview.Flex {
	return sm.layout
}
