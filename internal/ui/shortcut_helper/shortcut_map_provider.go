package shortcut_helper

type ShortcutMapProvider interface {
	GetShortcutMap() []ShortcutEntry
}
