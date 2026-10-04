package shortcut_helper

import (
	"strings"
	"testing"

	"github.com/rivo/tview"
)

func TestToggleShortcuts(t *testing.T) {
	collapsibleMaps = nil
	t.Cleanup(func() {
		shortcutsHidden.Store(false)
		collapsibleMaps = nil
	})
	InitShortcutVisibility()
	if ShortcutsHidden() {
		t.Fatal("expected shortcuts to be shown by default")
	}

	app := tview.NewApplication()
	entries := []ShortcutEntry{{KeyCombo: []string{Ctrl("q")}, Name: "Quit"}}
	newMap := func(collapsible bool) (*ShortcutMapComponent, *int) {
		shortcutMap := NewShortcutMap(app)
		if collapsible {
			shortcutMap.SetCollapsible()
		}
		height := -1
		shortcutMap.SetOnHeightChanged(func(h int) { height = h })
		shortcutMap.SetEntries(entries)
		return shortcutMap, &height
	}
	pageMap, pageHeight := newMap(true)
	dialogMap, dialogHeight := newMap(false)
	if *pageHeight <= 0 {
		t.Errorf("expected pageHeight > 0, got %d", *pageHeight)
	}
	if *dialogHeight <= 0 {
		t.Errorf("expected dialogHeight > 0, got %d", *dialogHeight)
	}

	ToggleShortcuts()
	if !ShortcutsHidden() {
		t.Error("expected ShortcutsHidden to be true")
	}
	if *pageHeight != 0 {
		t.Errorf("expected collapsible maps to take 0 space, got %d", *pageHeight)
	}
	if pageMap.shortcutEntriesTextView.GetText(false) != "" {
		t.Errorf("expected empty text when hidden, got %q", pageMap.shortcutEntriesTextView.GetText(false))
	}
	if h := pageMap.CalculateHeightForWidth(80); h != 0 {
		t.Errorf("expected height 0 when hidden, got %d", h)
	}
	if h := dialogMap.CalculateHeightForWidth(80); h <= 0 {
		t.Errorf("expected non-collapsible map to keep height > 0, got %d", h)
	}
	if dialogMap.shortcutEntriesTextView.GetText(false) == "" {
		t.Error("expected non-collapsible map to keep text")
	}

	// entries set while hidden are shown later
	pageMap.SetEntries([]ShortcutEntry{{KeyCombo: []string{"x"}, Name: "Other"}})
	if *pageHeight != 0 {
		t.Errorf("expected pageHeight to remain 0, got %d", *pageHeight)
	}

	ToggleShortcuts()
	if ShortcutsHidden() {
		t.Error("expected ShortcutsHidden to be false")
	}
	if *pageHeight <= 0 {
		t.Errorf("expected pageHeight > 0 after toggle, got %d", *pageHeight)
	}
	if text := pageMap.shortcutEntriesTextView.GetText(true); text == "" || !strings.Contains(text, "Other") {
		t.Errorf("expected text to contain 'Other', got %q", text)
	}
}

func TestSetCollapsibleRegistersOnce(t *testing.T) {
	collapsibleMaps = nil
	t.Cleanup(func() { collapsibleMaps = nil })

	shortcutMap := NewShortcutMap(tview.NewApplication())
	shortcutMap.SetCollapsible().SetCollapsible()
	if len(collapsibleMaps) != 1 {
		t.Errorf("expected 1 registered map, got %d", len(collapsibleMaps))
	}
}
