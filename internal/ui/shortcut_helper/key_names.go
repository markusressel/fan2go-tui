package shortcut_helper

// Key names for ShortcutEntry.KeyCombo, so all shortcuts look the same:
//   - characters as typed, so case matters: "h" is the plain key, "H" is shift+h
//   - named keys capitalized: KeyEnter, KeyEsc, ... and F-keys like "F1"
//   - modifiers lowercase with "+": Ctrl("q") is "ctrl+q", Shift(KeyTab) is "shift+⭾"
const (
	KeyEnter  = "Enter"
	KeyEsc    = "Esc"
	KeySpace  = "Space"
	KeyDelete = "Delete"
	KeyPgUp   = "PgUp"
	KeyPgDn   = "PgDn"
	KeyTab    = "⭾"
)

// Alt returns the name of the key with alt, e.g. "alt+f".
func Alt(key string) string {
	return "alt+" + key
}

// Ctrl returns the name of the key with ctrl, e.g. "ctrl+q".
func Ctrl(key string) string {
	return "ctrl+" + key
}

// Shift returns the name of the key with shift, e.g. "shift+⭾". For characters, use the shifted character instead,
// e.g. "H" instead of shift+h.
func Shift(key string) string {
	return "shift+" + key
}

var (
	// ShortcutHide hides or shows the shortcuts (see ToggleShortcuts).
	ShortcutHide = ShortcutEntry{KeyCombo: []string{"?", "F1"}, Name: "Hide shortcuts", Group: GroupGlobal}
)
