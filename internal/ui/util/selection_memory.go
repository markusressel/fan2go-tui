package util

import (
	"fan2go-tui/internal/uistate"
	"slices"
)

// SelectionInfo is a remembered selection: the entry (by its identity, see NewSelectionMemory) and its index, to
// select an entry close to it if it is gone.
type SelectionInfo struct {
	Index int `json:"index"`
	// Key identifies the selected entry, empty if the header row was selected
	Key string `json:"key,omitempty"`
	// Offset is the row of the entry within the visible entries, to restore the scroll position
	Offset int `json:"offset,omitempty"`
}

// SelectionMemory remembers the selection of a list by a key, e.g. the selected fan of each fan2go daemon, so the
// selection is restored when the list is shown again (see Restore). With Persist, the selections are remembered
// between runs as well.
type SelectionMemory[T any] struct {
	identify   func(entry *T) string
	selections map[string]SelectionInfo

	// store and namespace are where the selections are persisted, see Persist
	store     *uistate.Store
	namespace string
}

// NewSelectionMemory creates a memory of selections, identifying the entries with identify (e.g. by their name),
// so an entry is found again after the list was reloaded.
func NewSelectionMemory[T any](identify func(entry *T) string) *SelectionMemory[T] {
	return &SelectionMemory[T]{
		identify:   identify,
		selections: map[string]SelectionInfo{},
	}
}

// Persist saves the selections in the namespace of store, keeping the maxSelections most recent ones, so they are
// remembered between runs. Does nothing if store is nil.
func (s *SelectionMemory[T]) Persist(store *uistate.Store, namespace string, maxSelections int) *SelectionMemory[T] {
	if store == nil {
		return s
	}
	s.store = store
	s.namespace = namespace
	store.SetLimit(namespace, maxSelections)
	return s
}

// Remember remembers entry at index as the selection of the list with the given key; nil means the header row is
// selected.
func (s *SelectionMemory[T]) Remember(key string, index int, entry *T) {
	info := SelectionInfo{Index: index}
	if entry != nil {
		info.Key = s.identify(entry)
	}
	s.RememberInfo(key, info)
}

// RememberInfo remembers the selection of the list with the given key, e.g. an entry that is not loaded yet.
func (s *SelectionMemory[T]) RememberInfo(key string, info SelectionInfo) {
	if current, ok := s.selections[key]; ok && current == info {
		return
	}
	s.selections[key] = info
	if s.store != nil {
		uistate.Set(s.store, s.namespace, key, info)
	}
}

// Get returns the remembered selection of the list with the given key, nil if there is none.
func (s *SelectionMemory[T]) Get(key string) *SelectionInfo {
	info, ok := s.selections[key]
	if !ok && s.store != nil {
		info, ok = uistate.Get[SelectionInfo](s.store, s.namespace, key)
		if ok {
			s.selections[key] = info
		}
	}
	if !ok {
		return nil
	}
	return &info
}

// Restore returns the entry of entries to select for the list with the given key: the remembered one, or the one
// at its index if it is gone, or the first entry if nothing is remembered. header is true if the header row was
// selected. Returns nil if entries is empty.
func (s *SelectionMemory[T]) Restore(key string, entries []*T) (entry *T, header bool) {
	if len(entries) == 0 {
		return nil, false
	}
	info := s.Get(key)
	if info == nil {
		return entries[0], false
	}
	if info.Key == "" {
		return nil, true
	}
	index := slices.IndexFunc(entries, func(entry *T) bool { return s.identify(entry) == info.Key })
	if index < 0 {
		index = min(max(info.Index, 0), len(entries)-1)
	}
	return entries[index], false
}
