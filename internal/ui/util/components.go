package util

import (
	"fmt"

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

// SetupReactiveResize calls onResize on the UI thread after the inner size of box changed. The size is checked
// while box is drawn (on the UI thread, tview widgets are not thread safe), and onResize runs after the draw, when
// the children of box are laid out with the new size. Replaces the draw function of box.
func SetupReactiveResize(app *tview.Application, box *tview.Box, onResize func()) {
	tracker := &ResizeTracker{}
	box.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// the inner rect of the last draw is cached, SetRect clears it so it is calculated for the current size
		box.SetRect(x, y, width, height)
		innerX, innerY, innerWidth, innerHeight := box.GetInnerRect()
		if innerWidth > 1 && innerHeight > 1 {
			tracker.OnResize(innerWidth, innerHeight, func() {
				// queued from another goroutine: QueueUpdateDraw waits for the UI thread, which is drawing
				go app.QueueUpdateDraw(onResize)
			})
		}
		return innerX, innerY, innerWidth, innerHeight
	})
}
