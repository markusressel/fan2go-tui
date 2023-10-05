package theme

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	Colors = Color{
		Header: HeaderColors{
			Name:           tcell.ColorWhite,
			NameBackground: tcell.ColorDodgerBlue,

			PageIndicator:           tcell.ColorWhite,
			PageIndicatorBackground: tcell.ColorBlue,

			UpdateInterval:           tcell.ColorWhite,
			UpdateIntervalBackground: tcell.ColorSteelBlue,

			Version:           tcell.ColorBlack,
			VersionBackground: tcell.ColorGreenYellow,
		},
		Dialog: DialogColors{
			Border: tcell.ColorSteelBlue,
		},
		Layout: LayoutColors{
			Border: tcell.ColorSteelBlue,
			Title:  tcell.ColorBlue,
		},
		Graphs: GraphsColors{
			Rpm:    tcell.ColorBlue,
			Pwm:    tcell.ColorSteelBlue,
			Curve:  tcell.ColorBlue,
			Sensor: tcell.ColorBlue,
		},
	}

	Style = StyleStruct{
		Layout: LayoutStyle{
			TitleAlign:       tview.AlignCenter,
			DialogTitleAlign: tview.AlignCenter,
		},
	}
)

type HeaderColors struct {
	Name           tcell.Color
	NameBackground tcell.Color

	PageIndicator           tcell.Color
	PageIndicatorBackground tcell.Color

	UpdateInterval           tcell.Color
	UpdateIntervalBackground tcell.Color

	Version           tcell.Color
	VersionBackground tcell.Color
}

type DialogColors struct {
	Border tcell.Color
}

type StyleStruct struct {
	Layout LayoutStyle
}

type LayoutStyle struct {
	TitleAlign       int
	DialogTitleAlign int
}

type Color struct {
	Header HeaderColors
	Dialog DialogColors
	Layout LayoutColors
	Graphs GraphsColors
}

type GraphsColors struct {
	Rpm    tcell.Color
	Pwm    tcell.Color
	Curve  tcell.Color
	Sensor tcell.Color
}

type LayoutColors struct {
	Border tcell.Color
	Title  tcell.Color
}

func CreateTitleText(text string) string {
	titleText := fmt.Sprintf(" %s ", text)
	return titleText
}
