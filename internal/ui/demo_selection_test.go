package ui

import (
	"context"
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/demo"
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/testutil"
	"fan2go-tui/internal/ui/dialog"
	"fan2go-tui/internal/ui/fan"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func TestDemoMode_SelectionPreservedAcrossPollTicks(t *testing.T) {
	sim := demo.NewSimulator()
	server := demo.NewServer("127.0.0.1", 0, sim)
	if err := server.Start(); err != nil {
		t.Fatalf("failed to start demo server: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	host := "127.0.0.1"
	port := server.Port()
	apiClient := client.NewApiClient(host, port)
	store := state.NewStore()

	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init simScreen: %v", err)
	}
	simScreen.SetSize(120, 40)

	app := tview.NewApplication()
	app.SetScreen(simScreen)

	mainPage := NewMainPage(app, store)
	connectingDialog := dialog.NewConnectingDialog(client.JoinHostPort(host, port))

	hasReceivedData := false
	dismissedConnectingDialog := false
	var pagesLayout *tview.Pages

	poller := state.NewPoller(apiClient, store, func(err error) {
		app.QueueUpdateDraw(func() {
			if err != nil {
				if !hasReceivedData && !store.HasData() {
					connectingDialog.SetError(err)
					if !dismissedConnectingDialog {
						pagesLayout.ShowPage(string(ConnectingDialog))
					} else {
						mainPage.SetConnectionStatus(false, err.Error())
					}
				} else {
					mainPage.SetConnectionStatus(false, err.Error())
				}
			} else {
				hasReceivedData = true
				pagesLayout.HidePage(string(ConnectingDialog))
				mainPage.SetConnectionStatus(true, "")
			}
			mainPage.Refresh()
		})
	})

	pagesLayout = tview.NewPages().
		AddPage(string(Main), mainPage.layout, true, true).
		AddPage(string(ConnectingDialog), connectingDialog.GetLayout(), true, true)

	app.SetRoot(pagesLayout, true)

	testutil.RunApplication(t, app)

	// Initial fetch
	poller.FetchAndUpdate()
	testutil.OnUiThread(t, app, mainPage.Init)

	// Get FansPage
	var fansPage *fan.FansPage
	testutil.OnUiThread(t, app, func() {
		fansPage, _ = mainPage.GetCurrentPage().(*fan.FansPage)
	})
	if fansPage == nil {
		t.Fatalf("expected FansPage")
	}

	// Press KeyDown twice to navigate to the 3rd fan
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	var selectedAfterDown2 *fan.FanListItemComponent
	testutil.WaitFor(t, app, "third fan selected", func() bool {
		selectedAfterDown2 = fansPage.GetSelectedItem()
		return fansPage.GetSelectedIndex() == 2
	})

	focusBeforeTicks := app.GetFocus()
	t.Logf("Focus before poll ticks: %T (%+v)", focusBeforeTicks, focusBeforeTicks)

	// Now simulate poller ticks
	for i := 0; i < 3; i++ {
		sim.Tick(0.25)
		poller.FetchAndUpdate()
		time.Sleep(100 * time.Millisecond)
	}

	focusAfterTicks := app.GetFocus()
	t.Logf("Focus after poll ticks: %T (%+v)", focusAfterTicks, focusAfterTicks)

	var selectedAfterTicks *fan.FanListItemComponent
	testutil.OnUiThread(t, app, func() { selectedAfterTicks = fansPage.GetSelectedItem() })

	if selectedAfterTicks != selectedAfterDown2 {
		t.Fatalf("Selection was reset! Expected %v, got %v", selectedAfterDown2, selectedAfterTicks)
	}
	if focusAfterTicks != focusBeforeTicks {
		t.Fatalf("Focus was reset! Expected %v, got %v", focusBeforeTicks, focusAfterTicks)
	}
}

func TestRealCreateUi_DemoModeSelection(t *testing.T) {
	sim := demo.NewSimulator()
	server := demo.NewServer("127.0.0.1", 0, sim)
	if err := server.Start(); err != nil {
		t.Fatalf("failed to start demo server: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	config := configuration.CurrentConfig
	t.Cleanup(func() { configuration.CurrentConfig = config })
	configuration.CurrentConfig.Api.Host = "127.0.0.1"
	configuration.CurrentConfig.Api.Port = server.Port()
	configuration.CurrentConfig.Ui.UpdateInterval = 200 * time.Millisecond

	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init simScreen: %v", err)
	}
	simScreen.SetSize(120, 40)

	app := CreateUi(t.Context(), false)
	app.SetScreen(simScreen)

	testutil.RunApplication(t, app)

	// Wait for poller to receive initial data and hide connecting dialog
	time.Sleep(500 * time.Millisecond)

	inspectFocus := func(prefix string) {
		f := app.GetFocus()
		if f == nil {
			t.Logf("%s: nil focus", prefix)
			return
		}
		var title string
		if b, ok := f.(interface{ GetTitle() string }); ok {
			title = b.GetTitle()
		}
		t.Logf("%s: %T (title=%q, ptr=%p)", prefix, f, title, f)
	}

	inspectFocus("Focus before KeyDown")

	// Press KeyDown
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	time.Sleep(100 * time.Millisecond)

	inspectFocus("Focus after KeyDown 1")

	// Press KeyDown again
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	time.Sleep(100 * time.Millisecond)

	inspectFocus("Focus after KeyDown 2")
	savedFocus2 := app.GetFocus()
	// Sleep 600ms to let UpdateTicker fire multiple times
	time.Sleep(600 * time.Millisecond)
	focusAfterTicks := app.GetFocus()

	if focusAfterTicks != savedFocus2 {
		t.Fatalf("Fans focus changed after poll ticks! Expected %p, got %p", savedFocus2, focusAfterTicks)
	}

	// Switch to Curves page (Tab or '2')
	simScreen.InjectKey(tcell.KeyRune, '2', tcell.ModNone)
	time.Sleep(200 * time.Millisecond)

	// Press KeyDown on Curves page
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	time.Sleep(100 * time.Millisecond)
	curvesFocus := app.GetFocus()
	t.Logf("Curves focus after KeyDown: %T (ptr=%p)", curvesFocus, curvesFocus)

	time.Sleep(600 * time.Millisecond)
	curvesFocusAfterTicks := app.GetFocus()
	if curvesFocusAfterTicks != curvesFocus {
		t.Fatalf("Curves focus changed after poll ticks! Expected %p, got %p", curvesFocus, curvesFocusAfterTicks)
	}

	// Switch to Sensors page ('3')
	simScreen.InjectKey(tcell.KeyRune, '3', tcell.ModNone)
	time.Sleep(200 * time.Millisecond)

	// Press KeyDown on Sensors page
	simScreen.InjectKey(tcell.KeyDown, 0, tcell.ModNone)
	time.Sleep(100 * time.Millisecond)
	sensorsFocus := app.GetFocus()
	t.Logf("Sensors focus after KeyDown: %T (ptr=%p)", sensorsFocus, sensorsFocus)

	time.Sleep(600 * time.Millisecond)
	sensorsFocusAfterTicks := app.GetFocus()
	if sensorsFocusAfterTicks != sensorsFocus {
		t.Fatalf("Sensors focus changed after poll ticks! Expected %p, got %p", sensorsFocus, sensorsFocusAfterTicks)
	}
}
