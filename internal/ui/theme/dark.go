package theme

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CreateTitleText(text string) string {
	titleText := fmt.Sprintf(" %s ", text)
	return titleText
}

func GetTitleColor() tcell.Color {
	return tcell.ColorBlue
}

func GetDialogTitleAlign() int {
	return tview.AlignCenter
}

func GetTitleAlign() int {
	return tview.AlignLeft
}

func GetDialogBorderColor() tcell.Color {
	return tcell.ColorSteelBlue
}
