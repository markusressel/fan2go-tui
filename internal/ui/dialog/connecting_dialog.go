package dialog

import (
	"fmt"

	"github.com/rivo/tview"
)

type ConnectingDialog struct {
	layout          *tview.Flex
	serverAddress   string
	statusTextView  *tview.TextView
	detailsTextView *tview.TextView
	lastError       error
}

func NewConnectingDialog(serverAddress string) *ConnectingDialog {
	d := &ConnectingDialog{
		serverAddress: serverAddress,
	}
	d.createLayout()
	return d
}

func (d *ConnectingDialog) GetName() string {
	return "connecting"
}

func (d *ConnectingDialog) LastError() error {
	return d.lastError
}

func (d *ConnectingDialog) createLayout() {
	contentLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	statusTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	statusText := fmt.Sprintf("Connecting to fan2go daemon at\n[steelblue]%s[white]...", d.serverAddress)
	statusTextView.SetText(statusText)

	detailsTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	detailsTextView.SetText("[gray]Waiting for initial connection...[white]")

	shortcutsTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	shortcutsTextView.SetText("[yellow]Esc[white]: Inspect UI  •  [yellow]Ctrl+Q[white]: Quit")

	contentLayout.AddItem(statusTextView, 2, 0, false)
	contentLayout.AddItem(nil, 1, 0, false)
	contentLayout.AddItem(detailsTextView, 3, 0, false)
	contentLayout.AddItem(nil, 1, 0, false)
	contentLayout.AddItem(shortcutsTextView, 1, 0, false)

	d.statusTextView = statusTextView
	d.detailsTextView = detailsTextView
	d.layout = createModal(" Connecting ", contentLayout, 56, 11)
}

func (d *ConnectingDialog) SetConnecting() {
	d.lastError = nil
	d.detailsTextView.SetText("[gray]Waiting for initial connection...[white]")
}

func (d *ConnectingDialog) SetError(err error) {
	d.lastError = err
	if err == nil {
		d.SetConnecting()
		return
	}
	d.detailsTextView.SetText(fmt.Sprintf("[red]Connection failed:[white]\n%s\n[yellow]Retrying automatically...[white]", err.Error()))
}

func (d *ConnectingDialog) GetLayout() *tview.Flex {
	return d.layout
}
