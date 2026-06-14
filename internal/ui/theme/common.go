package theme

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
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
	Header              HeaderColors
	Dialog              DialogColors
	Layout              LayoutColors
	Graph               GraphsColors
	List                ListsColors
	ConfigInfoComponent ConfigInfoComponentColors
	ShortcutMap         ShortcutMapColors
}

type ShortcutMapColors struct {
	KeyCombo tcell.Color
	Name     tcell.Color
}

type ConfigInfoComponentColors struct {
	SectionDefault     tcell.Color
	SectionGeneral     tcell.Color
	SectionSource      tcell.Color
	SectionCurve       tcell.Color
	SectionMap         tcell.Color
	SectionControlLoop tcell.Color
	FieldKey           tcell.Color
	ValueText          tcell.Color
	ValueNumber        tcell.Color
	ValueBool          tcell.Color
	ValuePath          tcell.Color
	ValueSpecial       tcell.Color
}

type GraphsColors struct {
	Default tcell.Color

	XAxisValueLabelText       tcell.Color
	XAxisValueLabelBackground tcell.Color

	YAxisValueLabelText       tcell.Color
	YAxisValueLabelBackground tcell.Color

	// Fan
	Rpm              tcell.Color
	Pwm              tcell.Color
	CurrentPwmLine   tcell.Color
	CurrentRpmMarker tcell.Color
	HeatmapBase      tcell.Color

	// Curve
	Curve    tcell.Color
	CurveMin tcell.Color
	CurveMax tcell.Color

	// Sensor
	Sensor    tcell.Color
	SensorMin tcell.Color
	SensorMax tcell.Color
}

type LayoutColors struct {
	Border tcell.Color
	Title  tcell.Color
}

type ListsColors struct {
	Scrollbar ScrollbarColors
}

type ScrollbarColors struct {
	Bar               tcell.Color
	IndicatorInactive tcell.Color
	IndicatorActive   tcell.Color
	Background        tcell.Color
}

func CreateTitleText(text string) string {
	if len(text) <= 0 {
		return ""
	}
	titleText := fmt.Sprintf(" %s ", text)
	return titleText
}
