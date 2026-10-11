package shortcut_helper

import (
	"fan2go-tui/internal/uistate"
	"sync/atomic"
	"weak"
)

// toggleHideShortcuts is the key of the setting in uistate.Current.
const toggleHideShortcuts = "ui.hideShortcuts"

var (
	// shortcutsHidden is whether collapsible shortcut maps are hidden.
	// Changed on the UI thread, atomic so it can be read anywhere.
	shortcutsHidden atomic.Bool

	// collapsibleMaps are the shortcut maps that are hidden with the setting. Weak, so closed components
	// are not kept alive. Only accessed on the UI thread.
	collapsibleMaps []weak.Pointer[ShortcutMapComponent]
)

// InitShortcutVisibility loads the setting from uistate.Current. Called when the UI is created.
func InitShortcutVisibility() {
	shortcutsHidden.Store(uistate.Current.Toggle(toggleHideShortcuts, false))
	collapsibleMaps = nil
}

// ShortcutsHidden returns whether collapsible shortcut maps are hidden.
func ShortcutsHidden() bool {
	return shortcutsHidden.Load()
}

// ToggleShortcuts hides or shows all collapsible shortcut maps, and remembers the setting.
// Must be called on the UI thread.
func ToggleShortcuts() {
	hidden := !shortcutsHidden.Load()
	shortcutsHidden.Store(hidden)
	uistate.Current.SetToggle(toggleHideShortcuts, hidden)

	alive := collapsibleMaps[:0]
	for _, pointer := range collapsibleMaps {
		if shortcutMap := pointer.Value(); shortcutMap != nil {
			shortcutMap.applyVisibility()
			alive = append(alive, pointer)
		}
	}
	collapsibleMaps = alive
}

// SetCollapsible makes the map follow the setting of ToggleShortcuts: hidden, it has no text and a height of 0
// (reported to the callback of SetOnHeightChanged, so its container can resize it).
// Must be called on the UI thread.
func (sm *ShortcutMapComponent) SetCollapsible() *ShortcutMapComponent {
	if !sm.collapsible {
		sm.collapsible = true
		collapsibleMaps = append(collapsibleMaps, weak.Make(sm))
	}
	return sm
}

// isHidden returns whether the map is hidden by the setting.
func (sm *ShortcutMapComponent) isHidden() bool {
	return sm.collapsible && shortcutsHidden.Load()
}

// applyVisibility shows the entries again or hides them, after the setting changed.
func (sm *ShortcutMapComponent) applyVisibility() {
	sm.SetEntries(sm.ShortCutEntries)
}
