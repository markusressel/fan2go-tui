package util

import (
	"fan2go-tui/internal/ui/theme"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Window[T any] interface {
	SetTitle(title string) *T
	SetTitleColor(color tcell.Color) *T
	SetTitleAlign(align int) *T
	SetBorderColor(color tcell.Color) *T
}

func SetupWindow[T Window[tview.Box]](window T, text string) T {
	window.SetTitle(theme.CreateTitleText(text))
	window.SetTitleColor(theme.Colors.Layout.Title)
	window.SetTitleAlign(theme.Style.Layout.TitleAlign)
	window.SetBorderColor(theme.Colors.Layout.Border)
	return window
}

func SetupDialogWindow[T Window[tview.Box]](window T, text string) T {
	window.SetTitle(theme.CreateTitleText(text))
	window.SetTitleColor(theme.Colors.Layout.Title)
	window.SetTitleAlign(theme.Style.Layout.DialogTitleAlign)
	window.SetBorderColor(theme.Colors.Dialog.Border)
	return window
}
