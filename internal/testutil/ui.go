// Package testutil contains helpers for tests of the UI: tview is not thread safe, so tests only touch widgets on
// the UI thread of the running application.
package testutil

import (
	"testing"
	"time"

	"github.com/rivo/tview"
)

// timeout is how long the helpers wait for the UI thread, or for a condition.
const timeout = 5 * time.Second

// RunApplication runs application in the background. At the end of the test, it is stopped, and the cleanup waits
// until Run returned, so the cleanups registered before (e.g. restoring the configuration) do not race with it.
func RunApplication(t *testing.T, application *tview.Application) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = application.Run()
	}()
	t.Cleanup(func() {
		application.Stop()
		select {
		case <-done:
		case <-time.After(timeout):
			t.Errorf("timeout waiting for the application to stop")
		}
	})
}

// OnUiThread runs f on the UI thread of application, redraws, and waits for it. Must not be called from the UI
// thread.
func OnUiThread(t *testing.T, application *tview.Application, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		application.QueueUpdateDraw(func() {
			f()
			close(done)
		})
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		t.Fatalf("timeout waiting for the UI thread")
	}
}

// WaitFor waits until condition, checked on the UI thread, is true, e.g. after an injected key or a background
// update. Fails the test with message after a timeout.
func WaitFor(t *testing.T, application *tview.Application, message string, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for {
		var fulfilled bool
		OnUiThread(t, application, func() { fulfilled = condition() })
		if fulfilled {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for: %s", message)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
