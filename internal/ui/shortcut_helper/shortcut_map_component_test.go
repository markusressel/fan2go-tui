package shortcut_helper

import (
	"fan2go-tui/internal/ui/theme"
	"fan2go-tui/internal/ui/txwidgets"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestNewShortcutMap(t *testing.T) {
	app := tview.NewApplication()
	sm := NewShortcutMap(app)

	if sm == nil {
		t.Fatal("expected non-nil ShortcutMap")
	}
	if sm.GetLayout() == nil {
		t.Fatal("expected non-nil Layout")
	}
	if sm.shortcutEntriesTextView == nil {
		t.Fatal("expected non-nil shortcutEntriesTextView")
	}
}

func TestSetOnHeightChanged(t *testing.T) {
	app := tview.NewApplication()
	sm := NewShortcutMap(app)

	heightCalled := -1
	sm.SetOnHeightChanged(func(height int) {
		heightCalled = height
	})

	entries := []ShortcutEntry{
		{KeyCombo: []string{Ctrl("q")}, Name: "Quit"},
	}

	sm.SetEntries(entries)
	if heightCalled <= 0 {
		t.Errorf("expected heightCalled > 0, got %d", heightCalled)
	}
}

func TestSetEntriesAndClear(t *testing.T) {
	app := tview.NewApplication()
	sm := NewShortcutMap(app)

	entries := []ShortcutEntry{
		{KeyCombo: []string{KeyTab, Shift(KeyTab)}, Name: "Cycle focus"},
		{KeyCombo: []string{Ctrl("q")}, Name: "Quit"},
	}

	sm.SetEntries(entries)

	if len(sm.ShortCutEntries) != len(entries) {
		t.Fatalf("expected %d entries, got %d", len(entries), len(sm.ShortCutEntries))
	}

	text := sm.shortcutEntriesTextView.GetText(false)
	if !strings.Contains(text, "Cycle\u00a0focus") {
		t.Errorf("expected text to contain 'Cycle\u00a0focus', got %q", text)
	}
	if !strings.Contains(text, "⭾\u01c0shift+⭾") {
		t.Errorf("expected text to contain '⭾\u01c0shift+⭾', got %q", text)
	}

	sm.Clear()
	if sm.shortcutEntriesTextView.GetText(false) != "" {
		t.Errorf("expected empty text after clear, got %q", sm.shortcutEntriesTextView.GetText(false))
	}
}

func TestCalculateHeightForWidth(t *testing.T) {
	app := tview.NewApplication()
	sm := NewShortcutMap(app)

	sm.ShortCutEntries = nil
	if h := sm.CalculateHeightForWidth(80); h != 1 {
		t.Errorf("expected height 1 for empty entries, got %d", h)
	}

	sm.ShortCutEntries = []ShortcutEntry{
		{KeyCombo: []string{Ctrl("q")}, Name: "Quit"},
	}
	if h := sm.CalculateHeightForWidth(30); h != 1 {
		t.Errorf("expected height 1, got %d", h)
	}
	if h := sm.CalculateHeightForWidth(10); h != 2 {
		t.Errorf("expected height 2, got %d", h)
	}

	sm.ShortCutEntries = []ShortcutEntry{
		{KeyCombo: []string{KeyTab, Shift(KeyTab)}, Name: "Cycle focus"},
	}
	if h := sm.CalculateHeightForWidth(30); h != 1 {
		t.Errorf("expected height 1, got %d", h)
	}
	if h := sm.CalculateHeightForWidth(20); h != 2 {
		t.Errorf("expected height 2, got %d", h)
	}

	sm.ShortCutEntries = []ShortcutEntry{
		{KeyCombo: []string{KeyTab, Shift(KeyTab)}, Name: "Cycle focus"},
		{KeyCombo: []string{"F5"}, Name: "Refresh"},
		{KeyCombo: []string{Ctrl("q")}, Name: "Quit"},
		{KeyCombo: []string{"↑", "k"}, Name: "Move up"},
		{KeyCombo: []string{"↓", "j"}, Name: "Move down"},
	}

	if h := sm.CalculateHeightForWidth(200); h != 1 {
		t.Errorf("expected height 1 for 200 width, got %d", h)
	}
	if h := sm.CalculateHeightForWidth(40); h <= 1 {
		t.Errorf("expected height > 1 for 40 width, got %d", h)
	}
}

func TestCalculateHeightFromTerminal(t *testing.T) {
	app := tview.NewApplication()
	sm := NewShortcutMap(app)

	sm.ShortCutEntries = []ShortcutEntry{
		{KeyCombo: []string{Ctrl("q")}, Name: "Quit"},
	}

	height := sm.CalculateHeightFromTerminal()
	if height < 1 {
		t.Errorf("expected height >= 1, got %d", height)
	}
}

func TestDrawFuncHeightResize(t *testing.T) {
	app := tview.NewApplication()
	sm := NewShortcutMap(app)

	sm.SetOnHeightChanged(func(height int) {})

	sm.ShortCutEntries = []ShortcutEntry{
		{KeyCombo: []string{Ctrl("q")}, Name: "Quit"},
	}

	drawFunc := sm.shortcutEntriesTextView.GetDrawFunc()
	if drawFunc == nil {
		t.Fatal("expected non-nil DrawFunc")
	}

	drawFunc(nil, 0, 0, 8, 1)
}

func TestFormatEntriesGroupsShortcuts(t *testing.T) {
	entries := []ShortcutEntry{
		{KeyCombo: []string{"q"}, Name: "Quit", Group: GroupGlobal},
		{KeyCombo: []string{"↑", "↓"}, Name: "Move", Group: GroupNavigation},
		{KeyCombo: []string{"h"}, Name: "History"},
		{KeyCombo: []string{"F2"}, Name: "Columns", Group: GroupView},
		{KeyCombo: []string{"r"}, Name: "Restore file"},
	}

	wantPlain := "[h]:\u00a0History  [r]:\u00a0Restore\u00a0file  │  [F2]:\u00a0Columns  │  [↑ǀ↓]:\u00a0Move  │  [q]:\u00a0Quit"
	gotPlain := formatEntries(entries, false)
	if gotPlain != wantPlain {
		t.Errorf("expected %q, got %q", wantPlain, gotPlain)
	}

	singleGroupWant := "[h]:\u00a0History"
	singleGroupGot := formatEntries(entries[2:3], false)
	if singleGroupGot != singleGroupWant {
		t.Errorf("expected %q, got %q", singleGroupWant, singleGroupGot)
	}

	styled := formatEntries(entries, true)
	for _, color := range []tcell.Color{
		theme.Colors.ShortcutMap.KeyCombo, theme.Colors.ShortcutMap.ViewKeyCombo,
		theme.Colors.ShortcutMap.NavigationKeyCombo, theme.Colors.ShortcutMap.GlobalKeyCombo,
		theme.Colors.ShortcutMap.Separator,
	} {
		colorTag := txwidgets.ColorTag(color)
		if !strings.Contains(styled, colorTag) {
			t.Errorf("expected styled to contain color tag %q", colorTag)
		}
	}

	view := tview.NewTextView().SetDynamicColors(true).SetText(styled)
	if view.GetText(true) != gotPlain {
		t.Errorf("stripped dynamic colors %q != plain text %q", view.GetText(true), gotPlain)
	}
	if entries[0].Name != "Quit" {
		t.Errorf("original entries should not be mutated in-place, got %q", entries[0].Name)
	}
}

func TestGroupColorsDiffer(t *testing.T) {
	colors := map[tcell.Color]ShortcutGroup{}
	for _, group := range []ShortcutGroup{GroupAction, GroupView, GroupNavigation, GroupGlobal} {
		color := group.keyColor()
		if previous, exists := colors[color]; exists {
			t.Errorf("group %d has same color (%v) as group %d", group, color, previous)
		}
		colors[color] = group
	}
}
