package state

import (
	"errors"
	"fan2go-tui/internal/client"
	"testing"
)

type mockApiClient struct {
	fansErr    error
	curvesErr  error
	sensorsErr error

	fans    map[string]*client.Fan
	curves  map[string]*client.Curve
	sensors map[string]*client.Sensor
}

func (m *mockApiClient) GetFans() (*map[string]*client.Fan, error) {
	if m.fansErr != nil {
		return nil, m.fansErr
	}
	return &m.fans, nil
}

func (m *mockApiClient) GetFan(label string) (*client.Fan, error) {
	if m.fansErr != nil {
		return nil, m.fansErr
	}
	f, ok := m.fans[label]
	if !ok {
		return nil, errors.New("not found")
	}
	return f, nil
}

func (m *mockApiClient) GetCurves() (*map[string]*client.Curve, error) {
	if m.curvesErr != nil {
		return nil, m.curvesErr
	}
	return &m.curves, nil
}

func (m *mockApiClient) GetCurve(label string) (*client.Curve, error) {
	if m.curvesErr != nil {
		return nil, m.curvesErr
	}
	c, ok := m.curves[label]
	if !ok {
		return nil, errors.New("not found")
	}
	return c, nil
}

func (m *mockApiClient) GetSensors() (*map[string]*client.Sensor, error) {
	if m.sensorsErr != nil {
		return nil, m.sensorsErr
	}
	return &m.sensors, nil
}

func (m *mockApiClient) GetSensor(label string) (*client.Sensor, error) {
	if m.sensorsErr != nil {
		return nil, m.sensorsErr
	}
	s, ok := m.sensors[label]
	if !ok {
		return nil, errors.New("not found")
	}
	return s, nil
}

func TestPoller_FetchAndUpdate(t *testing.T) {
	store := NewStore()
	mockClient := &mockApiClient{
		fans: map[string]*client.Fan{
			"fan-1": {Config: client.FanConfig{ID: "fan-1"}, Rpm: 1200},
		},
		curves: map[string]*client.Curve{
			"curve-1": {Config: client.CurveConfig{ID: "curve-1"}, Value: 100},
		},
		sensors: map[string]*client.Sensor{
			"sensor-1": {Config: client.SensorConfig{ID: "sensor-1"}, MovingAvg: 40000},
		},
	}

	updated := false
	var passedErr error
	onUpdate := func(err error) {
		updated = true
		passedErr = err
	}

	poller := NewPoller(mockClient, store, onUpdate)
	poller.FetchAndUpdate()

	if !updated {
		t.Fatalf("expected onUpdate to be called")
	}
	if passedErr != nil {
		t.Fatalf("expected no error, got: %v", passedErr)
	}

	fans := store.GetFans()
	if len(fans) != 1 || fans["fan-1"].Rpm != 1200 {
		t.Fatalf("expected fan-1 with RPM 1200")
	}

	curves := store.GetCurves()
	if len(curves) != 1 || curves["curve-1"].Value != 100 {
		t.Fatalf("expected curve-1 with Value 100")
	}

	sensors := store.GetSensors()
	if len(sensors) != 1 || sensors["sensor-1"].MovingAvg != 40000 {
		t.Fatalf("expected sensor-1 with MovingAvg 40000")
	}
}

func TestPoller_FetchAndUpdate_WithErrors(t *testing.T) {
	store := NewStore()
	expectedErr := errors.New("fan error")
	mockClient := &mockApiClient{
		fansErr: expectedErr,
	}

	updated := false
	var passedErr error
	onUpdate := func(err error) {
		updated = true
		passedErr = err
	}

	poller := NewPoller(mockClient, store, onUpdate)
	poller.FetchAndUpdate()

	if !updated {
		t.Fatalf("expected onUpdate to be called")
	}
	if passedErr != expectedErr {
		t.Fatalf("expected fan error, got: %v", passedErr)
	}
}
