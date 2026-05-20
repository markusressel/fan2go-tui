package fan

import (
	"fan2go-tui/internal/client"
	"testing"
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
