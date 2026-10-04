package demo

import (
	"testing"
)

func TestNewSimulator(t *testing.T) {
	sim := NewSimulator()

	fans := sim.GetFans()
	if len(fans) != 4 {
		t.Fatalf("expected 4 fans, got %d", len(fans))
	}

	curves := sim.GetCurves()
	if len(curves) != 4 {
		t.Fatalf("expected 4 curves, got %d", len(curves))
	}

	sensors := sim.GetSensors()
	if len(sensors) != 4 {
		t.Fatalf("expected 4 sensors, got %d", len(sensors))
	}

	// Verify links between fans and curves
	for _, f := range fans {
		curveID := f.Config.Curve
		if curveID == "" {
			t.Fatalf("fan %s has no curve configured", f.Config.ID)
		}
		if _, ok := curves[curveID]; !ok {
			t.Fatalf("fan %s references non-existent curve %s", f.Config.ID, curveID)
		}
	}

	// Verify linear/pid curves reference valid sensors
	for _, c := range curves {
		if c.Config.Linear != nil {
			sensorID := c.Config.Linear.Sensor
			if _, ok := sensors[sensorID]; !ok {
				t.Fatalf("curve %s references non-existent sensor %s", c.Config.ID, sensorID)
			}
		}
		if c.Config.PID != nil {
			sensorID := c.Config.PID.Sensor
			if _, ok := sensors[sensorID]; !ok {
				t.Fatalf("curve %s references non-existent sensor %s", c.Config.ID, sensorID)
			}
		}
		if c.Config.Function != nil {
			for _, refCurve := range c.Config.Function.Curves {
				if _, ok := curves[refCurve]; !ok {
					t.Fatalf("function curve %s references non-existent curve %s", c.Config.ID, refCurve)
				}
			}
		}
	}

	// Verify that we have fans with and without fan curve data
	hasCurveDataCount := 0
	for _, f := range fans {
		if f.FanCurveData != nil && len(*f.FanCurveData) > 0 {
			hasCurveDataCount++
		}
	}
	if hasCurveDataCount == 0 {
		t.Fatalf("expected at least one fan with FanCurveData")
	}
	if hasCurveDataCount == len(fans) {
		t.Fatalf("expected at least one fan without FanCurveData to test both graph variants")
	}
}

func TestSimulator_Tick(t *testing.T) {
	sim := NewSimulator()

	initialCpuFan, _ := sim.GetFan("cpu-fan")
	initialCpuTemp, _ := sim.GetSensor("cpu-temp")

	for i := 0; i < 20; i++ {
		sim.Tick(0.25)
	}

	fans := sim.GetFans()
	for id, f := range fans {
		if f.Pwm < 0 || f.Pwm > 255 {
			t.Fatalf("fan %s PWM %d out of range [0, 255]", id, f.Pwm)
		}
		if f.Rpm < 0 || f.Rpm > 6000 {
			t.Fatalf("fan %s RPM %d out of reasonable range", id, f.Rpm)
		}
	}

	curves := sim.GetCurves()
	for id, c := range curves {
		if c.Value < 0 || c.Value > 255 {
			t.Fatalf("curve %s value %f out of range [0, 255]", id, c.Value)
		}
	}

	sensors := sim.GetSensors()
	for id, s := range sensors {
		tempC := s.MovingAvg / 1000.0
		if tempC < 10.0 || tempC > 110.0 {
			t.Fatalf("sensor %s temp %f C out of reasonable range", id, tempC)
		}
	}

	currentCpuFan, _ := sim.GetFan("cpu-fan")
	currentCpuTemp, _ := sim.GetSensor("cpu-temp")

	// Dynamic values should have evolved over 20 ticks
	if currentCpuTemp.MovingAvg == initialCpuTemp.MovingAvg && currentCpuFan.Pwm == initialCpuFan.Pwm {
		t.Logf("Note: values identical after ticks (could be coincidental, but expected to vary)")
	}
}

func TestSimulator_Getters(t *testing.T) {
	sim := NewSimulator()

	if _, ok := sim.GetFan("cpu-fan"); !ok {
		t.Fatalf("expected cpu-fan to exist")
	}
	if _, ok := sim.GetFan("unknown-fan"); ok {
		t.Fatalf("expected unknown-fan to not exist")
	}

	if _, ok := sim.GetCurve("cpu-linear-curve"); !ok {
		t.Fatalf("expected cpu-linear-curve to exist")
	}
	if _, ok := sim.GetCurve("unknown-curve"); ok {
		t.Fatalf("expected unknown-curve to not exist")
	}

	if _, ok := sim.GetSensor("cpu-temp"); !ok {
		t.Fatalf("expected cpu-temp to exist")
	}
	if _, ok := sim.GetSensor("unknown-sensor"); ok {
		t.Fatalf("expected unknown-sensor to not exist")
	}
}
