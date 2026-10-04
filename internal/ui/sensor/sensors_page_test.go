package sensor

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/demo"
	"fan2go-tui/internal/state"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestSensorsPage_RefreshAndSelection(t *testing.T) {
	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init simScreen: %v", err)
	}
	simScreen.SetSize(120, 40)

	app := tview.NewApplication()
	app.SetScreen(simScreen)

	store := state.NewStore()
	sim := demo.NewSimulator()
	store.UpdateSensors(sim.GetSensors())

	sensorsPage := NewSensorsPage(app, store)
	pages := tview.NewPages().AddPage("sensors", sensorsPage.GetLayout(), true, true)
	app.SetRoot(pages, true)
	go func() { _ = app.Run() }()
	defer app.Stop()

	_ = sensorsPage.Refresh()
	sensorsPage.ScrollToItem()
	app.Draw()

	initial := sensorsPage.GetSelectedItem()
	if initial == nil {
		t.Fatalf("expected initial sensor selected, got nil")
	}
	initialID := initial.SensorState.Sensor.Config.ID

	// KeyDown to select next
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	time.Sleep(100 * time.Millisecond)

	selectedAfterDown := sensorsPage.GetSelectedItem()
	if selectedAfterDown == nil {
		t.Fatalf("expected selection after KeyDown, got nil")
	}
	nextID := selectedAfterDown.SensorState.Sensor.Config.ID
	if nextID == initialID {
		t.Fatalf("expected different sensor selected after KeyDown")
	}

	// SelectSensorByID
	ok := sensorsPage.SelectSensorByID(initialID)
	if !ok {
		t.Fatalf("expected SelectSensorByID(%s) to succeed", initialID)
	}
	if sensorsPage.GetSelectedItem().SensorState.Sensor.Config.ID != initialID {
		t.Fatalf("expected %s selected after SelectSensorByID", initialID)
	}

	// SelectSensorByID with unknown ID
	if sensorsPage.SelectSensorByID("non-existent-sensor") {
		t.Fatalf("expected SelectSensorByID with non-existent ID to return false")
	}

	// Refresh preserves selection
	sim.Tick(0.25)
	store.UpdateSensors(sim.GetSensors())
	_ = sensorsPage.Refresh()
	app.Draw()

	if sensorsPage.GetSelectedItem().SensorState.Sensor.Config.ID != initialID {
		t.Fatalf("expected selection preserved across refresh")
	}
}

func TestSensorsPage_ShortcutMap(t *testing.T) {
	app := tview.NewApplication()
	store := state.NewStore()
	sensorsPage := NewSensorsPage(app, store)
	shortcuts := sensorsPage.GetShortcutMap()
	if len(shortcuts) == 0 {
		t.Fatalf("expected non-empty shortcut map")
	}
}

func TestSensorListItemComponent_ScrollAndSetSensor(t *testing.T) {
	app := tview.NewApplication()
	sensor := &client.Sensor{
		Name:      "CPU Temp",
		MovingAvg: 45000,
		Config: client.SensorConfig{
			ID: "cpu-temp",
			HwMon: &client.HwMonSensorConfig{
				Platform: "k10temp",
			},
		},
	}
	sensorState := &state.SensorState{Sensor: sensor}

	item := NewSensorListItemComponent(app, sensorState)
	if item.GetLayout() == nil {
		t.Fatalf("expected non-nil layout")
	}

	// Test ScrollHorizontal does not panic
	item.ScrollHorizontal(4)
	item.ScrollHorizontal(-4)

	// Test SetSensor update
	sensor2 := &client.Sensor{
		Name:      "CPU Temp Updated",
		MovingAvg: 50000,
		Config: client.SensorConfig{
			ID: "cpu-temp",
			HwMon: &client.HwMonSensorConfig{
				Platform: "k10temp",
			},
		},
	}
	item.SetSensor(&state.SensorState{Sensor: sensor2})
}
