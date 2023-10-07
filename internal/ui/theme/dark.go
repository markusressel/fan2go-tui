package theme

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	PrimaryColor        = tcell.ColorDodgerBlue
	PrimaryVariantColor = tcell.ColorSteelBlue
	SecondaryColor      = tcell.ColorGreenYellow

	OnPrimaryColor = tcell.ColorWhite
	OnSecondary    = tcell.ColorBlack

	Colors = Color{
		Header: HeaderColors{
			Name:           OnPrimaryColor,
			NameBackground: PrimaryColor,

			PageIndicator:           OnPrimaryColor,
			PageIndicatorBackground: tcell.ColorBlue,

			UpdateInterval:           OnPrimaryColor,
			UpdateIntervalBackground: PrimaryVariantColor,

			Version:           OnSecondary,
			VersionBackground: SecondaryColor,
		},
		Dialog: DialogColors{
			Border: PrimaryVariantColor,
		},
		Layout: LayoutColors{
			Border: PrimaryVariantColor,
			Title:  tcell.ColorBlue,
		},
		Graphs: GraphsColors{
			Rpm:    tcell.ColorBlue,
			Pwm:    PrimaryVariantColor,
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
	if len(text) <= 0 {
		return ""
	}
	titleText := fmt.Sprintf(" %s ", text)
	return titleText
}
