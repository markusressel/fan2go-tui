package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
	"fmt"
	"github.com/rivo/tview"
)

type FanInfoComponent struct {
	application *tview.Application

	Fan *client.Fan

	layout *tview.Flex

	configTextView   *tview.TextView
	pwmValueTextView *tview.TextView
	rpmValueTextView *tview.TextView
}

func NewFanInfoComponent(application *tview.Application, fan *client.Fan) *FanInfoComponent {
	c := &FanInfoComponent{
		application: application,
		Fan:         fan,
	}

	c.layout = c.createLayout()

	return c
}

func (c *FanInfoComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	configTextView := tview.NewTextView()
	layout.AddItem(configTextView, 0, 1, false)
	c.configTextView = configTextView

	rpmValueTextView := tview.NewTextView().SetTextColor(theme.Colors.Graph.Rpm)
	layout.AddItem(rpmValueTextView, 1, 0, false)
	c.rpmValueTextView = rpmValueTextView

	pwmValueTextView := tview.NewTextView().SetTextColor(theme.Colors.Graph.Pwm)
	layout.AddItem(pwmValueTextView, 1, 0, false)
	c.pwmValueTextView = pwmValueTextView

	return layout
}

func (c *FanInfoComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *FanInfoComponent) SetFan(fan *client.Fan) {
	c.Fan = fan
	c.refresh()
}

func (c *FanInfoComponent) refresh() {
	// print basic info
	pwmText := fmt.Sprintf("PWM: %d", c.Fan.Pwm)
	c.pwmValueTextView.SetText(pwmText)

	rpmText := fmt.Sprintf("RPM: %d", c.Fan.Rpm)
	c.rpmValueTextView.SetText(rpmText)

	// print config
	config := c.Fan.Config

	configText := ""
	// configText += fmt.Sprintf("ID: %s\n", config.ID)
	configText += fmt.Sprintf("Curve: %s\n", config.Curve)
	configText += fmt.Sprintf("Pwm:\n")
	minPwmText := "N/A"
	if config.MinPwm != nil {
		minPwmText = fmt.Sprintf("%d", *config.MinPwm)
	}
	configText += fmt.Sprintf("  Min: %s\n", minPwmText)
	startPwmText := "N/A"
	if config.StartPwm != nil {
		startPwmText = fmt.Sprintf("%d", *config.StartPwm)
	}
	configText += fmt.Sprintf("  Start: %s\n", startPwmText)
	maxPwmText := "N/A"
	if config.MaxPwm != nil {
		maxPwmText = fmt.Sprintf("%d", *config.MaxPwm)
	}
	configText += fmt.Sprintf("  Max: %s\n", maxPwmText)
	configText += fmt.Sprintf("NeverStop: %v\n", config.NeverStop)

	// value = strconv.FormatFloat(config.MinPwm, 'f', -1, 64)

	if config.File != nil {
		configText += fmt.Sprintf("Type: File\n")
		configText += fmt.Sprintf("  PwmPath: %s\n", config.File.Path)
		configText += fmt.Sprintf("  RpmPath: %s\n", config.File.RpmPath)
	} else if config.HwMon != nil {
		configText += fmt.Sprintf("Type: HwMon\n")
		configText += fmt.Sprintf("  Platform: %s\n", config.HwMon.Platform)
		configText += fmt.Sprintf("  Index: %d\n", config.HwMon.Index)
		configText += fmt.Sprintf("  PwmChannel: %d\n", config.HwMon.PwmChannel)
		configText += fmt.Sprintf("  RpmChannel: %d\n", config.HwMon.RpmChannel)
		configText += fmt.Sprintf("  SysfsPath: %s\n", config.HwMon.SysfsPath)
		configText += fmt.Sprintf("  PwmPath: %s\n", config.HwMon.PwmPath)
		configText += fmt.Sprintf("  PwmEnablePath: %s\n", config.HwMon.PwmEnablePath)
		configText += fmt.Sprintf("  RpmInputPath: %s\n", config.HwMon.RpmInputPath)
	} else if config.Cmd != nil {
		configText += fmt.Sprintf("Type: Cmd\n")
		configText += fmt.Sprintf("  GetPwm: %s\n", config.Cmd.GetPwm)
		configText += fmt.Sprintf("  SetPwm: %s\n", config.Cmd.SetPwm)
		configText += fmt.Sprintf("  GetRpm: %s\n", config.Cmd.GetRpm)
	}
	c.configTextView.SetText(configText)
}
