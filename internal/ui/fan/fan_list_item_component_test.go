package fan

import "testing"

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
