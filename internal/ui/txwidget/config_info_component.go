package txwidget

import (
	"fan2go-tui/internal/ui/theme"
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ConfigInfoAccent string

const (
	ConfigInfoAccentDefault     ConfigInfoAccent = "default"
	ConfigInfoAccentGeneral     ConfigInfoAccent = "general"
	ConfigInfoAccentSource      ConfigInfoAccent = "source"
	ConfigInfoAccentCurve       ConfigInfoAccent = "curve"
	ConfigInfoAccentMap         ConfigInfoAccent = "map"
	ConfigInfoAccentControlLoop ConfigInfoAccent = "controlLoop"
)

type ConfigInfoField struct {
	Label string
	Value string
}

type ConfigInfoSection struct {
	Title  string
	Fields []ConfigInfoField
	Accent ConfigInfoAccent
}

type ConfigInfoComponent struct {
	layout *tview.TextView
}

func NewConfigInfoComponent() *ConfigInfoComponent {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetScrollable(true)
	textView.SetBorderPadding(0, 0, 0, 0)

	c := &ConfigInfoComponent{layout: textView}
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			c.scrollHorizontal(4)
			return nil
		case tcell.KeyLeft:
			c.scrollHorizontal(-4)
			return nil
		}
		switch event.Rune() {
		case 'l':
			c.scrollHorizontal(4)
			return nil
		case 'h':
			c.scrollHorizontal(-4)
			return nil
		}
		return event
	})

	return c
}

func (w *ConfigInfoComponent) GetPrimitive() tview.Primitive {
	return w.layout
}

func (w *ConfigInfoComponent) SetSections(sections []ConfigInfoSection) {
	type renderSection struct {
		headerTitle string
		typeValue   string
		accent      ConfigInfoAccent
		fields      []ConfigInfoField
	}

	renderSections := make([]renderSection, 0, len(sections))
	maxKeyLen := 0
	for _, section := range sections {
		headerTitle, typeValue := sectionHeadline(section)
		fields := make([]ConfigInfoField, 0, len(section.Fields))
		for _, field := range section.Fields {
			if strings.EqualFold(field.Label, "Type") {
				continue
			}
			fields = append(fields, field)
			if len(field.Label) > maxKeyLen {
				maxKeyLen = len(field.Label)
			}
		}
		renderSections = append(renderSections, renderSection{
			headerTitle: headerTitle,
			typeValue:   typeValue,
			accent:      section.Accent,
			fields:      fields,
		})
	}

	keyColorTag := colorTag(theme.Colors.ConfigInfoComponent.FieldKey)
	var out strings.Builder
	for i, section := range renderSections {
		if i > 0 {
			out.WriteString("\n")
		}
		accentTag := colorTag(resolveAccentColor(section.accent))
		out.WriteString(fmt.Sprintf("%s[%s][-]\n", accentTag, tview.Escape(section.headerTitle)))
		for _, field := range section.fields {
			valueColor := resolveValueColor(field.Label, field.Value, section.typeValue)
			valueColorTag := colorTag(valueColor)
			out.WriteString(fmt.Sprintf("%s%-*s[-] %s%s[-]\n",
				keyColorTag,
				maxKeyLen,
				tview.Escape(field.Label),
				valueColorTag,
				tview.Escape(field.Value),
			))
		}
	}

	w.layout.SetText(out.String())
}

func sectionHeadline(section ConfigInfoSection) (string, string) {
	for _, field := range section.Fields {
		if strings.EqualFold(field.Label, "Type") {
			if field.Value == "" || strings.EqualFold(field.Value, "N/A") {
				return section.Title, field.Value
			}
			return field.Value, field.Value
		}
	}
	return section.Title, ""
}

func (w *ConfigInfoComponent) scrollHorizontal(delta int) {
	rowOffset, colOffset := w.layout.GetScrollOffset()
	nextCol := colOffset + delta
	if nextCol < 0 {
		nextCol = 0
	}
	w.layout.ScrollTo(rowOffset, nextCol)
}

func colorTag(color tcell.Color) string {
	r, g, b := color.RGB()
	return fmt.Sprintf("[#%02x%02x%02x]", uint8(r), uint8(g), uint8(b))
}

func resolveValueColor(label, value, typeValue string) tcell.Color {
	if value == "" || strings.EqualFold(value, "N/A") {
		return theme.Colors.ConfigInfoComponent.ValueSpecial
	}
	if typeValue != "" && strings.EqualFold(value, typeValue) {
		return theme.Colors.ConfigInfoComponent.SectionDefault
	}
	if strings.HasPrefix(value, "/") || strings.Contains(strings.ToLower(label), "path") {
		return theme.Colors.ConfigInfoComponent.ValuePath
	}
	if _, err := strconv.ParseFloat(value, 64); err == nil {
		return theme.Colors.ConfigInfoComponent.ValueNumber
	}
	if strings.EqualFold(value, "true") || strings.EqualFold(value, "false") {
		return theme.Colors.ConfigInfoComponent.ValueText
	}
	return theme.Colors.ConfigInfoComponent.ValueText
}

func resolveAccentColor(accent ConfigInfoAccent) tcell.Color {
	switch accent {
	case ConfigInfoAccentGeneral:
		return theme.Colors.ConfigInfoComponent.SectionGeneral
	case ConfigInfoAccentSource:
		return theme.Colors.ConfigInfoComponent.SectionSource
	case ConfigInfoAccentCurve:
		return theme.Colors.ConfigInfoComponent.SectionCurve
	case ConfigInfoAccentMap:
		return theme.Colors.ConfigInfoComponent.SectionMap
	case ConfigInfoAccentControlLoop:
		return theme.Colors.ConfigInfoComponent.SectionControlLoop
	default:
		if theme.Colors.ConfigInfoComponent.SectionDefault != tcell.ColorDefault {
			return theme.Colors.ConfigInfoComponent.SectionDefault
		}
		return theme.Colors.Layout.Title
	}
}
