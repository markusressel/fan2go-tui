// Package uistate persists what the application remembers between runs (e.g. the last page or selections), as
// opposed to the user-written configuration in internal/configuration and the data of the fan2go daemon in
// internal/state. The state is a small JSON file in $XDG_STATE_HOME, written by the application only.
//
// The state is a set of namespaces (e.g. "toggles"), each mapping keys to JSON values, see Get and Set. Several
// running instances share the file: each one keeps its own state in memory, and only writes the values it changed,
// merged into the file as it is at that time, so the changes of other instances are not lost.
package uistate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fan2go-tui/internal/logging"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"syscall"
	"time"

	"github.com/mitchellh/go-homedir"
)

// currentVersion is the version of the file format written by this version of the application.
const currentVersion = 1

// compatibleVersion is the oldest version of the file format whose readers can read and write files of
// currentVersion without losing values. Only increase it with changes older versions cannot handle (e.g. a changed
// structure, or new fields besides the namespaces); new namespaces or values need not, they are kept by older
// versions.
const compatibleVersion = 1

// defaultSaveDelay debounces saving, so fast changes (e.g. scrolling through a list) are written once.
const defaultSaveDelay = 500 * time.Millisecond

// namespaceToggles is the namespace of the on/off settings, see Toggle.
const namespaceToggles = "toggles"

// Current is the state of the running application, set up at startup by Init.
// While it is nil (e.g. in tests), nothing is restored or saved.
var Current *Store

// entry is a saved value, with the time it was saved, to merge the changes of several instances and to drop the
// least recently saved values of a namespace with a limit, see Store.SetLimit.
type entry struct {
	Value     json.RawMessage `json:"value"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type namespace map[string]entry

// fileContent is the content of the state file.
type fileContent struct {
	Version int `json:"version"`
	// CompatibleVersion is the oldest version that can read and write the file, see compatibleVersion. Files written
	// before it existed are only compatible with their own version.
	CompatibleVersion int                  `json:"compatibleVersion,omitempty"`
	Namespaces        map[string]namespace `json:"namespaces,omitempty"`
}

// change is a value changed (or deleted) by this instance and not yet written to the file.
type change struct {
	namespace string
	key       string
}

// Store holds the state in memory and saves it to its file in the background after every change.
// All methods are safe for concurrent use and never block on file access, except Flush.
// The methods of a nil *Store do nothing.
type Store struct {
	path      string
	saveDelay time.Duration
	// now returns the current time, replaced in tests
	now func() time.Time

	mutex      sync.Mutex
	namespaces map[string]namespace
	// limits are the maximum numbers of values of namespaces, see SetLimit
	limits map[string]int
	// changes are the values changed since the last save, with the time of the change
	changes   map[change]time.Time
	saveTimer *time.Timer

	// newerVersion is the version of the state file, if it was written by a newer version of the application that
	// this one cannot write (see compatibleVersion), 0 otherwise. Then nothing is saved, so its values are not lost.
	newerVersion int

	// saveMutex serializes writing the file
	saveMutex sync.Mutex
}

// formatVersion is the version of a state file, see fileContent.
type formatVersion struct {
	version    int
	compatible int
}

// isSupported returns whether this version of the application can read and write files of the version.
func (version formatVersion) isSupported() bool {
	return version.compatible <= currentVersion
}

// Init loads the state from the default path into Current. Errors are logged, the application then starts
// with an empty state.
func Init() {
	path, err := DefaultPath()
	if err != nil {
		logging.Error("Cannot determine the state file path, the state is not saved: %v", err)
		return
	}
	Current = Load(path)
}

// DefaultPath returns the path of the state file: $XDG_STATE_HOME/fan2go-tui/state.json,
// with $XDG_STATE_HOME defaulting to ~/.local/state.
func DefaultPath() (string, error) {
	stateHome := os.Getenv("XDG_STATE_HOME")
	if !filepath.IsAbs(stateHome) {
		// relative paths are invalid according to the XDG spec
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		stateHome = filepath.Join(home, ".local", "state")
	}
	return filepath.Join(stateHome, "fan2go-tui", "state.json"), nil
}

// Load reads the state file at path. A missing file results in an empty state. An unreadable file is logged and
// results in an empty state as well; it is overwritten with the next change.
func Load(path string) *Store {
	store := &Store{
		path:       path,
		saveDelay:  defaultSaveDelay,
		now:        time.Now,
		namespaces: map[string]namespace{},
		limits:     map[string]int{},
		changes:    map[change]time.Time{},
	}

	namespaces, version, err := readFile(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
	case err != nil && !version.isSupported():
		// the values cannot be read either, but must not be lost
		store.setNewerVersion(version.version)
	case err != nil:
		logging.Error("Cannot read state file, starting with an empty state: %v", err)
	default:
		store.namespaces = namespaces
		if !version.isSupported() {
			store.setNewerVersion(version.version)
		}
	}
	return store
}

// setNewerVersion stops saving, because the state file was written by a newer version of the application (see
// NewerVersion).
func (store *Store) setNewerVersion(version int) {
	if store.newerVersion == 0 {
		logging.Warning("State file %s was written by a newer version of the application (format %d, this one "+
			"supports %d), changes are not saved", store.path, version, currentVersion)
	}
	store.newerVersion = version
}

// NewerVersion returns the format version of the state file and true, if it was written by a newer version of the
// application which this one cannot write. Then changes are not saved, so the values of the newer version are not
// lost; the user should update the application.
func (store *Store) NewerVersion() (int, bool) {
	if store == nil {
		return 0, false
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	return store.newerVersion, store.newerVersion != 0
}

// readFile returns the namespaces of the state file at path and the version of the file. If the file cannot be decoded, the version is still returned if it can be read.
func readFile(path string) (map[string]namespace, formatVersion, error) {
	var version formatVersion
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, version, err
	}
	var content fileContent
	if err := json.Unmarshal(data, &content); err != nil {
		var header struct {
			Version           int `json:"version"`
			CompatibleVersion int `json:"compatibleVersion"`
		}
		if json.Unmarshal(data, &header) == nil {
			version = newFormatVersion(header.Version, header.CompatibleVersion)
		}
		return nil, version, fmt.Errorf("%s: %w", path, err)
	}
	version = newFormatVersion(content.Version, content.CompatibleVersion)

	namespaces := content.Namespaces
	if namespaces == nil {
		namespaces = map[string]namespace{}
	}
	return namespaces, version, nil
}

func newFormatVersion(version int, compatible int) formatVersion {
	if compatible == 0 {
		// written before compatibleVersion existed (or by hand)
		compatible = version
	}
	return formatVersion{version: version, compatible: compatible}
}

// Get returns the value saved under key in the namespace ns, and whether there is one (that can be decoded as T).
func Get[T any](store *Store, ns string, key string) (T, bool) {
	var value T
	if store == nil {
		return value, false
	}
	store.mutex.Lock()
	saved, ok := store.namespaces[ns][key]
	store.mutex.Unlock()
	if !ok {
		return value, false
	}
	if err := json.Unmarshal(saved.Value, &value); err != nil {
		logging.Warning("Ignoring the saved state %s/%s: %v", ns, key, err)
		var zero T
		return zero, false
	}
	return value, true
}

// Set saves value under key in the namespace ns. Saving the value that is already saved does nothing.
func Set[T any](store *Store, ns string, key string, value T) {
	if store == nil {
		return
	}
	data, err := json.Marshal(value)
	if err != nil {
		logging.Error("Cannot save the state %s/%s: %v", ns, key, err)
		return
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if saved, ok := store.namespaces[ns][key]; ok && bytes.Equal(saved.Value, data) {
		return
	}
	if store.namespaces[ns] == nil {
		store.namespaces[ns] = namespace{}
	}
	now := store.now()
	store.namespaces[ns][key] = entry{Value: data, UpdatedAt: now}
	applyLimit(store.namespaces[ns], store.limits[ns])
	store.changes[change{ns, key}] = now
	store.scheduleSave()
}

// Delete removes the value saved under key in the namespace ns.
func (store *Store) Delete(ns string, key string) {
	if store == nil {
		return
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	if _, ok := store.namespaces[ns][key]; !ok {
		return
	}
	delete(store.namespaces[ns], key)
	store.changes[change{ns, key}] = store.now()
	store.scheduleSave()
}

// SetLimit keeps at most maxValues values in the namespace ns (e.g. the selections for each fan2go daemon, which are
// not needed forever), dropping the least recently saved ones. 0 means no limit.
func (store *Store) SetLimit(ns string, maxValues int) {
	if store == nil {
		return
	}
	store.mutex.Lock()
	defer store.mutex.Unlock()
	store.limits[ns] = maxValues
	applyLimit(store.namespaces[ns], maxValues)
}

// applyLimit removes the least recently saved values of values, so at most maxValues are left.
func applyLimit(values namespace, maxValues int) {
	if maxValues <= 0 || len(values) <= maxValues {
		return
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	slices.SortFunc(keys, func(a, b string) int { return values[a].UpdatedAt.Compare(values[b].UpdatedAt) })
	for _, key := range keys[:len(keys)-maxValues] {
		delete(values, key)
	}
}

// Toggle returns the saved value of the on/off setting with the given key, or defaultValue if it is not saved.
func (store *Store) Toggle(key string, defaultValue bool) bool {
	value, ok := Get[bool](store, namespaceToggles, key)
	if !ok {
		return defaultValue
	}
	return value
}

// SetToggle saves the value of the on/off setting with the given key.
func (store *Store) SetToggle(key string, value bool) {
	Set(store, namespaceToggles, key, value)
}

// scheduleSave saves the state in the background after saveDelay, restarting the delay on every change.
// Must be called with mutex held.
func (store *Store) scheduleSave() {
	if store.saveTimer != nil {
		store.saveTimer.Stop()
	}
	store.saveTimer = time.AfterFunc(store.saveDelay, func() {
		if err := store.Flush(); err != nil {
			logging.Error("Cannot save state: %v", err)
		}
	})
}

// Flush writes pending changes to the file immediately. Called on shutdown, so the last changes are not lost.
// The changes are merged into the file as it is now, so the changes of other instances are kept: a value changed
// by another instance after this one changed it is kept as well.
// This accesses the file system, so it must not be called on the UI thread.
func (store *Store) Flush() error {
	if store == nil {
		return nil
	}
	store.saveMutex.Lock()
	defer store.saveMutex.Unlock()

	store.mutex.Lock()
	if store.saveTimer != nil {
		store.saveTimer.Stop()
		store.saveTimer = nil
	}
	if len(store.changes) == 0 || store.newerVersion != 0 {
		// with a newer version, the changes are kept in memory only
		store.changes = map[change]time.Time{}
		store.mutex.Unlock()
		return nil
	}
	changes := store.changes
	store.changes = map[change]time.Time{}
	pending := make(map[change]*entry, len(changes))
	for changed := range changes {
		if saved, ok := store.namespaces[changed.namespace][changed.key]; ok {
			pending[changed] = &saved
		} else {
			pending[changed] = nil
		}
	}
	limits := make(map[string]int, len(store.limits))
	for ns, limit := range store.limits {
		limits[ns] = limit
	}
	store.mutex.Unlock()

	err := store.merge(changes, pending, limits)
	if err != nil {
		// try again with the next change or flush, unless changed again meanwhile
		store.mutex.Lock()
		for changed, changedAt := range changes {
			if _, ok := store.changes[changed]; !ok {
				store.changes[changed] = changedAt
			}
		}
		store.mutex.Unlock()
	}
	return err
}

// merge writes the pending values (nil for deleted ones) into the file, while holding the lock of the file.
func (store *Store) merge(changes map[change]time.Time, pending map[change]*entry, limits map[string]int) error {
	dir := filepath.Dir(store.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	unlock, err := lockFile(store.path + ".lock")
	if err != nil {
		return err
	}
	defer unlock()

	namespaces, version, err := readFile(store.path)
	if (err == nil || version.version != 0) && !version.isSupported() {
		// written by a newer version meanwhile
		store.mutex.Lock()
		store.setNewerVersion(version.version)
		store.mutex.Unlock()
		return nil
	}
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			logging.Warning("Replacing the unreadable state file: %v", err)
		}
		namespaces = map[string]namespace{}
		version = formatVersion{}
	}
	for changed, changedAt := range changes {
		values := namespaces[changed.namespace]
		if saved, ok := values[changed.key]; ok && saved.UpdatedAt.After(changedAt) {
			// saved by another instance after this one changed it
			continue
		}
		value := pending[changed]
		if value == nil {
			delete(values, changed.key)
			continue
		}
		if values == nil {
			values = namespace{}
			namespaces[changed.namespace] = values
		}
		values[changed.key] = *value
	}
	for ns, limit := range limits {
		applyLimit(namespaces[ns], limit)
	}
	for ns, values := range namespaces {
		if len(values) == 0 {
			delete(namespaces, ns)
		}
	}

	// a newer, compatible version stays marked as such
	content := fileContent{
		Version:           max(version.version, currentVersion),
		CompatibleVersion: max(version.compatible, compatibleVersion),
		Namespaces:        namespaces,
	}
	data, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomically(store.path, append(data, '\n'))
}

// lockFile takes an exclusive lock on the file at path (created if missing), so instances do not merge their
// changes into the state file at the same time. Returns the function releasing it.
func lockFile(path string) (func(), error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, err
	}
	return func() {
		_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		_ = file.Close()
	}, nil
}

// writeFileAtomically replaces the file at path with data, so readers (e.g. another running instance) never see
// a partially written file.
func writeFileAtomically(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	tempFile, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()
	defer func() { _ = os.Remove(tempPath) }()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Sync(); err != nil {
		_ = tempFile.Close()
		return err
	}
	if err := tempFile.Close(); err != nil {
		return err
	}
	return os.Rename(tempPath, path)
}
