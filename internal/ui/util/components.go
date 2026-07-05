package util

import (
	"fmt"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Page string

func CreateAttentionText(text string) string {
	return fmt.Sprintf("  %s  ", text)
}

func CreateAttentionTextView(text string) *tview.TextView {
	abortText := CreateAttentionText(text)
	return tview.NewTextView().SetText(abortText).SetTextColor(tcell.ColorYellow).SetTextAlign(tview.AlignRight)
}

type PagesPage interface {
	GetLayout() *tview.Flex
	Refresh() error
}

type CanScrollToItem interface {
	ScrollToItem()
}

type ResizeTracker struct {
	lastW int
	lastH int
}

func (r *ResizeTracker) OnResize(width, height int, action func()) {
	if width > 0 && height > 0 && (width != r.lastW || height != r.lastH) {
		r.lastW = width
		r.lastH = height
		action()
	}
}

func SetupReactiveResize(app *tview.Application, box *tview.Box, onResize func()) {
	tracker := &ResizeTracker{}
	go func() {
		for {
			time.Sleep(50 * time.Millisecond)
			_, _, w, h := box.GetInnerRect()
			if w > 1 && h > 1 {
				tracker.OnResize(w, h, func() {
					app.QueueUpdateDraw(onResize)
				})
			}
		}
	}()
}
