package ui

import (
	"fan2go-tui/internal/ui/status_message"
	uiutil "fan2go-tui/internal/ui/util"
	"strings"
	"testing"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/rivo/tview"
)

type dummyPage struct {
	layout *tview.Flex
}

func (d *dummyPage) GetLayout() *tview.Flex {
	return d.layout
}

func (d *dummyPage) Refresh() error {
	return nil
}

func TestApplicationHeaderComponent_ConnectionStatus(t *testing.T) {
	app := tview.NewApplication()
	pagesMap := orderedmap.NewOrderedMap[Page, uiutil.PagesPage]()
	pagesMap.Set(FansPage, &dummyPage{layout: tview.NewFlex()})

	header := NewApplicationHeader(app, *pagesMap)

	if !header.connected {
		t.Fatalf("expected initially connected to be true")
	}
	if text := header.connectionStatusTextView.GetText(true); text != "" {
		t.Fatalf("expected empty connectionStatusTextView initially, got %q", text)
	}

	// Disconnect
	errMsg := "Cannot reach fan2go daemon"
	header.SetConnectionStatus(false, errMsg)

	if header.connected {
		t.Fatalf("expected connected to be false")
	}
	if text := header.connectionStatusTextView.GetText(true); !strings.Contains(text, "DISCONNECTED") {
		t.Fatalf("expected DISCONNECTED badge text, got %q", text)
	}
	if text := header.statusTextView.GetText(true); !strings.Contains(text, errMsg) {
		t.Fatalf("expected statusTextView to contain %q, got %q", errMsg, text)
	}

	// Normal status clear while disconnected should NOT clear disconnection error
	header.SetStatus(status_message.NewInfoStatusMessage(""))
	if text := header.statusTextView.GetText(true); !strings.Contains(text, errMsg) {
		t.Fatalf("expected statusTextView to keep error message, got %q", text)
	}

	header.ResetStatus()
	if text := header.statusTextView.GetText(true); !strings.Contains(text, errMsg) {
		t.Fatalf("expected statusTextView to keep error message after ResetStatus, got %q", text)
	}

	// Reconnect
	header.SetConnectionStatus(true, "")
	if !header.connected {
		t.Fatalf("expected connected to be true")
	}
	if text := header.connectionStatusTextView.GetText(true); text != "" {
		t.Fatalf("expected empty connectionStatusTextView after reconnect, got %q", text)
	}
	if text := header.statusTextView.GetText(true); text != "" {
		t.Fatalf("expected cleared statusTextView after reconnect, got %q", text)
	}
}
