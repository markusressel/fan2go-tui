package txwidget

import (
	"fan2go-tui/internal/client"
	"fmt"
	"maps"
	"sort"
	"strconv"
	"strings"
)

func FanConfigSections(config client.FanConfig) []ConfigInfoSection {
	sections := []ConfigInfoSection{
		{
			Title:  "General",
			Accent: ConfigInfoAccentGeneral,
			Fields: []ConfigInfoField{
				{Label: "ID", Value: config.ID},
				{Label: "Curve", Value: config.Curve},
				{Label: "NeverStop", Value: fmt.Sprintf("%v", config.NeverStop)},
				{Label: "Min PWM", Value: intPointer(config.MinPwm)},
				{Label: "Start PWM", Value: intPointer(config.StartPwm)},
				{Label: "Max PWM", Value: intPointer(config.MaxPwm)},
			},
		},
	}

	if config.PwmMap != nil && len(*config.PwmMap) > 0 {
		sections = append(sections, ConfigInfoSection{
			Title:  "PWM Map",
			Accent: ConfigInfoAccentMap,
			Fields: sortedPwmMap(*config.PwmMap),
		})
	}

	if config.ControlLoop != nil {
		sections = append(sections, ConfigInfoSection{
			Title:  "Control Loop",
			Accent: ConfigInfoAccentControlLoop,
			Fields: []ConfigInfoField{
				{Label: "P", Value: floatText(config.ControlLoop.P)},
				{Label: "I", Value: floatText(config.ControlLoop.I)},
				{Label: "D", Value: floatText(config.ControlLoop.D)},
			},
		})
	}

	sourceFields, sourceName := fanSource(config)
	sections = append(sections, ConfigInfoSection{Title: sourceName, Accent: ConfigInfoAccentSource, Fields: sourceFields})

	return sections
}

func SensorConfigSections(config client.SensorConfig) []ConfigInfoSection {
	sourceFields, sourceName := sensorSource(config)
	return []ConfigInfoSection{
		{
			Title:  "General",
			Accent: ConfigInfoAccentGeneral,
			Fields: []ConfigInfoField{
				{Label: "ID", Value: config.ID},
			},
		},
		{
			Title:  sourceName,
			Accent: ConfigInfoAccentSource,
			Fields: sourceFields,
		},
	}
}

func CurveConfigSections(config client.CurveConfig) []ConfigInfoSection {
	sections := []ConfigInfoSection{
		{
			Title:  "General",
			Accent: ConfigInfoAccentGeneral,
			Fields: []ConfigInfoField{
				{Label: "ID", Value: config.ID},
			},
		},
	}

	typeFields, typeName := curveType(config)
	sections = append(sections, ConfigInfoSection{Title: typeName, Accent: ConfigInfoAccentCurve, Fields: typeFields})

	return sections
}

func fanSource(config client.FanConfig) ([]ConfigInfoField, string) {
	if config.File != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "File"},
			{Label: "PWM Path", Value: config.File.Path},
			{Label: "RPM Path", Value: config.File.RpmPath},
		}, "Source"
	}
	if config.HwMon != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "HwMon"},
			{Label: "Platform", Value: config.HwMon.Platform},
			{Label: "Index", Value: strconv.Itoa(config.HwMon.Index)},
			{Label: "PWM Channel", Value: strconv.Itoa(config.HwMon.PwmChannel)},
			{Label: "RPM Channel", Value: strconv.Itoa(config.HwMon.RpmChannel)},
			{Label: "Sysfs Path", Value: config.HwMon.SysfsPath},
			{Label: "PWM Path", Value: config.HwMon.PwmPath},
			{Label: "PWM Enable", Value: config.HwMon.PwmEnablePath},
			{Label: "RPM Input", Value: config.HwMon.RpmInputPath},
		}, "Source"
	}
	if config.Cmd != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "Command"},
			{Label: "Set PWM", Value: execText(config.Cmd.SetPwm)},
			{Label: "Get PWM", Value: execText(config.Cmd.GetPwm)},
			{Label: "Get RPM", Value: execText(config.Cmd.GetRpm)},
		}, "Source"
	}

	return []ConfigInfoField{{Label: "Type", Value: "N/A"}}, "Source"
}

func sensorSource(config client.SensorConfig) ([]ConfigInfoField, string) {
	if config.HwMon != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "HwMon"},
			{Label: "Platform", Value: config.HwMon.Platform},
			{Label: "Index", Value: strconv.Itoa(config.HwMon.Index)},
			{Label: "Temp Input", Value: config.HwMon.TempInput},
		}, "Source"
	}
	if config.File != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "File"},
			{Label: "Path", Value: config.File.Path},
		}, "Source"
	}
	if config.Cmd != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "Command"},
			{Label: "Exec", Value: config.Cmd.Exec},
			{Label: "Args", Value: strings.Join(config.Cmd.Args, " ")},
		}, "Source"
	}

	return []ConfigInfoField{{Label: "Type", Value: "N/A"}}, "Source"
}

func curveType(config client.CurveConfig) ([]ConfigInfoField, string) {
	if config.PID != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "PID"},
			{Label: "Sensor", Value: config.PID.Sensor},
			{Label: "Set Point", Value: floatText(config.PID.SetPoint)},
			{Label: "P", Value: floatText(config.PID.P)},
			{Label: "I", Value: floatText(config.PID.I)},
			{Label: "D", Value: floatText(config.PID.D)},
		}, "Curve"
	}
	if config.Linear != nil {
		fields := []ConfigInfoField{
			{Label: "Type", Value: "Linear"},
			{Label: "Sensor", Value: config.Linear.Sensor},
		}
		if len(config.Linear.Steps) > 0 {
			fields = append(fields, sortedCurveSteps(config.Linear.Steps)...)
		} else {
			fields = append(fields,
				ConfigInfoField{Label: "Min", Value: strconv.Itoa(config.Linear.Min)},
				ConfigInfoField{Label: "Max", Value: strconv.Itoa(config.Linear.Max)},
			)
		}
		return fields, "Curve"
	}
	if config.Function != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "Function"},
			{Label: "Function", Value: config.Function.Type},
			{Label: "Curves", Value: strings.Join(config.Function.Curves, ", ")},
		}, "Curve"
	}

	return []ConfigInfoField{{Label: "Type", Value: "N/A"}}, "Curve"
}

func sortedPwmMap(values map[int]int) []ConfigInfoField {
	keys := make([]int, 0, len(values))
	for key := range maps.Keys(values) {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	fields := make([]ConfigInfoField, 0, len(keys))
	for _, key := range keys {
		fields = append(fields, ConfigInfoField{
			Label: fmt.Sprintf("Temp %d", key),
			Value: strconv.Itoa(values[key]),
		})
	}
	return fields
}

func sortedCurveSteps(values map[int]float64) []ConfigInfoField {
	keys := make([]int, 0, len(values))
	for key := range maps.Keys(values) {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	fields := make([]ConfigInfoField, 0, len(keys))
	for _, key := range keys {
		fields = append(fields, ConfigInfoField{
			Label: fmt.Sprintf("Step %d", key),
			Value: floatText(values[key]),
		})
	}
	return fields
}

func intPointer(value *int) string {
	if value == nil {
		return "N/A"
	}
	return strconv.Itoa(*value)
}

func execText(config *client.ExecConfig) string {
	if config == nil {
		return "N/A"
	}
	if len(config.Args) == 0 {
		return config.Exec
	}
	return fmt.Sprintf("%s %s", config.Exec, strings.Join(config.Args, " "))
}

func floatText(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
