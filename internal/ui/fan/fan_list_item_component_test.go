package fan

import (
	"fan2go-tui/internal/client"
	"math"
	"testing"

	"github.com/rivo/tview"
)

func TestFanListItemComponentSetFanWithNilSubcomponentsDoesNotPanic(t *testing.T) {
	c := &FanListItemComponent{}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SetFan panicked with nil subcomponents: %v", r)
		}
	}()

	c.SetFan(nil)
}

func TestFanRpmCurveComponentSetFanOnNilReceiverDoesNotPanic(t *testing.T) {
	var c *FanRpmCurveComponent

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SetFan panicked on nil receiver: %v", r)
		}
	}()

	c.SetFan(nil)
}

func TestHasFanCurveData(t *testing.T) {
	if hasFanCurveData(nil) {
		t.Fatalf("expected no curve data for nil fan")
	}

	empty := map[int]float64{}
	if hasFanCurveData(&client.Fan{FanCurveData: &empty}) {
		t.Fatalf("expected no curve data for empty map")
	}

	withData := map[int]float64{42: 1700}
	if !hasFanCurveData(&client.Fan{FanCurveData: &withData}) {
		t.Fatalf("expected curve data when map is non-empty")
	}
}

func TestFanListItemComponentSwitchesFromHistoryToCurve(t *testing.T) {
	app := tview.NewApplication()

	historyFan := &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Pwm: 30, Rpm: 1200}
	c := NewFanListItemComponent(app, historyFan, nil)

	if c.fanGraphComponent == nil || c.fanRpmCurveComponent != nil {
		t.Fatalf("expected history variant on initialization")
	}

	curveData := map[int]float64{30: 1200}
	curveFan := &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Pwm: 30, Rpm: 1200, FanCurveData: &curveData}
	c.SetFan(curveFan)

	if c.fanGraphComponent != nil || c.fanRpmCurveComponent == nil {
		t.Fatalf("expected curve variant after curve data appears")
	}
}

func TestFanListItemComponentSwitchesFromCurveToHistory(t *testing.T) {
	app := tview.NewApplication()

	curveData := map[int]float64{20: 900}
	curveFan := &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Pwm: 20, Rpm: 900, FanCurveData: &curveData}
	c := NewFanListItemComponent(app, curveFan, nil)

	if c.fanGraphComponent != nil || c.fanRpmCurveComponent == nil {
		t.Fatalf("expected curve variant on initialization")
	}

	historyFan := &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Pwm: 15, Rpm: 800}
	c.SetFan(historyFan)

	if c.fanGraphComponent == nil || c.fanRpmCurveComponent != nil {
		t.Fatalf("expected history variant after curve data disappears")
	}
}

func TestFanListItemComponentUpdatesCurveSeriesWithoutRebuild(t *testing.T) {
	app := tview.NewApplication()

	curveDataA := map[int]float64{10: 1000}
	fanA := &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Pwm: 10, Rpm: 1000, FanCurveData: &curveDataA}
	c := NewFanListItemComponent(app, fanA, nil)

	curveComponent := c.fanRpmCurveComponent
	if curveComponent == nil {
		t.Fatalf("expected curve variant on initialization")
	}

	curveDataB := map[int]float64{20: 2000, 30: 3000}
	fanB := &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Pwm: 20, Rpm: 2000, FanCurveData: &curveDataB}
	c.SetFan(fanB)

	if c.fanRpmCurveComponent != curveComponent {
		t.Fatalf("expected existing curve component to be reused")
	}
	if got := curveComponent.seriesValueProvider.F(20); got != 2000 {
		t.Fatalf("expected updated curve value 2000 for pwm 20, got %v", got)
	}
	if got := curveComponent.seriesValueProvider.F(10); !math.IsNaN(got) {
		t.Fatalf("expected old curve point to be removed, got %v", got)
	}
}
