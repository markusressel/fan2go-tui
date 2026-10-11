package uistate

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/mitchellh/go-homedir"
)

func newTestStore(t *testing.T) (*Store, string) {
	path := filepath.Join(t.TempDir(), "sub", "state.json")
	store := Load(path)
	// long enough to never fire during a test, unless the test waits for it
	store.saveDelay = time.Hour
	return store, path
}

func flush(t *testing.T, store *Store) {
	t.Helper()
	if err := store.Flush(); err != nil {
		t.Fatalf("flush failed: %v", err)
	}
}

func readFileContent(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read %s: %v", path, err)
	}
	return string(content)
}

func writeFileContent(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("cannot write %s: %v", path, err)
	}
}

func assertValue(t *testing.T, store *Store, ns string, key string, expected int) {
	t.Helper()
	value, ok := Get[int](store, ns, key)
	if !ok || value != expected {
		t.Errorf("expected %s/%s to be %d, got %d (saved: %v)", ns, key, expected, value, ok)
	}
}

func assertMissing(t *testing.T, store *Store, ns string, key string) {
	t.Helper()
	if value, ok := Get[int](store, ns, key); ok {
		t.Errorf("expected %s/%s to be missing, got %d", ns, key, value)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	store, path := newTestStore(t)

	assertMissing(t, store, "session", "server")
	// nothing is written without changes
	flush(t, store)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected no state file, got %v", err)
	}
}

type testSelection struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

func TestStore_SaveAndLoad(t *testing.T) {
	store, path := newTestStore(t)

	if _, ok := Get[testSelection](store, "selection", "fans"); ok {
		t.Fatalf("expected no value before it is saved")
	}
	Set(store, "selection", "fans", testSelection{Name: "cpu", Index: 2})
	Set(store, "selection", "other", 1)
	store.Delete("selection", "other")
	flush(t, store)

	loaded := Load(path)
	value, ok := Get[testSelection](loaded, "selection", "fans")
	if !ok || value != (testSelection{Name: "cpu", Index: 2}) {
		t.Errorf("unexpected loaded value %+v (saved: %v)", value, ok)
	}
	assertMissing(t, loaded, "selection", "other")
	// a value of another type is ignored
	if _, ok := Get[int](loaded, "selection", "fans"); ok {
		t.Errorf("expected a value of another type to be ignored")
	}

	content := readFileContent(t, path)
	if !strings.Contains(content, `"version": 1`) || !strings.Contains(content, `"compatibleVersion": 1`) {
		t.Errorf("expected the versions in the file, got %s", content)
	}

	// only the state file and its lock are left, no temporary files
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	if !slices.Equal(names, []string{"state.json", "state.json.lock"}) {
		t.Errorf("unexpected files %v", names)
	}
}

func TestStore_SavesInBackground(t *testing.T) {
	store, path := newTestStore(t)
	store.saveDelay = 10 * time.Millisecond

	Set(store, "session", "server", 1)

	deadline := time.Now().Add(2 * time.Second)
	for {
		if _, ok := Get[int](Load(path), "session", "server"); ok {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("the value was not saved in the background")
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestStore_ReturnsCopies(t *testing.T) {
	store, _ := newTestStore(t)
	ids := []string{"cpu", "gpu"}
	Set(store, "selection", "fans", ids)
	ids[0] = "changed"

	loaded, _ := Get[[]string](store, "selection", "fans")
	if !slices.Equal(loaded, []string{"cpu", "gpu"}) {
		t.Errorf("expected the saved value to be a copy, got %v", loaded)
	}
}

func TestLoad_CorruptFileStartsEmptyAndIsReplaced(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	writeFileContent(t, path, "{not json")

	store := Load(path)
	assertMissing(t, store, "session", "server")

	Set(store, "session", "server", 1)
	flush(t, store)
	assertValue(t, Load(path), "session", "server", 1)
}

// A file of a newer version that this one cannot write is read as far as possible, but never written, so its
// values are not lost.
func TestLoad_NewerVersionIsNotWritten(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	content := `{"version": 99, "compatibleVersion": 99, "namespaces": {"session": {"server": {"value": 3}}}, "unknown": true}`
	writeFileContent(t, path, content)

	store := Load(path)
	assertValue(t, store, "session", "server", 3)
	version, newer := store.NewerVersion()
	if !newer || version != 99 {
		t.Errorf("expected newer version 99, got %d (newer: %v)", version, newer)
	}

	// changes work in memory, but are not saved
	store.SetToggle("hideShortcuts", true)
	if !store.Toggle("hideShortcuts", false) {
		t.Errorf("expected the change in memory")
	}
	flush(t, store)
	if written := readFileContent(t, path); written != content {
		t.Errorf("expected the file to be unchanged, got %s", written)
	}
}

func TestLoad_UndecodableNewerVersionIsNotWritten(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	content := `{"version": 5, "compatibleVersion": 5, "namespaces": ["changed structure"]}`
	writeFileContent(t, path, content)

	store := Load(path)
	version, newer := store.NewerVersion()
	if !newer || version != 5 {
		t.Errorf("expected newer version 5, got %d (newer: %v)", version, newer)
	}

	store.SetToggle("hideShortcuts", true)
	flush(t, store)
	if written := readFileContent(t, path); written != content {
		t.Errorf("expected the file to be unchanged, got %s", written)
	}
}

// A newer version that is compatible with this one (e.g. one with additional namespaces) is written, keeping its
// values and its version.
func TestLoad_CompatibleNewerVersionIsWritten(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	writeFileContent(t, path, `{"version": 3, "compatibleVersion": 1, "namespaces": {"future": {"key": {"value": 1}}}}`)

	store := Load(path)
	if _, newer := store.NewerVersion(); newer {
		t.Errorf("expected a compatible version not to be reported as newer")
	}
	store.SetToggle("hideShortcuts", true)
	flush(t, store)

	loaded := Load(path)
	if !loaded.Toggle("hideShortcuts", false) {
		t.Errorf("expected the change to be saved")
	}
	assertValue(t, loaded, "future", "key", 1)
	written := readFileContent(t, path)
	if !strings.Contains(written, `"version": 3`) || !strings.Contains(written, `"compatibleVersion": 1`) {
		t.Errorf("expected the versions to be kept, got %s", written)
	}
}

// A newer version started meanwhile wrote the file: it is not overwritten.
func TestStore_DoesNotOverwriteNewerVersionWrittenMeanwhile(t *testing.T) {
	store, path := newTestStore(t)
	store.SetToggle("hideShortcuts", true)
	flush(t, store)
	if _, newer := store.NewerVersion(); newer {
		t.Fatalf("expected no newer version")
	}

	content := `{"version": 4, "compatibleVersion": 4, "namespaces": {}}`
	writeFileContent(t, path, content)
	store.SetToggle("hideShortcuts", false)
	flush(t, store)

	if written := readFileContent(t, path); written != content {
		t.Errorf("expected the file to be unchanged, got %s", written)
	}
	if version, newer := store.NewerVersion(); !newer || version != 4 {
		t.Errorf("expected newer version 4, got %d (newer: %v)", version, newer)
	}
}

func TestStore_Nil(t *testing.T) {
	var store *Store
	Set(store, "session", "server", 1)
	store.Delete("session", "server")
	store.SetLimit("session", 1)
	store.SetToggle("hideShortcuts", false)
	assertMissing(t, store, "session", "server")
	if !store.Toggle("hideShortcuts", true) {
		t.Errorf("expected the default value")
	}
	if _, newer := store.NewerVersion(); newer {
		t.Errorf("expected no newer version")
	}
	if err := store.Flush(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestStore_FlushErrorKeepsChanges(t *testing.T) {
	dir := t.TempDir()
	// the parent "directory" of the state file is a file, so writing fails
	blocker := filepath.Join(dir, "blocker")
	writeFileContent(t, blocker, "")
	store := Load(filepath.Join(blocker, "state.json"))
	store.saveDelay = time.Hour

	Set(store, "session", "server", 1)
	if err := store.Flush(); err == nil {
		t.Fatalf("expected an error")
	}

	// the change is written once writing works again
	if err := os.Remove(blocker); err != nil {
		t.Fatal(err)
	}
	flush(t, store)
	assertValue(t, Load(filepath.Join(blocker, "state.json")), "session", "server", 1)
}

func TestDefaultPath(t *testing.T) {
	t.Setenv("XDG_STATE_HOME", "/custom/state")
	path, err := DefaultPath()
	if err != nil || path != "/custom/state/fan2go-tui/state.json" {
		t.Errorf("unexpected path %s (error: %v)", path, err)
	}

	// relative paths are invalid according to the XDG spec, so the default is used
	t.Setenv("XDG_STATE_HOME", "relative")
	t.Setenv("HOME", "/home/test")
	// go-homedir caches the home directory
	homedir.Reset()
	t.Cleanup(homedir.Reset)
	path, err = DefaultPath()
	if err != nil || path != "/home/test/.local/state/fan2go-tui/state.json" {
		t.Errorf("unexpected path %s (error: %v)", path, err)
	}
}

func TestStore_Toggles(t *testing.T) {
	store, path := newTestStore(t)

	if !store.Toggle("hideShortcuts", true) || store.Toggle("hideShortcuts", false) {
		t.Errorf("expected the default value without a saved one")
	}

	store.SetToggle("hideShortcuts", false)
	store.SetToggle("other", true)
	flush(t, store)

	loaded := Load(path)
	if loaded.Toggle("hideShortcuts", true) || !loaded.Toggle("other", false) {
		t.Errorf("expected the saved values")
	}
}

func TestStore_SetWithoutChangeDoesNotSave(t *testing.T) {
	store, path := newTestStore(t)
	store.SetToggle("hideShortcuts", true)
	flush(t, store)
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	store.SetToggle("hideShortcuts", true)
	flush(t, store)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected no state file, got %v", err)
	}
}

// newClock returns a clock for Store.now that advances by a second on every call.
func newClock() func() time.Time {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	return func() time.Time {
		now = now.Add(time.Second)
		return now
	}
}

func TestStore_SetLimitDropsLeastRecentlySaved(t *testing.T) {
	store, path := newTestStore(t)
	store.now = newClock()
	store.SetLimit("selection", 2)

	Set(store, "selection", "a", 1)
	Set(store, "selection", "b", 2)
	Set(store, "selection", "c", 3)
	assertMissing(t, store, "selection", "a")

	flush(t, store)
	loaded := Load(path)
	assertMissing(t, loaded, "selection", "a")
	assertValue(t, loaded, "selection", "b", 2)
	assertValue(t, loaded, "selection", "c", 3)
}

func TestStore_InstancesKeepTheChangesOfEachOther(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	clock := newClock()
	first := Load(path)
	first.saveDelay, first.now = time.Hour, clock
	second := Load(path)
	second.saveDelay, second.now = time.Hour, clock

	Set(first, "session", "server", 1)
	first.SetToggle("tree", true)
	second.SetToggle("hidden", true)
	flush(t, first)
	// the second instance started before the first one saved, it must not drop its values
	flush(t, second)

	loaded := Load(path)
	assertValue(t, loaded, "session", "server", 1)
	if !loaded.Toggle("tree", false) || !loaded.Toggle("hidden", false) {
		t.Errorf("expected the values of both instances")
	}

	// the most recent change wins, even if it was saved first
	second.SetToggle("tree", false)
	first.SetToggle("tree", true)
	first.SetToggle("tree", false)
	first.SetToggle("tree", true)
	flush(t, first)
	flush(t, second)
	if !Load(path).Toggle("tree", false) {
		t.Errorf("expected the most recent change")
	}

	// a deletion is kept as well
	first.Delete("session", "server")
	flush(t, first)
	second.SetToggle("other", true)
	flush(t, second)
	assertMissing(t, Load(path), "session", "server")

	// a value saved by another instance can be deleted, once it is known
	third := Load(path)
	third.saveDelay, third.now = time.Hour, clock
	third.SetToggle("hidden", false)
	third.Delete(namespaceToggles, "other")
	flush(t, third)
	loaded = Load(path)
	if loaded.Toggle("hidden", true) || loaded.Toggle("other", false) {
		t.Errorf("expected the changes of the third instance")
	}
}
