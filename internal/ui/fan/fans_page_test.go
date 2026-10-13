package fan

import (
	"fan2go-tui/internal/demo"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/testutil"
	"testing"

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
	testutil.RunApplication(t, app)

	var firstID string
	testutil.OnUiThread(t, app, func() {
		// Initial refresh and scroll to item
		_ = fansPage.Refresh()
		fansPage.ScrollToItem()
		if initialSelected := fansPage.fanList.GetSelectedItem(); initialSelected != nil {
			firstID = initialSelected.FanState.Fan.Config.ID
		}
	})
	if firstID == "" {
		t.Fatalf("expected initial selection, got nil")
	}

	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	var secondID string
	testutil.WaitFor(t, app, "third fan selected after 2x KeyDown", func() bool {
		if fansPage.fanList.GetSelectedIndex() != 2 {
			return false
		}
		secondID = fansPage.fanList.GetSelectedItem().FanState.Fan.Config.ID
		return true
	})
	if secondID == firstID {
		t.Fatalf("selection did not change after KeyDown (still %s)", firstID)
	}

	// Simulate refresh
	sim.Tick(0.25)
	store.UpdateFans(sim.GetFans())
	testutil.OnUiThread(t, app, func() {
		_ = fansPage.Refresh()
	})
	testutil.OnUiThread(t, app, func() {
		selectedAfterRefresh := fansPage.fanList.GetSelectedItem()
		if selectedAfterRefresh == nil {
			t.Errorf("expected selection after Refresh, got nil")
		} else if refreshID := selectedAfterRefresh.FanState.Fan.Config.ID; refreshID != secondID {
			t.Errorf("expected %s selected after Refresh, got %s", secondID, refreshID)
		}
	})
}
