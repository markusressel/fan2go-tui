package fan

import (
	"fan2go-tui/internal/client"
	"fan2go-tui/internal/ui/theme"
	uiutil "fan2go-tui/internal/ui/util"
	"fmt"
	"github.com/rivo/tview"
)

type CurveComponent struct {
	application *tview.Application

	Curve *client.Curve

	layout *tview.Flex

	configTextView *tview.TextView
	valueTextView  *tview.TextView
}

func NewCurveComponent(application *tview.Application, curve *client.Curve) *CurveComponent {
	c := &CurveComponent{
		application: application,
		Curve:       curve,
	}

	c.layout = c.createLayout()

	return c
}

func (c *CurveComponent) createLayout() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)
	titleText := fmt.Sprintf("Curve: %s", c.Curve.Config.ID)

	layout.SetBorder(true)
	uiutil.SetupWindow(layout, titleText)

	configTextView := tview.NewTextView()
	layout.AddItem(configTextView, 0, 1, false)
	c.configTextView = configTextView

	curveValueTextView := tview.NewTextView().SetTextColor(theme.Colors.Graphs.Curve)
	layout.AddItem(curveValueTextView, 1, 0, false)
	c.valueTextView = curveValueTextView

	return layout
}

func (c *CurveComponent) Refresh() {
	// print basic info
	valueText := fmt.Sprintf("Value: %d", int(c.Curve.Value))
	c.valueTextView.SetText(valueText)

	// print config
	config := c.Curve.Config

	configText := ""
	// configText += fmt.Sprintf("Id: %s\n", config.Id)
	configText += fmt.Sprintf("Curve: %s\n", config.ID)
	// value = strconv.FormatFloat(config.MinPwm, 'f', -1, 64)

	if config.PID != nil {
		configText += fmt.Sprintf("Type: PID\n")
		configText += fmt.Sprintf("  Sensor: %s\n", config.PID.Sensor)
		configText += fmt.Sprintf("  P: %f\n", config.PID.P)
		configText += fmt.Sprintf("  I: %f\n", config.PID.I)
		configText += fmt.Sprintf("  D: %f\n", config.PID.D)
	} else if config.Linear != nil {
		configText += fmt.Sprintf("Type: Linear\n")
		configText += fmt.Sprintf("  Sensor: %s\n", config.Linear.Sensor)
		configText += fmt.Sprintf("  Min: %d\n", config.Linear.Min)
		configText += fmt.Sprintf("  Max: %d\n", config.Linear.Max)
		configText += fmt.Sprintf("  Steps: %v\n", config.Linear.Steps)
	} else if config.Function != nil {
		configText += fmt.Sprintf("Type: Function\n")
		configText += fmt.Sprintf("  Type: %s\n", config.Function.Type)
		configText += fmt.Sprintf("  Curves: %s\n", config.Function.Curves)
	}
	c.configTextView.SetText(configText)
}

func (c *CurveComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurveComponent) SetCurve(curve *client.Curve) {
	c.Curve = curve
}
