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
		Graph: GraphsColors{
			Default: tcell.ColorBlue,

			XAxisValueLabelText:       SecondaryColor,
			XAxisValueLabelBackground: OnSecondary,

			YAxisValueLabelText:       SecondaryColor,
			YAxisValueLabelBackground: OnSecondary,

			Rpm:              tcell.ColorBlue,
			Pwm:              PrimaryVariantColor,
			CurrentPwmLine:   PrimaryVariantColor,
			CurrentRpmMarker: SecondaryColor,
			HeatmapBase:      tcell.ColorGold,

			Curve:    tcell.ColorBlue,
			CurveMin: tcell.ColorGray,
			CurveMax: tcell.ColorGray,

			Sensor:    tcell.ColorBlue,
			SensorMin: tcell.ColorGray,
			SensorMax: tcell.ColorRed,
		},
		List: ListsColors{
			Scrollbar: ScrollbarColors{
				Bar:               tcell.ColorBlue,
				IndicatorActive:   tcell.ColorBlue,
				IndicatorInactive: tcell.ColorGray,
				Background:        tcell.ColorBlack,
			},
		},
		ConfigInfoComponent: ConfigInfoComponentColors{
			SectionDefault:     tcell.ColorSteelBlue,
			SectionGeneral:     tcell.ColorSkyblue,
			SectionSource:      tcell.ColorMediumTurquoise,
			SectionCurve:       tcell.ColorMediumPurple,
			SectionMap:         tcell.ColorKhaki,
			SectionControlLoop: tcell.ColorLightSeaGreen,
			FieldKey:           tcell.ColorSilver,
			ValueText:          tcell.ColorWhite,
			ValueNumber:        tcell.ColorLightSkyBlue,
			ValueBool:          tcell.ColorMediumPurple,
			ValuePath:          tcell.ColorTan,
			ValueSpecial:       tcell.ColorDarkGray,
		},
	}

	Style = StyleStruct{
		Layout: LayoutStyle{
			TitleAlign:       tview.AlignCenter,
			DialogTitleAlign: tview.AlignCenter,
		},
	}
)
