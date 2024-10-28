package client

import (
	"encoding/json"
	"errors"
	"fan2go-tui/internal/logging"
	"fmt"
	"io"
	"net/http"
)

type Fan struct {
	Pwm int `json:"pwm"`
	Rpm int `json:"rpm"`

	FanCurveData *map[int]float64 `json:"fanCurveData,omitempty"`

	Config FanConfig `json:"config"`
}

type FanConfig struct {
	ID          string             `json:"id"`
	NeverStop   bool               `json:"neverStop"`
	MinPwm      *int               `json:"minPwm,omitempty"`
	StartPwm    *int               `json:"startPwm,omitempty"`
	PwmMap      *map[int]int       `json:"pwmMap,omitempty"`
	MaxPwm      *int               `json:"maxPwm,omitempty"`
	Curve       string             `json:"curve"`
	HwMon       *HwMonFanConfig    `json:"hwMon,omitempty"`
	File        *FileFanConfig     `json:"file,omitempty"`
	Cmd         *CmdFanConfig      `json:"cmd,omitempty"`
	ControlLoop *ControlLoopConfig `json:"controlLoop,omitempty"`
}

type HwMonFanConfig struct {
	Platform      string `json:"platform"`
	Index         int    `json:"index"`
	RpmChannel    int    `json:"rpmChannel"`
	PwmChannel    int    `json:"pwmChannel"`
	SysfsPath     string
	RpmInputPath  string
	PwmPath       string
	PwmEnablePath string
}

type FileFanConfig struct {
	Path    string `json:"path"`
	RpmPath string `json:"rpmPath"`
}

type CmdFanConfig struct {
	SetPwm *ExecConfig `json:"setPwm,omitempty"`
	GetPwm *ExecConfig `json:"getPwm,omitempty"`
	GetRpm *ExecConfig `json:"getRpm,omitempty"`
}

type ExecConfig struct {
	Exec string   `json:"exec"`
	Args []string `json:"args"`
}

type ControlLoopConfig struct {
	P float64 `json:"p"`
	I float64 `json:"i"`
	D float64 `json:"d"`
}

type Curve struct {
	Config CurveConfig `json:"config"`
	Value  float64     `json:"value"`
}

type CurveConfig struct {
	ID       string               `json:"id"`
	Linear   *LinearCurveConfig   `json:"linear,omitempty"`
	PID      *PidCurveConfig      `json:"pid,omitempty"`
	Function *FunctionCurveConfig `json:"function,omitempty"`
}

type LinearCurveConfig struct {
	Sensor string          `json:"sensor"`
	Min    int             `json:"min"`
	Max    int             `json:"max"`
	Steps  map[int]float64 `json:"steps,omitempty"`
}

type PidCurveConfig struct {
	Sensor   string  `json:"sensor"`
	SetPoint float64 `json:"setPoint"`
	P        float64 `json:"p"`
	I        float64 `json:"i"`
	D        float64 `json:"d"`
}

const (
	// FunctionSum computes the sum of all referenced curves
	FunctionSum = "sum"
	// FunctionDifference computes the difference of all referenced curves
	FunctionDifference = "difference"
	// FunctionAverage computes the average value of all referenced
	// curves using the arithmetic mean
	FunctionAverage = "average"
	// FunctionDelta computes the difference between the biggest and the smallest
	// value of all referenced curves
	FunctionDelta = "delta"
	// FunctionMinimum computes the smallest value of all referenced curves
	FunctionMinimum = "minimum"
	// FunctionMaximum computes the biggest value of all referenced curves
	FunctionMaximum = "maximum"
)

type FunctionCurveConfig struct {
	Type   string   `json:"type"`
	Curves []string `json:"curves"`
}

type Sensor struct {
	Name      string       `json:"name"`
	Config    SensorConfig `json:"configuration"`
	MovingAvg float64      `json:"movingAvg"`
}

type SensorConfig struct {
	ID    string             `json:"id"`
	HwMon *HwMonSensorConfig `json:"hwMon,omitempty"`
	File  *FileSensorConfig  `json:"file,omitempty"`
	Cmd   *CmdSensorConfig   `json:"cmd,omitempty"`
}

type HwMonSensorConfig struct {
	Platform  string `json:"platform"`
	Index     int    `json:"index"`
	TempInput string
}

type FileSensorConfig struct {
	Path string `json:"path"`
}

type CmdSensorConfig struct {
	Exec string   `json:"exec"`
	Args []string `json:"args"`
}

type HwMonFan struct {
	Label        string  `json:"label"`
	Index        int     `json:"index"`
	RpmMovingAvg float64 `json:"rpmMovingAvg"`
	//Config       configuration.FanConfig `json:"config"`
	MinPwm       *int             `json:"minPwm"`
	StartPwm     *int             `json:"startPwm"`
	MaxPwm       *int             `json:"maxPwm"`
	FanCurveData *map[int]float64 `json:"fanCurveData"`
	Rpm          int              `json:"rpm"`
	Pwm          int              `json:"pwm"`
}

type Fan2goApiClient interface {
	GetFans() (*map[string]*Fan, error)
	GetFan(label string) (*Fan, error)

	GetCurves() (*map[string]*Curve, error)
	GetCurve(label string) (*Curve, error)

	GetSensors() (*map[string]*Sensor, error)
	GetSensor(label string) (*Sensor, error)
}

type Fan2goApiClientEcho struct {
	baseUrl   string
	webclient *http.Client
}

func NewApiClient(baseUrl string, port int) Fan2goApiClient {
	baseUrl = fmt.Sprintf("%s:%d", baseUrl, port)
	return &Fan2goApiClientEcho{
		baseUrl:   baseUrl,
		webclient: createWebserver(),
	}
}

func createWebserver() *http.Client {
	return &http.Client{}
}

func (client *Fan2goApiClientEcho) GetFans() (*map[string]*Fan, error) {
	url := fmt.Sprintf("http://%s/fan", client.baseUrl)

	var data map[string]*Fan
	return doGet(client.webclient, url, data)
}

func (client *Fan2goApiClientEcho) GetFan(label string) (*Fan, error) {
	url := fmt.Sprintf("http://%s/fan/%s", client.baseUrl, label)

	var data *Fan
	result, err := doGet(client.webclient, url, data)
	if err != nil {
		return nil, err
	} else {
		return *result, err
	}
}

func (client *Fan2goApiClientEcho) GetCurves() (*map[string]*Curve, error) {
	url := fmt.Sprintf("http://%s/curve", client.baseUrl)

	var data map[string]*Curve
	result, err := doGet(client.webclient, url, data)
	if err != nil {
		return nil, err
	} else {
		return result, err
	}
}

func (client *Fan2goApiClientEcho) GetCurve(label string) (*Curve, error) {
	url := fmt.Sprintf("http://%s/curve/%s", client.baseUrl, label)

	var data *Curve
	result, err := doGet(client.webclient, url, data)
	if err != nil {
		return nil, err
	} else {
		return *result, err
	}
}

func (client *Fan2goApiClientEcho) GetSensors() (*map[string]*Sensor, error) {
	url := fmt.Sprintf("http://%s/sensor", client.baseUrl)

	var data map[string]*Sensor
	result, err := doGet(client.webclient, url, data)
	if err != nil {
		return nil, err
	} else {
		return result, err
	}
}

func (client *Fan2goApiClientEcho) GetSensor(label string) (*Sensor, error) {
	url := fmt.Sprintf("http://%s/sensor/%s", client.baseUrl, label)

	var data *Sensor
	result, err := doGet(client.webclient, url, data)
	if err != nil {
		return nil, err
	} else {
		return *result, err
	}
}

func doGet[T any](client *http.Client, url string, data T) (*T, error) {

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logging.Warning("Error is req: %v", err)
	}

	// Send it
	resp, err := client.Do(req)
	if err != nil {
		logging.Warning("error in send req: %v", err)
		return nil, err
	}

	// Defer the closing of the body
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.Warning("error in close resp: %v", err)
		}
	}(resp.Body)

	// Check the status code
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New(fmt.Sprintf("Cannot reach fan2go daemon, did you enable its API?"))
	} else if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Unexpected API status code: %s", resp.Status))
	}

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logging.Warning("%v", err)
	}

	return &data, err
}
