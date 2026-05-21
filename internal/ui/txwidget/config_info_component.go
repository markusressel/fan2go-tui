package txwidget

import (
	"fan2go-tui/internal/ui/theme"
	"fmt"
	"sort"
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
	Title           string
	Fields          []ConfigInfoField
	Accent          ConfigInfoAccent
	PreserveTypeRow bool
}

type ConfigInfoComponent struct {
	layout *tview.TextView

	isFieldClickable func(sectionTitle, label, value string) bool
	onFieldClick     func(sectionTitle, label, value string)
	fieldClickByID   map[string]clickableConfigField
}

type clickableConfigField struct {
	sectionTitle string
	field        ConfigInfoField
}

func NewConfigInfoComponent() *ConfigInfoComponent {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false).
		SetScrollable(true).
		SetRegions(true)
	textView.SetBorderPadding(0, 0, 0, 0)

	c := &ConfigInfoComponent{
		layout: textView,
		isFieldClickable: func(sectionTitle, label, value string) bool {
			return false
		},
		onFieldClick:   func(sectionTitle, label, value string) {},
		fieldClickByID: map[string]clickableConfigField{},
	}
	textView.SetHighlightedFunc(func(added, _, _ []string) {
		if len(added) == 0 {
			return
		}
		clickableField, ok := c.fieldClickByID[added[0]]
		if !ok {
			return
		}
		c.onFieldClick(clickableField.sectionTitle, clickableField.field.Label, clickableField.field.Value)
		c.layout.Highlight()
	})
	textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			c.ScrollHorizontal(4)
			return nil
		case tcell.KeyLeft:
			c.ScrollHorizontal(-4)
			return nil
		}
		switch event.Rune() {
		case 'l':
			c.ScrollHorizontal(4)
			return nil
		case 'h':
			c.ScrollHorizontal(-4)
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
	w.fieldClickByID = map[string]clickableConfigField{}

	type renderSection struct {
		headerTitle     string
		typeValue       string
		accent          ConfigInfoAccent
		sectionTitle    string
		preserveTypeRow bool
		maxKeyLen       int
		fields          []ConfigInfoField
	}

	renderSections := make([]renderSection, 0, len(sections))
	for _, section := range sections {
		headerTitle, typeValue := sectionHeadline(section)
		fields := make([]ConfigInfoField, 0, len(section.Fields))
		sectionMaxKeyLen := 0
		for _, field := range section.Fields {
			if strings.EqualFold(field.Label, "Type") && !section.PreserveTypeRow {
				continue
			}
			fields = append(fields, field)
			if len(field.Label) > sectionMaxKeyLen {
				sectionMaxKeyLen = len(field.Label)
			}
		}
		renderSections = append(renderSections, renderSection{
			headerTitle:     headerTitle,
			typeValue:       typeValue,
			accent:          section.Accent,
			sectionTitle:    section.Title,
			preserveTypeRow: section.PreserveTypeRow,
			maxKeyLen:       sectionMaxKeyLen,
			fields:          fields,
		})
	}

	priority := func(section renderSection) int {
		if strings.EqualFold(section.headerTitle, "General") || strings.EqualFold(section.sectionTitle, "General") {
			return 0
		}
		if strings.EqualFold(section.sectionTitle, "Source") || strings.EqualFold(section.sectionTitle, "Curve") {
			return 1
		}
		return 2
	}

	sort.SliceStable(renderSections, func(i, j int) bool {
		left := renderSections[i]
		right := renderSections[j]

		leftPriority := priority(left)
		rightPriority := priority(right)
		return leftPriority < rightPriority
	})

	keyColorTag := colorTag(theme.Colors.ConfigInfoComponent.FieldKey)
	var out strings.Builder
	for i, section := range renderSections {
		if i > 0 {
			out.WriteString("\n")
		}
		accentTag := colorTag(resolveAccentColor(section.accent))
		headerText := tview.Escape(fmt.Sprintf("[%s]", section.headerTitle))
		out.WriteString(fmt.Sprintf("%s%s[-]\n", accentTag, headerText))
		for _, field := range section.fields {
			valueColor := resolveValueColor(field.Label, field.Value, section.typeValue)
			valueColorTag := colorTag(valueColor)
			valueText := tview.Escape(field.Value)
			if w.isFieldClickable(section.sectionTitle, field.Label, field.Value) {
				regionID := fmt.Sprintf("field-%d", len(w.fieldClickByID))
				w.fieldClickByID[regionID] = clickableConfigField{
					sectionTitle: section.sectionTitle,
					field:        field,
				}
				valueText = fmt.Sprintf("[\"%s\"]%s[\"\"]", regionID, valueText)
			}
			out.WriteString(fmt.Sprintf("%s%-*s[-] %s%s[-]\n",
				keyColorTag,
				section.maxKeyLen,
				tview.Escape(field.Label),
				valueColorTag,
				valueText,
			))
		}
	}

	w.layout.SetText(out.String())
}

func sectionHeadline(section ConfigInfoSection) (string, string) {
	if section.PreserveTypeRow {
		return section.Title, ""
	}
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

func (w *ConfigInfoComponent) ScrollHorizontal(delta int) {
	rowOffset, colOffset := w.layout.GetScrollOffset()
	nextCol := colOffset + delta
	if nextCol < 0 {
		nextCol = 0
	}
	w.layout.ScrollTo(rowOffset, nextCol)
}

func (w *ConfigInfoComponent) SetFieldClickablePredicate(predicate func(sectionTitle, label, value string) bool) {
	if predicate == nil {
		w.isFieldClickable = func(sectionTitle, label, value string) bool {
			return false
		}
		return
	}
	w.isFieldClickable = predicate
}

func (w *ConfigInfoComponent) SetFieldClickHandler(handler func(sectionTitle, label, value string)) {
	if handler == nil {
		w.onFieldClick = func(sectionTitle, label, value string) {}
		return
	}
	w.onFieldClick = handler
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
		return theme.Colors.ConfigInfoComponent.ValueBool
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
