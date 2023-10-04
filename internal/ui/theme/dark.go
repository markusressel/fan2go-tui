package theme

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	Colors = Color{
		Layout: LayoutColors{
			Border: tcell.ColorSteelBlue,
			Title:  tcell.ColorBlue,
		},
		Graphs: GraphsColors{
			First:  tcell.ColorRed,
			Second: tcell.ColorGreen,
		},
	}
)

type Color struct {
	Layout LayoutColors
	Graphs GraphsColors
}

type GraphsColors struct {
	First  tcell.Color
	Second tcell.Color
}

type LayoutColors struct {
	Border tcell.Color
	Title  tcell.Color
}

func CreateTitleText(text string) string {
	titleText := fmt.Sprintf(" %s ", text)
	return titleText
}

func GetDialogTitleAlign() int {
	return tview.AlignCenter
}

func GetTitleAlign() int {
	return tview.AlignLeft
}
