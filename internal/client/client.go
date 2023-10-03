package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Fan struct {
	Label string `json:"label"`
}
type Curve struct {
}
type Sensor struct {
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
	GetFans() map[string]Fan
	GetCurves() map[string]Curve
	GetSensors() map[string]Sensor
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

func (client *Fan2goApiClientEcho) GetFans() map[string]Fan {
	url := fmt.Sprintf("http://%s/fan", client.baseUrl)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error is req: ", err)
	}

	// Do sends an HTTP request and
	resp, err := client.webclient.Do(req)
	if err != nil {
		fmt.Println("error in send req: ", err)
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the data with the data from the JSON
	var data map[string]Fan

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
	}

	return data
}

func (client *Fan2goApiClientEcho) GetCurves() map[string]Curve {
	url := fmt.Sprintf("http://%s/curve", client.baseUrl)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error is req: ", err)
	}

	// Do sends an HTTP request and
	resp, err := client.webclient.Do(req)
	if err != nil {
		fmt.Println("error in send req: ", err)
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the data with the data from the JSON
	var data map[string]Curve

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
	}

	return data
}
func (client *Fan2goApiClientEcho) GetSensors() map[string]Sensor {
	url := fmt.Sprintf("http://%s/sensor", client.baseUrl)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error is req: ", err)
	}

	// Do sends an HTTP request and
	resp, err := client.webclient.Do(req)
	if err != nil {
		fmt.Println("error in send req: ", err)
	}

	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the data with the data from the JSON
	var data map[string]Sensor

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
	}

	return data
}
