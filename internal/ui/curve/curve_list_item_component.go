package curve

import (
	"fan2go-tui/internal/state"
	"fan2go-tui/internal/ui/txwidgets"
	uiutil "fan2go-tui/internal/ui/util"

	"github.com/rivo/tview"
)

type CurveListItemComponent struct {
	application *tview.Application

	CurveState *state.CurveState

	layout *tview.Flex
	// window is the bordered box of the curve, its title shows whether the curve has warnings
	window *tview.Flex

	curveInfoComponent  *CurveInfoComponent
	curveGraphComponent *CurveGraphComponent
}

func NewCurveListItemComponent(application *tview.Application, curveState *state.CurveState, onOpenSensor func(sensorID string), onOpenCurve func(curveID string)) *CurveListItemComponent {
	c := &CurveListItemComponent{
		application: application,
		CurveState:  curveState,
	}

	c.layout = c.createLayout(onOpenSensor, onOpenCurve)

	return c
}

func (c *CurveListItemComponent) createLayout(onOpenSensor func(sensorID string), onOpenCurve func(curveID string)) *tview.Flex {
	rootLayout := tview.NewFlex().SetDirection(tview.FlexRow)

	curveColumnLayout := tview.NewFlex().SetDirection(tview.FlexColumn)
	uiutil.SetupWindow(curveColumnLayout, c.CurveState.Curve.Config.ID)
	c.window = curveColumnLayout
	c.refreshWindowTitle()
	curveColumnLayout.SetTitleAlign(tview.AlignLeft)
	curveColumnLayout.SetBorder(true)
	rootLayout.AddItem(curveColumnLayout, 0, 1, true)

	f := c.CurveState.Curve
	curveInfoComponent := NewCurveInfoComponent(c.application, f, onOpenSensor, onOpenCurve)
	c.curveInfoComponent = curveInfoComponent
	curveInfoComponent.SetCurve(f)
	layout := curveInfoComponent.GetLayout()
	curveColumnLayout.AddItem(layout, 0, 1, true)
	curveColumnLayout.AddItem(tview.NewBox(), 1, 0, false)

	curveGraphComponent := NewCurveGraphComponent(c.application, c.CurveState)
	c.curveGraphComponent = curveGraphComponent
	curveGraphComponent.SetCurve(c.CurveState)
	layout = curveGraphComponent.GetLayout()
	curveColumnLayout.AddItem(layout, 0, 3, true)

	return rootLayout
}

func (c *CurveListItemComponent) GetLayout() *tview.Flex {
	return c.layout
}

func (c *CurveListItemComponent) SetCurve(curveState *state.CurveState) {
	c.CurveState = curveState
	c.refresh()
}

func (c *CurveListItemComponent) refresh() {
	c.refreshWindowTitle()
	c.curveInfoComponent.SetCurve(c.CurveState.Curve)
	c.curveGraphComponent.SetCurve(c.CurveState)
}

func (c *CurveListItemComponent) ScrollHorizontal(delta int) {
	if c.curveInfoComponent != nil {
		c.curveInfoComponent.ScrollHorizontal(delta)
	}
}

// refreshWindowTitle marks the title of the window if the curve has warnings.
func (c *CurveListItemComponent) refreshWindowTitle() {
	if c.window == nil || c.CurveState == nil || c.CurveState.Curve == nil {
		return
	}
	txwidgets.SetWindowWarnings(c.window.Box, c.CurveState.Curve.Config.ID, c.CurveState.Curve.Warnings)
}
