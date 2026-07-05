package state

import (
	"fan2go-tui/internal/client"
	"testing"
)

func TestStore_Fans(t *testing.T) {
	store := NewStore()

	fanMap := map[string]*client.Fan{
		"fan-1": {Config: client.FanConfig{ID: "fan-1"}, Rpm: 1000, Pwm: 100},
	}
	store.UpdateFans(fanMap)

	fans := store.GetFans()
	if len(fans) != 1 || fans["fan-1"].Rpm != 1000 {
		t.Fatalf("expected fan-1 with rpm 1000")
	}

	fanState := store.GetFanState("fan-1")
	if len(fanState.RpmValues) != 1 || fanState.RpmValues[0] != 1000 {
		t.Fatalf("expected fan-1 RpmValues to have 1000")
	}
	if len(fanState.PwmValues) != 1 || fanState.PwmValues[0] != 100 {
		t.Fatalf("expected fan-1 PwmValues to have 100")
	}

	// Update again
	fanMap["fan-1"] = &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Rpm: 1500, Pwm: 150}
	store.UpdateFans(fanMap)

	fanState = store.GetFanState("fan-1")
	if len(fanState.RpmValues) != 2 || fanState.RpmValues[1] != 1500 {
		t.Fatalf("expected fan-1 RpmValues to append 1500")
	}

	// Test max history truncation
	store.maxHistory = 2
	fanMap["fan-1"] = &client.Fan{Config: client.FanConfig{ID: "fan-1"}, Rpm: 2000, Pwm: 200}
	store.UpdateFans(fanMap)

	fanState = store.GetFanState("fan-1")
	if len(fanState.RpmValues) != 2 {
		t.Fatalf("expected fan-1 RpmValues to be truncated to maxHistory 2")
	}
	if fanState.RpmValues[0] != 1500 || fanState.RpmValues[1] != 2000 {
		t.Fatalf("expected truncated RpmValues to be 1500, 2000")
	}

	// Test removal
	store.UpdateFans(map[string]*client.Fan{})
	if store.GetFanState("fan-1") != nil {
		t.Fatalf("expected fan-1 to be removed")
	}
}

func TestStore_Curves(t *testing.T) {
	store := NewStore()

	curveMap := map[string]*client.Curve{
		"curve-1": {Config: client.CurveConfig{ID: "curve-1"}, Value: 120},
	}
	store.UpdateCurves(curveMap)

	curves := store.GetCurves()
	if len(curves) != 1 || curves["curve-1"].Value != 120 {
		t.Fatalf("expected curve-1 with value 120")
	}

	curveState := store.GetCurveState("curve-1")
	if len(curveState.Values) != 1 || curveState.Values[0] != 120 {
		t.Fatalf("expected curve-1 Values to have 120")
	}

	store.UpdateCurves(map[string]*client.Curve{})
	if store.GetCurveState("curve-1") != nil {
		t.Fatalf("expected curve-1 to be removed")
	}
}

func TestStore_Sensors(t *testing.T) {
	store := NewStore()

	sensorMap := map[string]*client.Sensor{
		"sensor-1": {Config: client.SensorConfig{ID: "sensor-1"}, MovingAvg: 45000},
	}
	store.UpdateSensors(sensorMap)

	sensors := store.GetSensors()
	if len(sensors) != 1 || sensors["sensor-1"].MovingAvg != 45000 {
		t.Fatalf("expected sensor-1 with moving avg 45000")
	}

	sensorState := store.GetSensorState("sensor-1")
	if len(sensorState.Values) != 1 || sensorState.Values[0] != 45.0 {
		t.Fatalf("expected sensor-1 Values to have 45.0")
	}

	store.UpdateSensors(map[string]*client.Sensor{})
	if store.GetSensorState("sensor-1") != nil {
		t.Fatalf("expected sensor-1 to be removed")
	}
}
