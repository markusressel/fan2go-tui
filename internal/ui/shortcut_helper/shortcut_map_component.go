package shortcut_helper

import (
	"cmp"
	"fan2go-tui/internal/ui/theme"
	"fan2go-tui/internal/ui/txwidgets"
	"os"
	"slices"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

type ShortcutEntry struct {
	KeyCombo []string
	Name     string
	// Group decides the color of the keys and where the entry is shown, see ShortcutGroup
	Group ShortcutGroup
}

// ShortcutGroup groups the entries of a shortcut map, so the one looked for is found quickly: the groups are
// shown in this order, separated by a line, and the keys of each group have their own color
// (theme.Colors.ShortcutMap). Within a group, the entries keep their order.
type ShortcutGroup int

const (
	// GroupAction is for actions on the selection or the component (the default)
	GroupAction ShortcutGroup = iota
	// GroupView is for keys that change what is shown, e.g. interval or modes
	GroupView
	// GroupNavigation is for keys that move the selection or the focus
	GroupNavigation
	// GroupGlobal is for keys that work everywhere, e.g. switching pages or quitting
	GroupGlobal
)

// keyColor returns the color of the keys of the group.
func (group ShortcutGroup) keyColor() tcell.Color {
	switch group {
	case GroupView:
		return theme.Colors.ShortcutMap.ViewKeyCombo
	case GroupNavigation:
		return theme.Colors.ShortcutMap.NavigationKeyCombo
	case GroupGlobal:
		return theme.Colors.ShortcutMap.GlobalKeyCombo
	default:
		return theme.Colors.ShortcutMap.KeyCombo
	}
}

// groupSeparator is shown between groups. Surrounded by spaces, so lines may wrap around it.
const groupSeparator = "│"

// sortedByGroup returns the entries ordered by group, keeping their order within a group.
func sortedByGroup(entries []ShortcutEntry) []ShortcutEntry {
	sorted := slices.Clone(entries)
	slices.SortStableFunc(sorted, func(a, b ShortcutEntry) int { return cmp.Compare(a.Group, b.Group) })
	return sorted
}

// formatEntries returns the text of the entries, styled or (for calculating its size) plain.
func formatEntries(entries []ShortcutEntry, styled bool) string {
	var text strings.Builder
	sorted := sortedByGroup(entries)
	for i, entry := range sorted {
		if i > 0 {
			text.WriteString("  ")
			if sorted[i-1].Group != entry.Group {
				if styled {
					text.WriteString(txwidgets.Span(theme.Colors.ShortcutMap.Separator, "%s", groupSeparator))
				} else {
					text.WriteString(groupSeparator)
				}
				text.WriteString("  ")
			}
		}
		// alternative keys joined with a non-breaking vertical line, so an entry is never wrapped
		keys := "[" + strings.Join(entry.KeyCombo, "\u01c0") + "]"
		name := strings.ReplaceAll(entry.Name, " ", "\u00a0")
		if styled {
			keys = txwidgets.Span(entry.Group.keyColor(), "%s", keys)
			name = txwidgets.Span(theme.Colors.ShortcutMap.Name, "%s", name)
		}
		text.WriteString(keys + ":\u00a0" + name)
	}
	return text.String()
}

type ShortcutMapComponent struct {
	application *tview.Application

	layout                  *tview.Flex
	shortcutEntriesTextView *tview.TextView
	onHeightChanged         func(height int)
	// collapsible maps are hidden with ToggleShortcuts, see SetCollapsible
	collapsible bool

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

	// Set draw func to monitor and dynamically resize height on line wraps
	shortcutEntriesTextView.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		lines := sm.CalculateHeightForWidth(width)
		if height != lines {
			if sm.onHeightChanged != nil {
				go func() {
					sm.application.QueueUpdateDraw(func() {
						if sm.onHeightChanged != nil {
							sm.onHeightChanged(lines)
						}
					})
				}()
			}
		}
		return x, y, width, height
	})

	layout.AddItem(shortcutEntriesTextView, 0, 1, false)

	sm.shortcutEntriesTextView = shortcutEntriesTextView
	sm.layout = layout
}

func (sm *ShortcutMapComponent) SetOnHeightChanged(f func(height int)) {
	sm.onHeightChanged = f
}

func (sm *ShortcutMapComponent) SetEntries(entries []ShortcutEntry) {
	sm.ShortCutEntries = entries
	if sm.isHidden() {
		sm.shortcutEntriesTextView.SetText("")
		if sm.onHeightChanged != nil {
			sm.onHeightChanged(0)
		}
		return
	}
	sm.shortcutEntriesTextView.SetText(formatEntries(entries, true))
	if sm.onHeightChanged != nil {
		lines := sm.CalculateHeightFromTerminal()
		sm.onHeightChanged(lines)
	}
}

func (sm *ShortcutMapComponent) Clear() {
	sm.shortcutEntriesTextView.SetText("")
	if sm.onHeightChanged != nil {
		height := 1
		if sm.isHidden() {
			height = 0
		}
		sm.onHeightChanged(height)
	}
}

func (sm *ShortcutMapComponent) GetLayout() *tview.Flex {
	return sm.layout
}

func (sm *ShortcutMapComponent) CalculateHeightFromTerminal() int {
	_, _, width, _ := sm.layout.GetRect()
	if width <= 0 {
		var err error
		width, _, err = term.GetSize(int(os.Stdout.Fd()))
		if err != nil || width <= 0 {
			width = 80
		}
	}
	return sm.CalculateHeightForWidth(width)
}

func (sm *ShortcutMapComponent) CalculateHeightForWidth(width int) int {
	if sm.isHidden() {
		return 0
	}
	availableWidth := width - 2 // padding
	if availableWidth <= 0 {
		availableWidth = 80
	}

	visibleText := formatEntries(sm.ShortCutEntries, false)
	if len(visibleText) == 0 {
		return 1
	}

	runes := []rune(visibleText)
	lines := 1
	currentLineLength := 0

	i := 0
	for i < len(runes) {
		if runes[i] == ' ' {
			currentLineLength++
			i++
		} else {
			wordStart := i
			for i < len(runes) && runes[i] != ' ' {
				i++
			}
			wordWidth := i - wordStart

			if currentLineLength+wordWidth > availableWidth {
				lines++
				currentLineLength = wordWidth
			} else {
				currentLineLength += wordWidth
			}
		}
	}

	return lines
}
