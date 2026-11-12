package ui

import (
	"fan2go-tui/internal/state"
	"testing"

	"github.com/rivo/tview"
)

func TestMainPage_PageNavigation(t *testing.T) {
	mainPage := NewMainPage(tview.NewApplication(), state.NewStore())

	pages := []Page{FansPage, CurvesPage, SensorsPage}
	for i, want := range pages {
		if got := mainPage.GetPageAtIndex(i); got != want {
			t.Fatalf("GetPageAtIndex(%d) = %q, want %q", i, got, want)
		}
	}

	if mainPage.page != FansPage {
		t.Fatalf("expected initial page %q, got %q", FansPage, mainPage.page)
	}
	for _, want := range []Page{CurvesPage, SensorsPage, FansPage} {
		mainPage.NextPage()
		if mainPage.page != want {
			t.Fatalf("NextPage: expected %q, got %q", want, mainPage.page)
		}
	}
	for _, want := range []Page{SensorsPage, CurvesPage, FansPage} {
		mainPage.PreviousPage()
		if mainPage.page != want {
			t.Fatalf("PreviousPage: expected %q, got %q", want, mainPage.page)
		}
	}
}
