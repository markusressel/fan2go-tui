package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Fan struct {
	Label string `json:"label"`

	Pwm int `json:"pwm"`
	Rpm int `json:"rpm"`

	Config FanConfig `json:"config"`
}

type FanConfig struct {
	Id          string             `json:"id"`
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
	Steps  map[int]float64 `json:"steps"`
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
	GetFans() map[string]*Fan
	GetFan(label string) *Fan

	GetCurves() map[string]*Curve
	GetCurve(label string) *Curve

	GetSensors() map[string]*Sensor
	GetSensor(label string) *Sensor
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

func (client *Fan2goApiClientEcho) GetFans() map[string]*Fan {
	url := fmt.Sprintf("http://%s/fan", client.baseUrl)

	var data map[string]*Fan
	data = doGet(client.webclient, url, data)
	return data
}

func (client *Fan2goApiClientEcho) GetFan(label string) *Fan {
	url := fmt.Sprintf("http://%s/fan/%s", client.baseUrl, label)

	var data *Fan
	data = doGet(client.webclient, url, data)
	return data
}

func (client *Fan2goApiClientEcho) GetCurves() map[string]*Curve {
	url := fmt.Sprintf("http://%s/curve", client.baseUrl)

	var data map[string]*Curve
	data = doGet(client.webclient, url, data)
	return data
}

func (client *Fan2goApiClientEcho) GetCurve(label string) *Curve {
	url := fmt.Sprintf("http://%s/curve/%s", client.baseUrl, label)

	var data *Curve
	data = doGet(client.webclient, url, data)
	return data
}

func (client *Fan2goApiClientEcho) GetSensors() map[string]*Sensor {
	url := fmt.Sprintf("http://%s/sensor", client.baseUrl)

	var data map[string]*Sensor
	data = doGet(client.webclient, url, data)
	return data
}

func (client *Fan2goApiClientEcho) GetSensor(label string) *Sensor {
	url := fmt.Sprintf("http://%s/sensor/%s", client.baseUrl, label)

	var data *Sensor
	data = doGet(client.webclient, url, data)
	return data
}

func doGet[T any](client *http.Client, url string, data T) T {

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error is req: ", err)
	}

	// Do sends an HTTP request and
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error in send req: ", err)
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
	}

	return data
}
