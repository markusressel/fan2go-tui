package dialog

import (
	"errors"
	"strings"
	"testing"
)

func TestConnectingDialog(t *testing.T) {
	d := NewConnectingDialog("127.0.0.1:9001")

	if d.GetName() != "connecting" {
		t.Fatalf("expected name 'connecting', got %s", d.GetName())
	}

	if d.GetLayout() == nil {
		t.Fatalf("expected non-nil layout")
	}

	initialText := d.detailsTextView.GetText(false)
	if !strings.Contains(initialText, "Waiting for initial connection") {
		t.Fatalf("unexpected initial text: %s", initialText)
	}

	d.SetError(errors.New("connection refused"))
	if d.LastError() == nil || d.LastError().Error() != "connection refused" {
		t.Fatalf("expected LastError to return 'connection refused', got %v", d.LastError())
	}
	errText := d.detailsTextView.GetText(false)
	if !strings.Contains(errText, "connection refused") {
		t.Fatalf("expected error text to contain 'connection refused', got: %s", errText)
	}
	if !strings.Contains(errText, "Retrying automatically") {
		t.Fatalf("expected error text to contain 'Retrying automatically', got: %s", errText)
	}

	d.SetConnecting()
	if d.LastError() != nil {
		t.Fatalf("expected nil LastError after SetConnecting, got %v", d.LastError())
	}
	resetText := d.detailsTextView.GetText(false)
	if !strings.Contains(resetText, "Waiting for initial connection") {
		t.Fatalf("unexpected reset text: %s", resetText)
	}
}
