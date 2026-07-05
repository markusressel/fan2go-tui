package state

import (
	"fan2go-tui/internal/client"
	"sync"
)

type FanState struct {
	Fan       *client.Fan
	RpmValues []float64
	PwmValues []float64
}

type CurveState struct {
	Curve  *client.Curve
	Values []float64
}

type SensorState struct {
	Sensor *client.Sensor
	Values []float64
}

type Store struct {
	mutex   sync.RWMutex
	Fans    map[string]*FanState
	Curves  map[string]*CurveState
	Sensors map[string]*SensorState

	maxHistory int
}

func NewStore() *Store {
	return &Store{
		Fans:       make(map[string]*FanState),
		Curves:     make(map[string]*CurveState),
		Sensors:    make(map[string]*SensorState),
		maxHistory: 1000,
	}
}

func (s *Store) UpdateFans(fans map[string]*client.Fan) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, fan := range fans {
		if s.Fans[id] == nil {
			s.Fans[id] = &FanState{
				Fan:       fan,
				RpmValues: make([]float64, 0),
				PwmValues: make([]float64, 0),
			}
		} else {
			s.Fans[id].Fan = fan
		}

		fs := s.Fans[id]
		fs.RpmValues = append(fs.RpmValues, float64(fan.Rpm))
		if len(fs.RpmValues) > s.maxHistory {
			fs.RpmValues = fs.RpmValues[len(fs.RpmValues)-s.maxHistory:]
		}

		fs.PwmValues = append(fs.PwmValues, float64(fan.Pwm))
		if len(fs.PwmValues) > s.maxHistory {
			fs.PwmValues = fs.PwmValues[len(fs.PwmValues)-s.maxHistory:]
		}
	}

	// Remove fans that no longer exist
	for id := range s.Fans {
		if _, ok := fans[id]; !ok {
			delete(s.Fans, id)
		}
	}
}

func (s *Store) UpdateCurves(curves map[string]*client.Curve) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, curve := range curves {
		if s.Curves[id] == nil {
			s.Curves[id] = &CurveState{
				Curve:  curve,
				Values: make([]float64, 0),
			}
		} else {
			s.Curves[id].Curve = curve
		}

		cs := s.Curves[id]
		cs.Values = append(cs.Values, curve.Value)
		if len(cs.Values) > s.maxHistory {
			cs.Values = cs.Values[len(cs.Values)-s.maxHistory:]
		}
	}

	for id := range s.Curves {
		if _, ok := curves[id]; !ok {
			delete(s.Curves, id)
		}
	}
}

func (s *Store) UpdateSensors(sensors map[string]*client.Sensor) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for id, sensor := range sensors {
		if s.Sensors[id] == nil {
			s.Sensors[id] = &SensorState{
				Sensor: sensor,
				Values: make([]float64, 0),
			}
		} else {
			s.Sensors[id].Sensor = sensor
		}

		ss := s.Sensors[id]
		ss.Values = append(ss.Values, sensor.MovingAvg/1000.0)
		if len(ss.Values) > s.maxHistory {
			ss.Values = ss.Values[len(ss.Values)-s.maxHistory:]
		}
	}

	for id := range s.Sensors {
		if _, ok := sensors[id]; !ok {
			delete(s.Sensors, id)
		}
	}
}

func (s *Store) GetFanState(id string) *FanState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Fans[id]
}
func (s *Store) GetCurveState(id string) *CurveState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Curves[id]
}
func (s *Store) GetSensorState(id string) *SensorState {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Sensors[id]
}

func (s *Store) GetFans() map[string]*client.Fan {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	result := make(map[string]*client.Fan)
	for id, fs := range s.Fans {
		result[id] = fs.Fan
	}
	return result
}

func (s *Store) GetCurves() map[string]*client.Curve {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	result := make(map[string]*client.Curve)
	for id, cs := range s.Curves {
		result[id] = cs.Curve
	}
	return result
}

func (s *Store) GetSensors() map[string]*client.Sensor {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	result := make(map[string]*client.Sensor)
	for id, ss := range s.Sensors {
		result[id] = ss.Sensor
	}
	return result
}
