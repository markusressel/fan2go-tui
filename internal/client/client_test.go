package client

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
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
