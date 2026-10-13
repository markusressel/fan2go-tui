# AGENTS.md

Notes for working on this repository: how it is structured, the rules that are easy to get wrong, and the conventions
of the code. fan2go-tui is a terminal UI (tview/tcell) for the REST API of the [fan2go](https://github.com/markusressel/fan2go)
daemon: it shows the fans, curves and sensors of the daemon, with graphs of their recent values.

## Build and test

- `just build` (to `bin/fan2go-tui`), `just test`, `just coverage`. Or directly: `go build ./...`, `go vet ./...`,
  `go test ./...`. The `Makefile` only forwards to the `Justfile`.
- `just demo` (or `bin/fan2go-tui demo --ui`) runs the UI against a simulated daemon (`internal/demo`). Use it to try
  changes: no fan2go daemon or hardware is needed. `demo` alone only runs the simulated API (default `127.0.0.1:9001`),
  for a second terminal running `fan2go-tui`. `--port 0` picks a free port.
- Before committing: `gofmt -l .` is empty, `go vet ./...`, `go test ./...` and `just test-race` pass. CI
  (`.github/workflows/go.yml`) runs `just test` and `just test-race`, and builds; codespell checks the spelling (see
  `.codespellignore`).
- Tests never need a fan2go daemon: UI tests start the simulated one (`demo.NewServer("127.0.0.1", 0, ...)`).
- The log is written to `/tmp/fan2go-tui.log` (stdout belongs to the UI).

## Structure

- `main.go`, `cmd/`: the cobra commands. `root.go` loads the configuration and starts the UI, `demo.go` the simulated
  daemon (with `--ui` also the UI), `version.go`. `cmd/global` holds the flags and the version set by `-ldflags`.
- `internal/application.go`: `RunApplication` runs the UI (and the optional pprof server) in an `oklog/run` group, and
  saves `uistate` on exit.
- `internal/configuration/`: the user-written configuration (viper; `fan2go-tui.yaml` in `.`, `~`, `~/.config`,
  `~/.config/fan2go-tui`, `/etc/fan2go-tui`), in `configuration.CurrentConfig`. The documented example is
  `fan2go-tui.yaml` in the repository root; keep it in sync with new options.
- `internal/client/`: the API client (`Fan2goApiClient`: `GetFans`, `GetCurves`, `GetSensors`, ...) and the types of
  the API. They mirror the JSON of the fan2go daemon, so field names and shapes follow fan2go, not this project.
- `internal/state/`: the data of the daemon. `Poller.FetchAndUpdate` fetches fans, curves and sensors in parallel and
  appends their values to the history of `Store` (`FanState`, `CurveState`, `SensorState`, at most 1000 values).
- `internal/uistate/`: what the application remembers between runs (`$XDG_STATE_HOME/fan2go-tui/state.json`): the
  page and the selections per daemon address, toggles. Namespaced JSON values, saving is debounced, several instances
  share the file (see the package comment). `uistate.Current` is nil in tests unless set, then nothing is saved.
- `internal/demo/`: the simulated daemon. `Simulator` models sensors, curves (linear, PID, function, ...) and fans
  (`Tick`), `Server` serves them with the same JSON as fan2go. When the client types change, update the simulator
  so the demo shows the new data.
- `internal/ui/`: the UI.
  - `root.go`: `CreateUi` creates the application, the poller loop (`UpdateTicker`, the refresh interval, `+`/`-`),
    the connecting dialog, global keys (`ctrl+q`, `?`/`F1`, `⭾`/`shift+⭾`).
  - `main_page.go`: `MainPage`, the header, the pages (`FansPage`, `CurvesPage`, `SensorsPage`; keys `1`-`3`) and the
    shortcut map. `session.go` remembers the page and selections in `uistate`.
  - One package per page: `fan/`, `curve/`, `sensor/`. Each has a page (`*_page.go`, a `util.ListComponent` of list
    items), list items (info + graph components), and graph overlays (`overlays.go`). Pages implement
    `util.PagesPage` and `shortcut_helper.ShortcutMapProvider`.
  - `graph/`: the graph widget. `GraphComponent` draws `GraphLine`/`GraphBar` series on an `OverlayPlot` (extends
    `tvxwidgets.Plot`), with overlays (`HLine`, `VLine`, `Marker`, `Trail`, `Heatmap`, axis labels, legend) that
    implement `GraphComponentOverlay`. `graph_test/` is a page for trying graphs by hand (disabled in `main_page.go`).
  - `util/`: `ListComponent[T]` (the scrollable lists of the pages, with `ListSelection`/`SelectionMemory` to keep
    and restore the selection), `RowSelectionTable[T]`, `ScrollbarComponent`, windows (`SetupWindow`), colors.
  - `txwidgets/`: `ConfigInfoComponent` (the info panes) and `config_sections.go` (how the configuration of a fan,
    curve or sensor is shown), `Span`/`ColorTag` for colored text.
  - `shortcut_helper/`: shortcut entries, key names, the shortcut bar at the bottom (hidden with `?`).
  - `dialog/`, `status_message/`, `header.go`: the connecting dialog, status messages in the header.
  - `theme/`: all colors (`theme.Colors`), types in `common.go`, values in `dark.go`.

## Data flow and threads

tview is not thread safe, and every update redraws the whole screen. Most UI bugs break one of these rules:

- **The poller runs in the background.** `root.go` calls `poller.FetchAndUpdate()` on every tick of `UpdateTicker`
  in a goroutine, until the context given to `CreateUi` is done (cancel it when the application stopped). Its
  callback switches to the UI thread with `queueUpdateDraw` (`QueueUpdateDraw` that stops waiting when the context is
  done: a stopped application runs no updates) and calls `MainPage.Refresh`, which refreshes only the shown page; the
  page reads the `Store`.
- **Widgets only on the UI thread:** in input/mouse handlers, draw functions, or `QueueUpdateDraw`. Never call
  `QueueUpdateDraw` (or `QueueUpdate`) from the UI thread itself: it blocks forever. Do not block the UI thread with
  network calls; fetch in the background like the poller.
- **`state.Store` is shared** between the poller and the UI: read it only with its getters. They hold the lock and
  return copies (`GetFanState`, ...), as the poller keeps appending to the value histories.
- **Sizes are read while drawing.** `util.SetupReactiveResize` checks the size of a box in its draw function and
  queues the update from a goroutine; do not poll widget sizes from other goroutines.
- **Keep the selection across refreshes.** Pages rebuild their list items from the store on every tick; the
  selected entry, its row and the scroll position must not jump (see `ListSelection`, and the tests named
  `...SelectionPreserved...`). Entries are identified by their configuration ID.
- **Graph sizes must not flicker.** Values that decide the layout (axis widths, labels) must not change between
  redraws of the same data, e.g. depending on which plot columns hit a measured value (see the fix in #128).

## Shortcuts

- Each page returns its entries (`GetShortcutMap`); `MainPage.updateShortcutMap` appends the global ones.
- Key names follow `shortcut_helper/key_names.go` (characters as typed, named keys capitalized, modifiers lowercase:
  `ctrl+q`, `shift+⭾`, `F1`). `TestShortcutStyle` checks every `KeyCombo` in `internal/ui`.
- Entries are grouped (`GroupAction`, `GroupView`, `GroupNavigation`, `GroupGlobal`); each group has its key color.

## Code style

- Comments are full sentences that explain why, not what, and start with the name of what they describe. New and
  changed code follows this, older code does not always.
- Match the surrounding code: naming, receiver names (`mainPage`, `c` for components, `applicationHeader`), file
  layout (`New...`, `createLayout`, `GetLayout`, `refresh`/`Refresh`, setters).
- Escape text for tview with `tview.Escape`, or use `txwidgets.Span(color, format, args...)`, which escapes its
  arguments itself (do not escape twice). Text from the daemon (IDs, paths, commands) must be escaped.
- New UI elements take their colors from `theme.Colors`, not `tcell.Color...` constants.
- Packages must not import upwards: `ui/util`, `graph`, `txwidgets` do not import pages, pages do not import `ui`.

## Tests

- Standard `testing` (no testify). Table tests where several inputs share an assertion.
- UI tests run a real `tview.Application` on `tcell.NewSimulationScreen("UTF-8")` and talk to the simulated daemon
  (see `startDemoServer`/`startUi` in `internal/ui/session_test.go`). `internal/testutil`: `RunApplication` (stops the
  application at the end of the test, or with the returned function, and waits for it; do not call `Stop` again
  after `Run` returned), `OnUiThread` (read or change widgets only there), `WaitFor` (instead of sleeping after an
  injected key or for a background update).
- Tests that change `configuration.CurrentConfig` or `uistate.Current` restore them with `t.Cleanup`.
- Regression tests start their comment with `// Regression:` (with the issue, if there is one) and fail without the
  fix.

## Commits and pull requests

- Commit messages: a short summary in the imperative, sentence case (e.g. "Restore the scroll position of the lists
  between runs"), then a body that explains why. Older commits use `feat:`/`fix:` prefixes; new ones do not.
- Branches: `feature/...`, `fix/...`. Pull requests: `## Summary` (what and why), `## Notes` (known issues, things
  not done), `## Test plan` (checkboxes; only check what was actually done). Label `enhancement` or `bug`, used by
  release-drafter.
- The README is short: features, installation, configuration, demo mode. Details belong into the application
  (shortcut bar) or `fan2go-tui.yaml`.
