package theme

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	PrimaryColor        = tcell.ColorDodgerBlue
	PrimaryVariantColor = tcell.ColorSteelBlue
	SecondaryColor      = tcell.ColorGreenYellow

	OnPrimaryColor = tcell.ColorWhite
	OnSecondary    = tcell.ColorBlack
)

var (
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
