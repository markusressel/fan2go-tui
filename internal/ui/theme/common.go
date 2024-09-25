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
	Header HeaderColors
	Dialog DialogColors
	Layout LayoutColors
	Graph  GraphsColors
	List   ListsColors
}

type GraphsColors struct {
	Default tcell.Color

	// Fan
	Rpm tcell.Color
	Pwm tcell.Color

	// Curve
	Curve    tcell.Color
	CurveMin tcell.Color
	CurveMax tcell.Color

	// Sensor
	Sensor tcell.Color
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
