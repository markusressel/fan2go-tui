package ui

import (
	"context"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/testutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Regression: the poller kept polling the fan2go daemon until the process exited, e.g. after the UI of a test
// stopped.
func TestCreateUi_StopsPollingWhenCancelled(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("{}"))
	}))
	t.Cleanup(server.Close)
	host, port, err := net.SplitHostPort(server.Listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	config := configuration.CurrentConfig
	t.Cleanup(func() { configuration.CurrentConfig = config })
	configuration.CurrentConfig.Api.Host = host
	configuration.CurrentConfig.Api.Port, _ = strconv.Atoi(port)
	configuration.CurrentConfig.Ui.UpdateInterval = 20 * time.Millisecond

	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	app := CreateUi(ctx, false)
	app.SetScreen(screen)
	testutil.RunApplication(t, app)

	deadline := time.Now().Add(5 * time.Second)
	for requests.Load() < 6 { // fans, curves and sensors, twice
		if time.Now().After(deadline) {
			t.Fatalf("expected the daemon to be polled, got %d requests", requests.Load())
		}
		time.Sleep(10 * time.Millisecond)
	}

	cancel()
	time.Sleep(50 * time.Millisecond) // a poll that was running when cancelled
	polled := requests.Load()
	time.Sleep(10 * configuration.CurrentConfig.Ui.UpdateInterval)
	if requests.Load() != polled {
		t.Errorf("expected no more polls after cancelling, got %d more", requests.Load()-polled)
	}
}

// Regression: after the application stopped, the poller waited forever for it to run its update.
func TestQueueUpdateDraw_StopsWaitingWhenCancelled(t *testing.T) {
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	app := tview.NewApplication().SetScreen(screen).SetRoot(tview.NewBox(), true)
	stopApplication := testutil.RunApplication(t, app)

	ctx, cancel := context.WithCancel(t.Context())
	ran := false
	queueUpdateDraw(ctx, app, func() { ran = true })
	if !ran {
		t.Fatalf("expected the update to run while the application runs")
	}

	stopApplication()
	returned := make(chan struct{})
	go func() {
		queueUpdateDraw(ctx, app, func() {})
		close(returned)
	}()
	cancel()
	select {
	case <-returned:
	case <-time.After(5 * time.Second):
		t.Fatalf("expected queueUpdateDraw to return when cancelled")
	}
}
