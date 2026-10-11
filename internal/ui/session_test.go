package ui

import (
	"context"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/demo"
	"fan2go-tui/internal/ui/curve"
	"fan2go-tui/internal/ui/fan"
	"fan2go-tui/internal/ui/shortcut_helper"
	"fan2go-tui/internal/uistate"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// startDemoServer starts a simulated fan2go daemon and configures the UI to connect to it.
func startDemoServer(t *testing.T) {
	t.Helper()
	server := demo.NewServer("127.0.0.1", 0, demo.NewSimulator())
	if err := server.Start(); err != nil {
		t.Fatalf("failed to start demo server: %v", err)
	}
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	})

	config := configuration.CurrentConfig
	t.Cleanup(func() { configuration.CurrentConfig = config })
	configuration.CurrentConfig.Api.Host = "127.0.0.1"
	configuration.CurrentConfig.Api.Port = server.Port()
	configuration.CurrentConfig.Ui.UpdateInterval = 200 * time.Millisecond
	configuration.CurrentConfig.Ui.RestoreSession = true
}

// runningUi is the UI of a run of the application.
type runningUi struct {
	t        *testing.T
	app      *tview.Application
	mainPage *MainPage
	screen   tcell.SimulationScreen
}

// startUi starts the UI with the state at statePath, like the application does.
func startUi(t *testing.T, statePath string) *runningUi {
	t.Helper()
	uistate.Current = uistate.Load(statePath)
	t.Cleanup(func() { uistate.Current = nil })

	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatalf("failed to init screen: %v", err)
	}
	app, mainPage := createUi(true)
	app.SetScreen(screen)
	screen.SetSize(120, 40)
	go func() { _ = app.Run() }()
	t.Cleanup(app.Stop)
	return &runningUi{t: t, app: app, mainPage: mainPage, screen: screen}
}

// stop quits the application and saves the state, like the application does.
func (ui *runningUi) stop() {
	ui.app.Stop()
	if err := uistate.Current.Flush(); err != nil {
		ui.t.Fatalf("cannot save state: %v", err)
	}
}

// onUiThread runs f on the UI thread and waits for it.
func (ui *runningUi) onUiThread(f func()) {
	done := make(chan struct{})
	ui.app.QueueUpdate(func() {
		f()
		close(done)
	})
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		ui.t.Fatalf("timeout waiting for the UI thread")
	}
}

// selection returns the shown page and the ID of the selected entry of the fans and curves pages.
func (ui *runningUi) selection() (page Page, fanID string, curveID string) {
	ui.onUiThread(func() {
		page = ui.mainPage.page
		fansPage, _ := ui.mainPage.pagesMap.Get(FansPage)
		if item := fansPage.(*fan.FansPage).GetSelectedItem(); item != nil {
			fanID = item.FanState.Fan.Config.ID
		}
		curvesPage, _ := ui.mainPage.pagesMap.Get(CurvesPage)
		if item := curvesPage.(*curve.CurvesPage).GetSelectedItem(); item != nil {
			curveID = item.CurveState.Curve.Config.ID
		}
	})
	return page, fanID, curveID
}

// waitFor waits until condition is true for the shown page and the selections, see selection.
func (ui *runningUi) waitFor(message string, condition func(page Page, fanID string, curveID string) bool) {
	ui.t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		page, fanID, curveID := ui.selection()
		if condition(page, fanID, curveID) {
			return
		}
		if time.Now().After(deadline) {
			ui.t.Fatalf("%s: page %s, fan %q, curve %q", message, page, fanID, curveID)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func (ui *runningUi) pressKey(key tcell.Key, r rune) {
	ui.screen.InjectKey(key, r, tcell.ModNone)
}

// The page and the selections of the pages are restored when connecting to the same fan2go daemon again.
func TestSession_IsRestored(t *testing.T) {
	startDemoServer(t)
	statePath := filepath.Join(t.TempDir(), "state.json")

	ui := startUi(t, statePath)
	ui.waitFor("fans loaded", func(page Page, fanID string, _ string) bool {
		return page == FansPage && fanID == "case-exhaust"
	})
	ui.pressKey(tcell.KeyDown, 0)
	ui.pressKey(tcell.KeyDown, 0)
	ui.waitFor("third fan selected", func(_ Page, fanID string, _ string) bool { return fanID == "cpu-fan" })
	ui.pressKey(tcell.KeyRune, '2')
	ui.waitFor("curves shown", func(page Page, _ string, curveID string) bool {
		return page == CurvesPage && curveID == "case-function-curve"
	})
	ui.pressKey(tcell.KeyDown, 0)
	ui.waitFor("second curve selected", func(_ Page, _ string, curveID string) bool {
		return curveID == "cpu-linear-curve"
	})
	ui.stop()

	ui = startUi(t, statePath)
	ui.waitFor("curves page and curve restored", func(page Page, _ string, curveID string) bool {
		return page == CurvesPage && curveID == "cpu-linear-curve"
	})
	ui.pressKey(tcell.KeyRune, '1')
	ui.waitFor("fan restored", func(page Page, fanID string, _ string) bool {
		return page == FansPage && fanID == "cpu-fan"
	})
	ui.stop()

	// with --fresh, the first page and entries are shown, and the new selections are remembered
	configuration.CurrentConfig.Ui.RestoreSession = false
	ui = startUi(t, statePath)
	ui.waitFor("fresh start", func(page Page, fanID string, _ string) bool {
		return page == FansPage && fanID == "case-exhaust"
	})
	ui.pressKey(tcell.KeyDown, 0)
	ui.waitFor("second fan selected", func(_ Page, fanID string, _ string) bool { return fanID == "case-intake" })
	ui.stop()

	configuration.CurrentConfig.Ui.RestoreSession = true
	ui = startUi(t, statePath)
	ui.waitFor("selection of the fresh start restored", func(page Page, fanID string, _ string) bool {
		return page == FansPage && fanID == "case-intake"
	})
	ui.stop()
}

// visibleFans returns the indices of the first and the last fan shown.
func (ui *runningUi) visibleFans() (first int, last int) {
	ui.onUiThread(func() {
		fansPage, _ := ui.mainPage.pagesMap.Get(FansPage)
		first, last = fansPage.(*fan.FansPage).GetVisibleRange()
	})
	return first, last
}

// waitForVisibleFans waits until the fans first to last are shown.
func (ui *runningUi) waitForVisibleFans(message string, first int, last int) {
	ui.t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		shownFirst, shownLast := ui.visibleFans()
		if shownFirst == first && shownLast == last {
			return
		}
		if time.Now().After(deadline) {
			ui.t.Fatalf("%s: expected fans %d-%d to be shown, got %d-%d", message, first, last, shownFirst, shownLast)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

// The scroll position of the list is restored with the selection, not only scrolled far enough to show it.
func TestSession_ScrollPositionIsRestored(t *testing.T) {
	startDemoServer(t)
	statePath := filepath.Join(t.TempDir(), "state.json")

	ui := startUi(t, statePath)
	ui.waitFor("fans loaded", func(_ Page, fanID string, _ string) bool { return fanID == "case-exhaust" })
	ui.waitForVisibleFans("first fans shown", 0, 1)
	for range 3 {
		ui.pressKey(tcell.KeyDown, 0)
	}
	ui.waitFor("last fan selected", func(_ Page, fanID string, _ string) bool { return fanID == "gpu-fan" })
	ui.pressKey(tcell.KeyUp, 0)
	ui.waitFor("third fan selected", func(_ Page, fanID string, _ string) bool { return fanID == "cpu-fan" })
	ui.waitForVisibleFans("last fans shown", 2, 3)
	ui.stop()

	ui = startUi(t, statePath)
	ui.waitFor("fan restored", func(_ Page, fanID string, _ string) bool { return fanID == "cpu-fan" })
	ui.waitForVisibleFans("scroll position restored", 2, 3)
	// not undone by the next update
	time.Sleep(3 * configuration.CurrentConfig.Ui.UpdateInterval)
	ui.waitForVisibleFans("scroll position kept", 2, 3)
	ui.stop()
}

// The session of another fan2go daemon is not restored.
func TestSession_IsRememberedPerServer(t *testing.T) {
	startDemoServer(t)
	statePath := filepath.Join(t.TempDir(), "state.json")

	ui := startUi(t, statePath)
	ui.waitFor("fans loaded", func(_ Page, fanID string, _ string) bool { return fanID == "case-exhaust" })
	ui.pressKey(tcell.KeyRune, '2')
	ui.waitFor("curves shown", func(page Page, _ string, _ string) bool { return page == CurvesPage })
	ui.stop()

	// another daemon
	startDemoServer(t)
	ui = startUi(t, statePath)
	ui.waitFor("first page shown", func(page Page, fanID string, _ string) bool {
		return page == FansPage && fanID == "case-exhaust"
	})
	ui.stop()
}

func TestShortcutVisibility_IsRemembered(t *testing.T) {
	statePath := filepath.Join(t.TempDir(), "state.json")
	uistate.Current = uistate.Load(statePath)
	t.Cleanup(func() { uistate.Current = nil })

	shortcut_helper.InitShortcutVisibility()
	if shortcut_helper.ShortcutsHidden() {
		t.Fatalf("expected the shortcuts to be shown by default")
	}
	shortcut_helper.ToggleShortcuts()
	if err := uistate.Current.Flush(); err != nil {
		t.Fatal(err)
	}

	uistate.Current = uistate.Load(statePath)
	shortcut_helper.InitShortcutVisibility()
	if !shortcut_helper.ShortcutsHidden() {
		t.Errorf("expected the shortcuts to stay hidden")
	}
}

func TestNewerStateVersionIsShown(t *testing.T) {
	startDemoServer(t)
	statePath := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(statePath, []byte(`{"version": 7, "compatibleVersion": 7}`), 0o600); err != nil {
		t.Fatal(err)
	}

	ui := startUi(t, statePath)
	ui.waitFor("fans loaded", func(_ Page, fanID string, _ string) bool { return fanID == "case-exhaust" })
	var status string
	ui.onUiThread(func() { status = ui.mainPage.header.statusTextView.GetText(true) })
	if !strings.Contains(status, "newer fan2go-tui (format 7)") {
		t.Errorf("expected the warning in the header, got %q", status)
	}
	ui.stop()
}
