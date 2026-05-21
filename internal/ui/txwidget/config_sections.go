package txwidget

import (
	"fan2go-tui/internal/client"
	"fmt"
	"maps"
	"sort"
	"strconv"
	"strings"
	"time"
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
				{Label: "Use Unscaled Curve", Value: fmt.Sprintf("%v", config.UseUnscaledCurveValues)},
				{Label: "Min PWM", Value: intPointer(config.MinPwm)},
				{Label: "Start PWM", Value: intPointer(config.StartPwm)},
				{Label: "Max PWM", Value: intPointer(config.MaxPwm)},
				{Label: "PWM Set Delay", Value: durationPointer(config.PwmSetDelay)},
			},
		},
	}

	if config.PwmMap != nil {
		pwmMapFields := pwmMapConfigFields(config.PwmMap)
		sections = append(sections, ConfigInfoSection{
			Title:           "PwmMap",
			Accent:          ConfigInfoAccentMap,
			Fields:          pwmMapFields,
			PreserveTypeRow: true,
		})
	}

	if config.SetPwmToGetPwmMap != nil {
		sections = append(sections, ConfigInfoSection{
			Title:           "SetPwmToGetPwmMap",
			Accent:          ConfigInfoAccentMap,
			Fields:          setPwmToGetPwmMapFields(config.SetPwmToGetPwmMap),
			PreserveTypeRow: true,
		})
	}

	if config.ControlMode != nil {
		sections = append(sections, ConfigInfoSection{
			Title:  "Control Mode",
			Accent: ConfigInfoAccentControlLoop,
			Fields: controlModeFields(config.ControlMode),
		})
	}

	if config.ControlAlgorithm != nil {
		sections = append(sections, ConfigInfoSection{
			Title:  "Control Algorithm",
			Accent: ConfigInfoAccentControlLoop,
			Fields: controlAlgorithmFields(config.ControlAlgorithm),
		})
	}

	if config.SanityCheck != nil {
		sections = append(sections, ConfigInfoSection{
			Title:  "Sanity Check",
			Accent: ConfigInfoAccentControlLoop,
			Fields: sanityCheckFields(config.SanityCheck),
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
	if config.Nvidia != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "Nvidia"},
			{Label: "Device", Value: config.Nvidia.Device},
			{Label: "Index", Value: strconv.Itoa(config.Nvidia.Index)},
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
	if config.Acpi != nil {
		fields := []ConfigInfoField{{Label: "Type", Value: "ACPI"}}
		fields = append(fields, acpiFanCallFields("Set PWM", config.Acpi.SetPwm)...)
		fields = append(fields, acpiFanCallFields("Get PWM", config.Acpi.GetPwm)...)
		fields = append(fields, acpiFanCallFields("Get RPM", config.Acpi.GetRpm)...)
		return fields, "Source"
	}

	return []ConfigInfoField{{Label: "Type", Value: "N/A"}}, "Source"
}

func sensorSource(config client.SensorConfig) ([]ConfigInfoField, string) {
	if config.HwMon != nil {
		fields := []ConfigInfoField{
			{Label: "Type", Value: "HwMon"},
			{Label: "Platform", Value: config.HwMon.Platform},
			{Label: "Index", Value: strconv.Itoa(config.HwMon.Index)},
			{Label: "Temp Input", Value: config.HwMon.TempInput},
		}
		if config.HwMon.Channel > 0 {
			fields = append(fields, ConfigInfoField{Label: "Channel", Value: strconv.Itoa(config.HwMon.Channel)})
		}
		return fields, "Source"
	}
	if config.Nvidia != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "Nvidia"},
			{Label: "Device", Value: config.Nvidia.Device},
			{Label: "Index", Value: strconv.Itoa(config.Nvidia.Index)},
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
	if config.Disk != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "Disk"},
			{Label: "Device", Value: config.Disk.Device},
		}, "Source"
	}
	if config.Acpi != nil {
		conversion := string(config.Acpi.Conversion)
		if conversion == "" {
			conversion = string(client.AcpiSensorConversionCelsius)
		}
		return []ConfigInfoField{
			{Label: "Type", Value: "ACPI"},
			{Label: "Method", Value: config.Acpi.Method},
			{Label: "Args", Value: config.Acpi.Args},
			{Label: "Conversion", Value: conversion},
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
		fields := []ConfigInfoField{
			{Label: "Type", Value: "Function"},
			{Label: "Function", Value: config.Function.Type},
		}
		fields = append(fields, functionCurveFields(config.Function.Curves)...)
		return fields, "Curve"
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

func durationPointer(value *time.Duration) string {
	if value == nil {
		return "N/A"
	}
	return value.String()
}

func pwmMapConfigFields(config *client.PwmMapConfig) []ConfigInfoField {
	if config.Autodetect != nil {
		return []ConfigInfoField{{Label: "Type", Value: "Autodetect"}}
	}
	if config.Identity != nil {
		return []ConfigInfoField{{Label: "Type", Value: "Identity"}}
	}
	if config.Linear != nil {
		return []ConfigInfoField{{Label: "Type", Value: "Linear"}}
	}
	if config.Values != nil {
		return []ConfigInfoField{{Label: "Type", Value: "Values"}}
	}
	return []ConfigInfoField{{Label: "Type", Value: "N/A"}}
}

func setPwmToGetPwmMapFields(config *client.SetPwmToGetPwmMapConfig) []ConfigInfoField {
	if config.Autodetect != nil {
		return []ConfigInfoField{{Label: "Type", Value: "Autodetect"}}
	}
	if config.Identity != nil {
		return []ConfigInfoField{{Label: "Type", Value: "Identity"}}
	}
	if config.Linear != nil {
		return []ConfigInfoField{{Label: "Type", Value: "Linear"}}
	}
	if config.Values != nil {
		return []ConfigInfoField{{Label: "Type", Value: "Values"}}
	}
	return []ConfigInfoField{{Label: "Type", Value: "N/A"}}
}

func controlModeFields(config *client.ControlModeConfig) []ConfigInfoField {
	fields := []ConfigInfoField{}
	if config.Active != nil {
		fields = append(fields, ConfigInfoField{Label: "Active", Value: string(*config.Active)})
	}
	if config.OnExit != nil {
		onExit := config.OnExit
		onExitMode := "Custom"
		if onExit.Restore != nil {
			onExitMode = "Restore"
		} else if onExit.None != nil {
			onExitMode = "None"
		}
		fields = append(fields, ConfigInfoField{Label: "On Exit", Value: onExitMode})
		if onExit.ControlMode != nil {
			fields = append(fields, ConfigInfoField{Label: "Exit Mode", Value: string(*onExit.ControlMode)})
		}
		if onExit.Speed != nil {
			fields = append(fields, ConfigInfoField{Label: "Exit Speed", Value: strconv.Itoa(*onExit.Speed)})
		}
	}
	if len(fields) == 0 {
		return []ConfigInfoField{{Label: "Type", Value: "N/A"}}
	}
	return fields
}

func controlAlgorithmFields(config *client.ControlAlgorithmConfig) []ConfigInfoField {
	if config.Direct != nil {
		fields := []ConfigInfoField{{Label: "Type", Value: "Direct"}}
		fields = append(fields, ConfigInfoField{Label: "Max PWM Change/Cycle", Value: intPointer(config.Direct.MaxPwmChangePerCycle)})
		return fields
	}
	if config.Pid != nil {
		return []ConfigInfoField{
			{Label: "Type", Value: "PID"},
			{Label: "P", Value: floatText(config.Pid.P)},
			{Label: "I", Value: floatText(config.Pid.I)},
			{Label: "D", Value: floatText(config.Pid.D)},
		}
	}
	return []ConfigInfoField{{Label: "Type", Value: "N/A"}}
}

func sanityCheckFields(config *client.SanityCheckConfig) []ConfigInfoField {
	return []ConfigInfoField{
		{Label: "PWM Changed 3rd Party", Value: fmt.Sprintf("%v", bool(config.PwmValueChangedByThirdParty.Enabled))},
		{Label: "Mode Changed 3rd Party", Value: fmt.Sprintf("%v", bool(config.FanModeChangedByThirdParty.Enabled))},
		{Label: "Mode Check Throttle", Value: config.FanModeChangedByThirdParty.ThrottleDuration.String()},
	}
}

func acpiFanCallFields(prefix string, call *client.AcpiFanCallConfig) []ConfigInfoField {
	if call == nil {
		return nil
	}
	conversion := string(call.Conversion)
	if conversion == "" {
		conversion = string(client.AcpiFanConversionPwm)
	}
	return []ConfigInfoField{
		{Label: fmt.Sprintf("%s Method", prefix), Value: call.Method},
		{Label: fmt.Sprintf("%s Args", prefix), Value: call.Args},
		{Label: fmt.Sprintf("%s Conversion", prefix), Value: conversion},
	}
}

func sortedIntMap(labelPrefix string, values map[int]int) []ConfigInfoField {
	keys := make([]int, 0, len(values))
	for key := range maps.Keys(values) {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	fields := make([]ConfigInfoField, 0, len(keys))
	for _, key := range keys {
		fields = append(fields, ConfigInfoField{
			Label: fmt.Sprintf("%s %d", labelPrefix, key),
			Value: strconv.Itoa(values[key]),
		})
	}
	return fields
}

func functionCurveFields(curves []string) []ConfigInfoField {
	if len(curves) == 0 {
		return []ConfigInfoField{{Label: "Curves", Value: "N/A"}}
	}

	fields := make([]ConfigInfoField, 0, len(curves))
	for index, curve := range curves {
		label := ""
		if index == 0 {
			label = "Curves"
		}
		fields = append(fields, ConfigInfoField{
			Label: label,
			Value: fmt.Sprintf("- %s", curve),
		})
	}
	return fields
}
