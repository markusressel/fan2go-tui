package ui

import (
	"fan2go-tui/internal/testutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/elliotchance/orderedmap/v3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	uiutil "fan2go-tui/internal/ui/util"
)

func TestApplicationHeaderComponent_WarningCount(t *testing.T) {
	app := tview.NewApplication()
	pagesMap := orderedmap.NewOrderedMap[Page, uiutil.PagesPage]()
	pagesMap.Set(FansPage, &dummyPage{layout: tview.NewFlex()})
	header := NewApplicationHeader(app, pagesMap)

	for _, test := range []struct {
		name      string
		count     int
		supported bool
		expected  string
	}{
		{name: "two warnings", count: 2, supported: true, expected: " ⚠ 2 warnings "},
		{name: "one warning", count: 1, supported: true, expected: " ⚠ 1 warning "},
		{name: "no warnings", count: 0, supported: true, expected: ""},
		{name: "daemon without the warnings API", count: 0, supported: false, expected: ""},
	} {
		t.Run(test.name, func(t *testing.T) {
			header.SetWarningCount(test.count, test.supported)
			if text := header.warningCountTextView.GetText(true); text != test.expected {
				t.Fatalf("expected the warning badge %q, got %q", test.expected, text)
			}
		})
	}
}

// screenText returns the text on the screen, one line per row.
func screenText(screen tcell.SimulationScreen) string {
	cells, width, _ := screen.GetContents()
	var text strings.Builder
	for i, cell := range cells {
		if len(cell.Runes) > 0 {
			text.WriteString(string(cell.Runes))
		} else {
			text.WriteRune(' ')
		}
		if (i+1)%width == 0 {
			text.WriteRune('\n')
		}
	}
	return text.String()
}

// The demo has a warning on the fan case-exhaust and on the sensor nvme-temp: the header shows their count, the
// list marks the fan and shows its warning.
func TestUi_ShowsWarningsOfTheDemo(t *testing.T) {
	startDemoServer(t)
	ui := startUi(t, filepath.Join(t.TempDir(), "state.json"))

	testutil.WaitFor(t, ui.app, "the warning count in the header", func() bool {
		return ui.mainPage.header.warningCountTextView.GetText(true) == " ⚠ 2 warnings "
	})
	testutil.WaitFor(t, ui.app, "the marked fan and its warning", func() bool {
		text := screenText(ui.screen)
		return strings.Contains(text, "⚠ case-exhaust") && strings.Contains(text, "⚠ RPM stays at ~820")
	})
}
