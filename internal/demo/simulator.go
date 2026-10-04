package demo

import (
	"fan2go-tui/internal/client"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Simulator struct {
	mutex   sync.RWMutex
	fans    map[string]*client.Fan
	curves  map[string]*client.Curve
	sensors map[string]*client.Sensor

	elapsed    float64
	pidState   *pidController
	cpuFanCal  map[int]float64
	gpuFanCal  map[int]float64
	rng        *rand.Rand
}

type pidController struct {
	integral   float64
	lastError  float64
}

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}

func NewSimulator() *Simulator {
	s := &Simulator{
		fans:      make(map[string]*client.Fan),
		curves:    make(map[string]*client.Curve),
		sensors:   make(map[string]*client.Sensor),
		pidState:  &pidController{},
		rng:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	s.initData()
	return s
}

func (s *Simulator) initData() {
	s.sensors = map[string]*client.Sensor{
		"cpu-temp": {
			Name:      "AMD Ryzen 9 7950X - Tdie",
			MovingAvg: 48500.0,
			Config: client.SensorConfig{
				ID: "cpu-temp",
				HwMon: &client.HwMonSensorConfig{
					Platform:  "k10temp",
					Index:     1,
					Channel:   1,
					TempInput: "temp1_input",
				},
			},
		},
		"gpu-temp": {
			Name:      "NVIDIA GeForce RTX 4090 - GPU Core",
			MovingAvg: 42000.0,
			Config: client.SensorConfig{
				ID: "gpu-temp",
				Nvidia: &client.NvidiaSensorConfig{
					Device: "NVIDIA GeForce RTX 4090",
					Index:  0,
				},
			},
		},
		"case-temp": {
			Name:      "Motherboard Ambient",
			MovingAvg: 31200.0,
			Config: client.SensorConfig{
				ID: "case-temp",
				HwMon: &client.HwMonSensorConfig{
					Platform:  "nct6798",
					Index:     0,
					Channel:   2,
					TempInput: "temp2_input",
				},
			},
		},
		"nvme-temp": {
			Name:      "Samsung SSD 990 PRO 2TB",
			MovingAvg: 43500.0,
			Config: client.SensorConfig{
				ID: "nvme-temp",
				Disk: &client.DiskSensorConfig{
					Device: "/dev/nvme0n1",
				},
			},
		},
	}

	s.curves = map[string]*client.Curve{
		"cpu-linear-curve": {
			Value: 110.0,
			Config: client.CurveConfig{
				ID: "cpu-linear-curve",
				Linear: &client.LinearCurveConfig{
					Sensor: "cpu-temp",
					Min:    40,
					Max:    80,
				},
			},
		},
		"gpu-stepped-curve": {
			Value: 75.0,
			Config: client.CurveConfig{
				ID: "gpu-stepped-curve",
				Linear: &client.LinearCurveConfig{
					Sensor: "gpu-temp",
					Steps: map[int]float64{
						40: 0.0,
						50: 75.0,
						65: 155.0,
						80: 255.0,
					},
				},
			},
		},
		"case-function-curve": {
			Value: 95.0,
			Config: client.CurveConfig{
				ID: "case-function-curve",
				Function: &client.FunctionCurveConfig{
					Type:   client.FunctionMaximum,
					Curves: []string{"cpu-linear-curve", "gpu-stepped-curve"},
				},
			},
		},
		"exhaust-pid-curve": {
			Value: 80.0,
			Config: client.CurveConfig{
				ID: "exhaust-pid-curve",
				PID: &client.PidCurveConfig{
					Sensor:   "case-temp",
					SetPoint: 32.0,
					P:        6.0,
					I:        0.2,
					D:        1.5,
				},
			},
		},
	}

	s.cpuFanCal = map[int]float64{
		0:   0,
		40:  550,
		80:  850,
		120: 1200,
		160: 1550,
		200: 1850,
		255: 2200,
	}

	s.gpuFanCal = map[int]float64{
		0:   0,
		50:  700,
		100: 1150,
		150: 1600,
		200: 2100,
		255: 2600,
	}

	cpuFanCalCopy := make(map[int]float64, len(s.cpuFanCal))
	for k, v := range s.cpuFanCal {
		cpuFanCalCopy[k] = v
	}

	gpuFanCalCopy := make(map[int]float64, len(s.gpuFanCal))
	for k, v := range s.gpuFanCal {
		gpuFanCalCopy[k] = v
	}

	activeMode := client.ControlModeValue("curve")

	s.fans = map[string]*client.Fan{
		"cpu-fan": {
			Pwm:          110,
			Rpm:          1120,
			FanCurveData: &cpuFanCalCopy,
			Config: client.FanConfig{
				ID:        "cpu-fan",
				Curve:     "cpu-linear-curve",
				NeverStop: true,
				MinPwm:    intPtr(40),
				StartPwm:  intPtr(60),
				MaxPwm:    intPtr(255),
				HwMon: &client.HwMonFanConfig{
					Platform:      "nct6798",
					Index:         1,
					RpmChannel:    1,
					PwmChannel:    1,
					SysfsPath:     "/sys/devices/platform/nct6775.656/hwmon/hwmon2",
					PwmPath:       "/sys/devices/platform/nct6775.656/hwmon/hwmon2/pwm1",
					RpmInputPath:  "/sys/devices/platform/nct6775.656/hwmon/hwmon2/fan1_input",
					PwmEnablePath: "/sys/devices/platform/nct6775.656/hwmon/hwmon2/pwm1_enable",
				},
				ControlMode: &client.ControlModeConfig{
					Active: &activeMode,
					OnExit: &client.OnExitConfig{
						Restore: &client.OnExitRestoreConfig{},
					},
				},
				ControlAlgorithm: &client.ControlAlgorithmConfig{
					Direct: &client.DirectControlAlgorithmConfig{
						MaxPwmChangePerCycle: intPtr(15),
					},
				},
				SanityCheck: &client.SanityCheckConfig{
					PwmValueChangedByThirdParty: client.PwmValueChangedByThirdPartyConfig{
						Enabled: true,
					},
					FanModeChangedByThirdParty: client.FanModeChangedByThirdPartyConfig{
						Enabled:          true,
						ThrottleDuration: 5 * time.Second,
					},
				},
			},
		},
		"case-intake": {
			Pwm:          95,
			Rpm:          980,
			FanCurveData: nil,
			Config: client.FanConfig{
				ID:        "case-intake",
				Curve:     "case-function-curve",
				NeverStop: false,
				MinPwm:    intPtr(40),
				StartPwm:  intPtr(60),
				MaxPwm:    intPtr(255),
				HwMon: &client.HwMonFanConfig{
					Platform:   "nct6798",
					Index:      1,
					RpmChannel: 2,
					PwmChannel: 2,
				},
			},
		},
		"case-exhaust": {
			Pwm:          80,
			Rpm:          820,
			FanCurveData: nil,
			Config: client.FanConfig{
				ID:        "case-exhaust",
				Curve:     "exhaust-pid-curve",
				NeverStop: false,
				MinPwm:    intPtr(45),
				StartPwm:  intPtr(65),
				MaxPwm:    intPtr(255),
				File: &client.FileFanConfig{
					Path:    "/sys/class/hwmon/hwmon3/pwm1",
					RpmPath: "/sys/class/hwmon/hwmon3/fan1_input",
				},
			},
		},
		"gpu-fan": {
			Pwm:          75,
			Rpm:          920,
			FanCurveData: &gpuFanCalCopy,
			Config: client.FanConfig{
				ID:        "gpu-fan",
				Curve:     "gpu-stepped-curve",
				NeverStop: false,
				MinPwm:    intPtr(0),
				StartPwm:  intPtr(50),
				MaxPwm:    intPtr(255),
				Nvidia: &client.NvidiaFanConfig{
					Device: "NVIDIA GeForce RTX 4090",
					Index:  0,
				},
			},
		},
	}
}

func (s *Simulator) Tick(dt float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.elapsed += dt

	cpuTemp := 48.0 + 16.0*math.Sin(s.elapsed*0.35) + 6.0*math.Sin(s.elapsed*0.8) + (s.rng.Float64()-0.5)*0.6
	gpuTemp := 42.0 + 18.0*math.Sin(s.elapsed*0.25+0.8) + 4.0*math.Cos(s.elapsed*0.6) + (s.rng.Float64()-0.5)*0.5
	caseTemp := 31.0 + 4.0*math.Sin(s.elapsed*0.1) + (s.rng.Float64()-0.5)*0.3
	nvmeTemp := 43.0 + 5.0*math.Sin(s.elapsed*0.15+1.5) + (s.rng.Float64()-0.5)*0.4

	cpuTemp = math.Max(35.0, math.Min(85.0, cpuTemp))
	gpuTemp = math.Max(32.0, math.Min(88.0, gpuTemp))
	caseTemp = math.Max(24.0, math.Min(45.0, caseTemp))
	nvmeTemp = math.Max(35.0, math.Min(65.0, nvmeTemp))

	s.sensors["cpu-temp"].MovingAvg = cpuTemp * 1000.0
	s.sensors["gpu-temp"].MovingAvg = gpuTemp * 1000.0
	s.sensors["case-temp"].MovingAvg = caseTemp * 1000.0
	s.sensors["nvme-temp"].MovingAvg = nvmeTemp * 1000.0

	// 1. cpu-linear-curve (min 40, max 80, output 50 to 255)
	var cpuCurveVal float64
	if cpuTemp <= 40.0 {
		cpuCurveVal = 50.0
	} else if cpuTemp >= 80.0 {
		cpuCurveVal = 255.0
	} else {
		cpuCurveVal = 50.0 + ((cpuTemp-40.0)/40.0)*(255.0-50.0)
	}
	s.curves["cpu-linear-curve"].Value = math.Round(cpuCurveVal*10) / 10

	// 2. gpu-stepped-curve: steps {40: 0, 50: 75, 65: 155, 80: 255}
	var gpuCurveVal float64
	if gpuTemp <= 40.0 {
		gpuCurveVal = 0.0
	} else if gpuTemp <= 50.0 {
		gpuCurveVal = ((gpuTemp - 40.0) / 10.0) * 75.0
	} else if gpuTemp <= 65.0 {
		gpuCurveVal = 75.0 + ((gpuTemp-50.0)/15.0)*(155.0-75.0)
	} else if gpuTemp <= 80.0 {
		gpuCurveVal = 155.0 + ((gpuTemp-65.0)/15.0)*(255.0-155.0)
	} else {
		gpuCurveVal = 255.0
	}
	s.curves["gpu-stepped-curve"].Value = math.Round(gpuCurveVal*10) / 10

	// 3. case-function-curve: max(cpuCurve, gpuCurve) * 0.85
	caseCurveVal := math.Max(cpuCurveVal, gpuCurveVal) * 0.85
	caseCurveVal = math.Max(40.0, math.Min(255.0, caseCurveVal))
	s.curves["case-function-curve"].Value = math.Round(caseCurveVal*10) / 10

	// 4. exhaust-pid-curve: PID on caseTemp with setpoint 32.0
	setPoint := 32.0
	err := caseTemp - setPoint
	s.pidState.integral += err * dt
	s.pidState.integral = math.Max(-50.0, math.Min(50.0, s.pidState.integral))
	derivative := (err - s.pidState.lastError) / dt
	s.pidState.lastError = err

	pidVal := 70.0 + 6.0*err + 0.2*s.pidState.integral + 1.5*derivative
	pidVal = math.Max(45.0, math.Min(255.0, pidVal))
	s.curves["exhaust-pid-curve"].Value = math.Round(pidVal*10) / 10

	// Update fans towards curve values
	// cpu-fan
	cpuTargetPwm := int(math.Round(cpuCurveVal))
	s.fans["cpu-fan"].Pwm = approach(s.fans["cpu-fan"].Pwm, cpuTargetPwm, 0.25)
	s.fans["cpu-fan"].Rpm = interpolateRpm(s.fans["cpu-fan"].Pwm, s.cpuFanCal) + int((s.rng.Float64()-0.5)*16.0)

	// case-intake
	intakeTargetPwm := int(math.Round(caseCurveVal))
	s.fans["case-intake"].Pwm = approach(s.fans["case-intake"].Pwm, intakeTargetPwm, 0.20)
	if s.fans["case-intake"].Pwm < 40 {
		s.fans["case-intake"].Rpm = 0
	} else {
		ratio := float64(s.fans["case-intake"].Pwm-40) / float64(255-40)
		s.fans["case-intake"].Rpm = int(450.0+ratio*(1550.0-450.0)) + int((s.rng.Float64()-0.5)*18.0)
	}

	// case-exhaust
	exhaustTargetPwm := int(math.Round(pidVal))
	s.fans["case-exhaust"].Pwm = approach(s.fans["case-exhaust"].Pwm, exhaustTargetPwm, 0.20)
	if s.fans["case-exhaust"].Pwm < 45 {
		s.fans["case-exhaust"].Rpm = 0
	} else {
		ratio := float64(s.fans["case-exhaust"].Pwm-45) / float64(255-45)
		s.fans["case-exhaust"].Rpm = int(500.0+ratio*(1400.0-500.0)) + int((s.rng.Float64()-0.5)*16.0)
	}

	// gpu-fan (supports fan stop)
	gpuTargetPwm := int(math.Round(gpuCurveVal))
	if gpuTargetPwm <= 10 {
		gpuTargetPwm = 0
	}
	s.fans["gpu-fan"].Pwm = approach(s.fans["gpu-fan"].Pwm, gpuTargetPwm, 0.30)
	if s.fans["gpu-fan"].Pwm <= 0 {
		s.fans["gpu-fan"].Rpm = int(float64(s.fans["gpu-fan"].Rpm) * 0.75)
		if s.fans["gpu-fan"].Rpm < 30 {
			s.fans["gpu-fan"].Rpm = 0
		}
	} else {
		s.fans["gpu-fan"].Rpm = interpolateRpm(s.fans["gpu-fan"].Pwm, s.gpuFanCal) + int((s.rng.Float64()-0.5)*20.0)
	}
}

func approach(current, target int, rate float64) int {
	if current == target {
		return current
	}
	diff := float64(target - current)
	step := diff * rate
	if math.Abs(diff) <= 1.0 {
		return target
	}
	if math.Abs(step) < 1.0 {
		if diff > 0 {
			step = 1.0
		} else {
			step = -1.0
		}
	}
	return current + int(math.Round(step))
}

func interpolateRpm(pwm int, curve map[int]float64) int {
	if pwm <= 0 {
		return 0
	}

	var p1, p2 int = -1, -1
	for k := range curve {
		if k <= pwm {
			if p1 == -1 || k > p1 {
				p1 = k
			}
		}
		if k >= pwm {
			if p2 == -1 || k < p2 {
				p2 = k
			}
		}
	}

	if p1 == -1 {
		return int(curve[p2])
	}
	if p2 == -1 || p1 == p2 {
		return int(curve[p1])
	}

	ratio := float64(pwm-p1) / float64(p2-p1)
	rpm := curve[p1] + ratio*(curve[p2]-curve[p1])
	return int(math.Round(rpm))
}

func (s *Simulator) GetFans() map[string]*client.Fan {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make(map[string]*client.Fan, len(s.fans))
	for k, v := range s.fans {
		fanCopy := *v
		result[k] = &fanCopy
	}
	return result
}

func (s *Simulator) GetFan(id string) (*client.Fan, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	fan, ok := s.fans[id]
	if !ok {
		return nil, false
	}
	fanCopy := *fan
	return &fanCopy, true
}

func (s *Simulator) GetCurves() map[string]*client.Curve {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make(map[string]*client.Curve, len(s.curves))
	for k, v := range s.curves {
		curveCopy := *v
		result[k] = &curveCopy
	}
	return result
}

func (s *Simulator) GetCurve(id string) (*client.Curve, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	curve, ok := s.curves[id]
	if !ok {
		return nil, false
	}
	curveCopy := *curve
	return &curveCopy, true
}

func (s *Simulator) GetSensors() map[string]*client.Sensor {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	result := make(map[string]*client.Sensor, len(s.sensors))
	for k, v := range s.sensors {
		sensorCopy := *v
		result[k] = &sensorCopy
	}
	return result
}

func (s *Simulator) GetSensor(id string) (*client.Sensor, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	sensor, ok := s.sensors[id]
	if !ok {
		return nil, false
	}
	sensorCopy := *sensor
	return &sensorCopy, true
}
