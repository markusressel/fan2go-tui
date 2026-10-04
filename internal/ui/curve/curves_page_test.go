package curve

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/demo"
	"fan2go-tui/internal/state"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestCurvesPage_RefreshAndSelection(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init simScreen: %v", err)
	}
	simScreen.SetSize(120, 40)

	app := tview.NewApplication()
	app.SetScreen(simScreen)

	store := state.NewStore()
	sim := demo.NewSimulator()
	store.UpdateCurves(sim.GetCurves())

	curvesPage := NewCurvesPage(app, store, nil, nil)
	pages := tview.NewPages().AddPage("curves", curvesPage.GetLayout(), true, true)
	app.SetRoot(pages, true)
	go func() { _ = app.Run() }()
	defer app.Stop()

	_ = curvesPage.Refresh()
	curvesPage.ScrollToItem()
	app.Draw()

	initial := curvesPage.GetSelectedItem()
	if initial == nil {
		t.Fatalf("expected initial curve selected, got nil")
	}
	initialID := initial.CurveState.Curve.Config.ID

	// KeyDown to select next
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	time.Sleep(100 * time.Millisecond)

	selectedAfterDown := curvesPage.GetSelectedItem()
	if selectedAfterDown == nil {
		t.Fatalf("expected selection after KeyDown, got nil")
	}
	nextID := selectedAfterDown.CurveState.Curve.Config.ID
	if nextID == initialID {
		t.Fatalf("expected different curve selected after KeyDown")
	}

	// SelectCurveByID
	ok := curvesPage.SelectCurveByID(initialID)
	if !ok {
		t.Fatalf("expected SelectCurveByID(%s) to succeed", initialID)
	}
	if curvesPage.GetSelectedItem().CurveState.Curve.Config.ID != initialID {
		t.Fatalf("expected %s selected after SelectCurveByID", initialID)
	}

	// SelectCurveByID with unknown ID
	if curvesPage.SelectCurveByID("non-existent-curve") {
		t.Fatalf("expected SelectCurveByID with non-existent ID to return false")
	}

	// Refresh preserves selection
	sim.Tick(0.25)
	store.UpdateCurves(sim.GetCurves())
	_ = curvesPage.Refresh()
	app.Draw()

	if curvesPage.GetSelectedItem().CurveState.Curve.Config.ID != initialID {
		t.Fatalf("expected selection preserved across refresh")
	}
}

func TestCurvesPage_ShortcutMap(t *testing.T) {
	app := tview.NewApplication()
	store := state.NewStore()
	curvesPage := NewCurvesPage(app, store, nil, nil)
	shortcuts := curvesPage.GetShortcutMap()
	if len(shortcuts) == 0 {
		t.Fatalf("expected non-empty shortcut map")
	}
}

func TestCurveListItemComponent_ScrollAndSetCurve(t *testing.T) {
	app := tview.NewApplication()
	curve := &client.Curve{
		Value: 50.0,
		Config: client.CurveConfig{
			ID: "test-curve",
			Linear: &client.LinearCurveConfig{
				Sensor: "cpu-temp",
				Min:    40,
				Max:    80,
			},
		},
	}
	curveState := &state.CurveState{Curve: curve}

	item := NewCurveListItemComponent(app, curveState, nil, nil)
	if item.GetLayout() == nil {
		t.Fatalf("expected non-nil layout")
	}

	// Test ScrollHorizontal does not panic
	item.ScrollHorizontal(4)
	item.ScrollHorizontal(-4)

	// Test SetCurve update
	curve2 := &client.Curve{
		Value: 75.0,
		Config: client.CurveConfig{
			ID: "test-curve",
			Linear: &client.LinearCurveConfig{
				Sensor: "cpu-temp",
				Min:    30,
				Max:    90,
			},
		},
	}
	item.SetCurve(&state.CurveState{Curve: curve2})
}
