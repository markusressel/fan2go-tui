package fan

import (
	"fan2go-tui/internal/demo"
	"fan2go-tui/internal/state"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestFansPage_SelectionPreservedAcrossDrawAndRefresh(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init simScreen: %v", err)
	}
	simScreen.SetSize(120, 40)

	app := tview.NewApplication()
	app.SetScreen(simScreen)

	store := state.NewStore()
	sim := demo.NewSimulator()
	store.UpdateFans(sim.GetFans())

	fansPage := NewFansPage(app, store, nil)
	pages := tview.NewPages().AddPage("fans", fansPage.GetLayout(), true, true)
	app.SetRoot(pages, true)
	go func() { _ = app.Run() }()
	defer app.Stop()

	// Initial refresh and scroll to item
	_ = fansPage.Refresh()
	fansPage.ScrollToItem()

	app.Draw()

	initialSelected := fansPage.fanList.GetSelectedItem()
	if initialSelected == nil {
		t.Fatalf("expected initial selection, got nil")
	}
	firstID := initialSelected.FanState.Fan.Config.ID
	t.Logf("initial selected fan: %s", firstID)

	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	time.Sleep(100 * time.Millisecond)
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	time.Sleep(100 * time.Millisecond)

	selectedAfterDown := fansPage.fanList.GetSelectedItem()
	if selectedAfterDown == nil {
		t.Fatalf("expected selection after KeyDown, got nil")
	}
	secondID := selectedAfterDown.FanState.Fan.Config.ID
	t.Logf("selected fan after 2x KeyDown: %s", secondID)
	if secondID == firstID {
		t.Fatalf("selection did not change after KeyDown (still %s)", firstID)
	}

	// Simulate refresh
	sim.Tick(0.25)
	store.UpdateFans(sim.GetFans())
	_ = fansPage.Refresh()

	app.Draw()

	selectedAfterRefresh := fansPage.fanList.GetSelectedItem()
	if selectedAfterRefresh == nil {
		t.Fatalf("expected selection after Refresh, got nil")
	}
	refreshID := selectedAfterRefresh.FanState.Fan.Config.ID
	t.Logf("selected fan after Refresh: %s", refreshID)
	if refreshID != secondID {
		t.Fatalf("expected %s selected after Refresh, got %s", secondID, refreshID)
	}
}
