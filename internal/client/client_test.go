package client

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestJoinHostPort(t *testing.T) {
	for _, test := range []struct {
		host     string
		expected string
	}{
		{"127.0.0.1", "127.0.0.1:9001"},
		{"localhost", "localhost:9001"},
		{"[::1]", "[::1]:9001"},
		{"::1", "[::1]:9001"},
		// link-local, with the zone index (the interface), like the API host in the fan2go config
		{"[fe80::18c2:b5ff:fec0:71b8%ens19]", "[fe80::18c2:b5ff:fec0:71b8%ens19]:9001"},
		{"fe80::18c2:b5ff:fec0:71b8%ens19", "[fe80::18c2:b5ff:fec0:71b8%ens19]:9001"},
	} {
		if actual := joinHostPort(test.host, 9001); actual != test.expected {
			t.Errorf("joinHostPort(%q) = %q, expected %q", test.host, actual, test.expected)
		}
	}
}

// Regression (https://github.com/markusressel/fan2go-tui/issues/120): the zone index of a link-local address was put
// into the URL as is ("%ens19"), which is an invalid escape. The request could not be built, and sending it panicked.
func TestEndpointUrl(t *testing.T) {
	client := NewApiClient("[fe80::18c2:b5ff:fec0:71b8%ens19]", 9001).(*Fan2goApiClientEcho)

	endpoint := client.endpointUrl("fan", "cpu fan")
	if expected := "http://[fe80::18c2:b5ff:fec0:71b8%25ens19]:9001/fan/cpu%20fan"; endpoint != expected {
		t.Errorf("endpointUrl = %q, expected %q", endpoint, expected)
	}
	parsed, err := url.Parse(endpoint)
	if err != nil {
		t.Fatalf("the URL is invalid: %v", err)
	}
	if expected := "[fe80::18c2:b5ff:fec0:71b8%ens19]:9001"; parsed.Host != expected {
		t.Errorf("the host to connect to is %q, expected %q (with the zone index)", parsed.Host, expected)
	}
	if expected := "/fan/cpu fan"; parsed.Path != expected {
		t.Errorf("the path is %q, expected %q", parsed.Path, expected)
	}
}

// Regression (https://github.com/markusressel/fan2go-tui/issues/120): a request that could not be built was sent
// anyway, which panicked.
func TestDoGetWithAnInvalidUrl(t *testing.T) {
	var data map[string]*Fan
	result, err := doGet(&http.Client{}, "http://[fe80::1%ens19]:9001/fan", data)
	if err == nil || result != nil {
		t.Errorf("expected an error, got %v, %v", result, err)
	}
}

func TestGetFans(t *testing.T) {
	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/fan" {
			http.NotFound(writer, request)
			return
		}
		_, _ = writer.Write([]byte(`{"cpu": {"label": "cpu", "pwm": 128, "rpm": 900}}`))
	})
	listen := func(t *testing.T, network string, address string) (*httptest.Server, int) {
		listener, err := net.Listen(network, address)
		if err != nil {
			t.Skipf("cannot listen on %s: %v", address, err)
		}
		server := httptest.NewUnstartedServer(handler)
		server.Listener = listener
		server.Start()
		t.Cleanup(server.Close)
		return server, listener.Addr().(*net.TCPAddr).Port
	}
	getFans := func(t *testing.T, host string, port int) {
		fans, err := NewApiClient(host, port).GetFans()
		if err != nil {
			t.Fatalf("GetFans: %v", err)
		}
		if fan := (*fans)["cpu"]; fan == nil || fan.Rpm != 900 {
			t.Errorf("unexpected fans: %v", *fans)
		}
	}

	t.Run("IPv4", func(t *testing.T) {
		_, port := listen(t, "tcp4", "127.0.0.1:0")
		getFans(t, "127.0.0.1", port)
	})
	t.Run("IPv6 with a zone index", func(t *testing.T) {
		// the loopback interface as the zone, like the interface of a link-local address
		_, port := listen(t, "tcp6", "[::1]:0")
		loopback := loopbackInterface(t)
		getFans(t, "[::1%"+loopback+"]", port)
		getFans(t, "::1%"+loopback, port)
	})
}

// loopbackInterface returns the name of the loopback interface, or skips the test.
func loopbackInterface(t *testing.T) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Skipf("cannot list the interfaces: %v", err)
	}
	for _, networkInterface := range interfaces {
		if networkInterface.Flags&net.FlagLoopback != 0 {
			return networkInterface.Name
		}
	}
	t.Skip("no loopback interface")
	return ""
}

func TestApiClient_EndpointsAndErrors(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/fan/cpu", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"config": {"id": "cpu"}, "pwm": 100, "rpm": 1200}`))
	})
	mux.HandleFunc("/curve", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"curve1": {"value": 50, "config": {"id": "curve1"}}}`))
	})
	mux.HandleFunc("/curve/curve1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"value": 50, "config": {"id": "curve1"}}`))
	})
	mux.HandleFunc("/sensor", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"sensor1": {"name": "s1", "movingAvg": 45000, "config": {"id": "sensor1"}}}`))
	})
	mux.HandleFunc("/sensor/sensor1", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"name": "s1", "movingAvg": 45000, "config": {"id": "sensor1"}}`))
	})
	mux.HandleFunc("/error-500", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("failed to parse test server url: %v", err)
	}
	host := u.Hostname()
	port, _ := strconv.Atoi(u.Port())
	api := NewApiClient(host, port)

	// GetFan
	fan, err := api.GetFan("cpu")
	if err != nil || fan == nil || fan.Rpm != 1200 {
		t.Errorf("GetFan failed: %v, %v", fan, err)
	}

	// GetCurves
	curves, err := api.GetCurves()
	if err != nil || curves == nil || (*curves)["curve1"] == nil {
		t.Errorf("GetCurves failed: %v, %v", curves, err)
	}

	// GetCurve
	curve, err := api.GetCurve("curve1")
	if err != nil || curve == nil || curve.Config.ID != "curve1" {
		t.Errorf("GetCurve failed: %v, %v", curve, err)
	}

	// GetSensors
	sensors, err := api.GetSensors()
	if err != nil || sensors == nil || (*sensors)["sensor1"] == nil {
		t.Errorf("GetSensors failed: %v, %v", sensors, err)
	}

	// GetSensor
	sensor, err := api.GetSensor("sensor1")
	if err != nil || sensor == nil || sensor.Name != "s1" {
		t.Errorf("GetSensor failed: %v, %v", sensor, err)
	}

	// 404 error
	_, err = api.GetFan("non-existent-404")
	if err == nil || !strings.Contains(err.Error(), "Cannot reach fan2go daemon") {
		t.Errorf("expected daemon 404 error, got %v", err)
	}

	// 500 error
	var dummy map[string]*Fan
	_, err = doGet(http.DefaultClient, server.URL+"/error-500", dummy)
	if err == nil || !strings.Contains(err.Error(), "Unexpected API status code") {
		t.Errorf("expected 500 status code error, got %v", err)
	}
}
