package util

import (
	"fan2go-tui/internal/uistate"
	"path/filepath"
	"testing"
)

type dummyEntry struct {
	Name string
}

func newDummyMemory() *SelectionMemory[dummyEntry] {
	return NewSelectionMemory(func(entry *dummyEntry) string { return entry.Name })
}

func TestSelectionMemory_GetReturnsNilForUnknownKey(t *testing.T) {
	if selection := newDummyMemory().Get("unknown"); selection != nil {
		t.Errorf("expected nil, got %+v", selection)
	}
}

func TestSelectionMemory_RememberOverwritesExistingValue(t *testing.T) {
	memory := newDummyMemory()

	memory.Remember("server", 1, &dummyEntry{Name: "alpha"})
	memory.Remember("server", 5, &dummyEntry{Name: "beta"})

	selection := memory.Get("server")
	if selection == nil || *selection != (SelectionInfo{Index: 5, Key: "beta"}) {
		t.Errorf("unexpected selection %+v", selection)
	}
}

func TestSelectionMemory_Restore(t *testing.T) {
	alpha, beta, gamma := &dummyEntry{Name: "alpha"}, &dummyEntry{Name: "beta"}, &dummyEntry{Name: "gamma"}
	entries := []*dummyEntry{alpha, beta, gamma}
	memory := newDummyMemory()

	if entry, header := memory.Restore("server", nil); entry != nil || header {
		t.Errorf("expected nothing for no entries, got %+v (header: %v)", entry, header)
	}

	// nothing remembered: the first entry
	if entry, _ := memory.Restore("server", entries); entry != alpha {
		t.Errorf("expected alpha, got %+v", entry)
	}

	// by identity, after reloading the entries
	memory.Remember("server", 0, &dummyEntry{Name: "gamma"})
	if entry, _ := memory.Restore("server", entries); entry != gamma {
		t.Errorf("expected gamma, got %+v", entry)
	}

	// gone: the one at its index
	memory.Remember("server", 1, &dummyEntry{Name: "removed"})
	if entry, _ := memory.Restore("server", entries); entry != beta {
		t.Errorf("expected beta, got %+v", entry)
	}
	memory.Remember("server", 7, &dummyEntry{Name: "removed"})
	if entry, _ := memory.Restore("server", entries); entry != gamma {
		t.Errorf("expected gamma, got %+v", entry)
	}
}

func TestSelectionMemory_Persist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	store := uistate.Load(path)
	memory := newDummyMemory().Persist(store, "selection", 10)

	memory.Remember("server", 2, &dummyEntry{Name: "beta"})
	if err := store.Flush(); err != nil {
		t.Fatal(err)
	}

	restored := newDummyMemory().Persist(uistate.Load(path), "selection", 10)
	if selection := restored.Get("server"); selection == nil || *selection != (SelectionInfo{Index: 2, Key: "beta"}) {
		t.Errorf("unexpected restored selection %+v", selection)
	}
	if selection := restored.Get("other"); selection != nil {
		t.Errorf("expected nil, got %+v", selection)
	}

	// without a store, nothing is persisted
	if selection := newDummyMemory().Persist(nil, "selection", 10).Get("server"); selection != nil {
		t.Errorf("expected nil, got %+v", selection)
	}
}
