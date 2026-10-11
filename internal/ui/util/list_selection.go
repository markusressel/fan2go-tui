package util

import (
	"fan2go-tui/internal/uistate"
	"slices"
)

// ListSelection remembers the selection of a ListComponent (see SelectionMemory), and selects the remembered entry
// once the list has entries, e.g. after the first data was received from the fan2go daemon.
type ListSelection[T comparable] struct {
	list     *ListComponent[T]
	identify func(entry *T) string
	memory   *SelectionMemory[T]
	// key is the key of the selection in memory, see Persist
	key string
	// restored is true once the remembered entry was selected; until then, the selection is not remembered, so the
	// first entry selected by default does not replace the remembered one
	restored bool
}

// NewListSelection remembers the selection of list, identifying the entries with identify (e.g. by their ID).
func NewListSelection[T comparable](list *ListComponent[T], identify func(entry *T) string) *ListSelection[T] {
	selection := &ListSelection[T]{
		list:     list,
		identify: identify,
		memory:   NewSelectionMemory(identify),
	}
	list.SetSelectionChangedCallback(selection.remember)
	return selection
}

// Persist remembers the selection under key (e.g. the address of the fan2go daemon) in the namespace of store,
// keeping the maxSelections most recent keys, so it is restored in the next run. If restore is false, the selection
// of the last run is not restored, but replaced. Must be called before the list has entries. Does nothing if store
// is nil.
func (s *ListSelection[T]) Persist(store *uistate.Store, namespace string, key string, maxSelections int, restore bool) {
	s.memory.Persist(store, namespace, maxSelections)
	s.key = key
	if !restore {
		s.restored = true
	}
}

// Update must be called after the entries of the list were set: the first time the list has entries, the remembered
// entry is selected and scrolled to its remembered position, afterwards the selected entry and the scroll position
// are remembered (e.g. an entry selected with the mouse).
func (s *ListSelection[T]) Update() {
	if s.restored {
		s.remember(s.list.GetSelectedItem())
		return
	}
	entry, _ := s.memory.Restore(s.key, s.list.GetData())
	if entry == nil {
		return
	}
	s.restored = true
	offset := 0
	if info := s.memory.Get(s.key); info != nil {
		offset = info.Offset
	}
	s.list.RestoreSelection(entry, offset)
}

func (s *ListSelection[T]) remember(entry *T) {
	if !s.restored || entry == nil {
		return
	}
	s.memory.RememberInfo(s.key, SelectionInfo{
		Index:  slices.Index(s.list.GetData(), entry),
		Key:    s.identify(entry),
		Offset: s.list.GetScrollOffset(entry),
	})
}
