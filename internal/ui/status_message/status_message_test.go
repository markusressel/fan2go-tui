package status_message

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

func TestStatusMessage_Constructors(t *testing.T) {
	info := NewInfoStatusMessage("info msg")
	if info.Message != "info msg" || info.Color != tcell.ColorLightGray || info.Duration != StatusMessageDurationInfinite {
		t.Errorf("unexpected info message: %+v", info)
	}

	warn := NewWarningStatusMessage("warn msg")
	if warn.Message != "warn msg" || warn.Color != tcell.ColorYellow || warn.Duration != StatusMessageDurationInfinite {
		t.Errorf("unexpected warn message: %+v", warn)
	}

	err := NewErrorStatusMessage("err msg")
	if err.Message != "err msg" || err.Color != tcell.ColorRed || err.Duration != StatusMessageDurationInfinite {
		t.Errorf("unexpected err message: %+v", err)
	}

	custom := info.SetDuration(5 * time.Second).SetColor(tcell.ColorGreen)
	if custom.Duration != 5*time.Second || custom.Color != tcell.ColorGreen {
		t.Errorf("unexpected custom message: %+v", custom)
	}
}
