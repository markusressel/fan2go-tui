package dialog

import (
	uiutil "fan2go-tui/internal/ui/util"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Dialog interface {
	GetName() string
	GetLayout() *tview.Flex
}

type DialogOptionId int

type DialogOption struct {
	Id   DialogOptionId
	Name string
}

func createModal(title string, content tview.Primitive, width int, height int) *tview.Flex {
	dialogFrame := tview.NewFlex()
	dialogFrame.SetBorder(true)
	uiutil.SetupDialogWindow(dialogFrame, title)
	dialogFrame.AddItem(content, 0, 1, true)

	dialogContentColumnWrapper := tview.NewFlex()
	dialogContentColumnWrapper.AddItem(nil, 0, 1, false)

	dialogContentRowWrapper := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 0, 1, false).
		AddItem(dialogFrame, height, 1, true).
		AddItem(nil, 0, 1, false)

	dialogContentColumnWrapper.
		AddItem(dialogContentRowWrapper, width, 1, true).
		AddItem(nil, 0, 1, false)

	// tview hands mouse events to all pages until one consumes them, and neither the empty space around the dialog
	// nor its border do: without this, a click next to the dialog or on its border would reach the page behind it
	// (https://github.com/rivo/tview/issues/926). The content consumes the events on it itself.
	dialogContentColumnWrapper.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		x, y := event.Position()
		contentX, contentY, width, height := content.GetRect()
		if x >= contentX && x < contentX+width && y >= contentY && y < contentY+height {
			return action, event
		}
		if action == tview.MouseMove {
			// not consumed, so moving the mouse does not redraw the screen
			return action, nil
		}
		return tview.MouseConsumed, nil
	})

	return dialogContentColumnWrapper
}
