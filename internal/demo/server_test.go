package demo

import (
	"context"
	"fan2go-tui/internal/client"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestServer_IntegrationWithApiClient(t *testing.T) {
	sim := NewSimulator()
	server := NewServer("127.0.0.1", 0, sim)
	server.SetTickPeriod(50 * time.Millisecond)

	if err := server.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()

	port := server.Port()
	if port <= 0 {
		t.Fatalf("invalid server port: %d", port)
	}

	apiClient := client.NewApiClient("127.0.0.1", port)

	// 1. GetFans
	fans, err := apiClient.GetFans()
	if err != nil {
		t.Fatalf("GetFans failed: %v", err)
	}
	if fans == nil || len(*fans) != 4 {
		t.Fatalf("expected 4 fans, got %v", fans)
	}

	// 2. GetFan
	cpuFan, err := apiClient.GetFan("cpu-fan")
	if err != nil {
		t.Fatalf("GetFan failed: %v", err)
	}
	if cpuFan == nil || cpuFan.Config.ID != "cpu-fan" {
		t.Fatalf("unexpected cpu-fan: %v", cpuFan)
	}

	// 3. GetCurves
	curves, err := apiClient.GetCurves()
	if err != nil {
		t.Fatalf("GetCurves failed: %v", err)
	}
	if curves == nil || len(*curves) != 4 {
		t.Fatalf("expected 4 curves, got %v", curves)
	}

	// 4. GetCurve
	cpuCurve, err := apiClient.GetCurve("cpu-linear-curve")
	if err != nil {
		t.Fatalf("GetCurve failed: %v", err)
	}
	if cpuCurve == nil || cpuCurve.Config.ID != "cpu-linear-curve" {
		t.Fatalf("unexpected cpu-linear-curve: %v", cpuCurve)
	}

	// 5. GetSensors
	sensors, err := apiClient.GetSensors()
	if err != nil {
		t.Fatalf("GetSensors failed: %v", err)
	}
	if sensors == nil || len(*sensors) != 4 {
		t.Fatalf("expected 4 sensors, got %v", sensors)
	}

	// 6. GetSensor
	cpuTemp, err := apiClient.GetSensor("cpu-temp")
	if err != nil {
		t.Fatalf("GetSensor failed: %v", err)
	}
	if cpuTemp == nil || cpuTemp.Config.ID != "cpu-temp" {
		t.Fatalf("unexpected cpu-temp: %v", cpuTemp)
	}

	// 7. Nonexistent items
	_, err = apiClient.GetFan("nonexistent-fan")
	if err == nil {
		t.Fatalf("expected error for nonexistent fan, got nil")
	}

	// 8. Index endpoint
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/", port))
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK on /, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if len(body) == 0 {
		t.Fatalf("expected non-empty body on /")
	}

	// Let simulation tick for a bit and verify values evolve
	time.Sleep(200 * time.Millisecond)

	updatedCpuFan, err := apiClient.GetFan("cpu-fan")
	if err != nil {
		t.Fatalf("GetFan after ticks failed: %v", err)
	}
	if updatedCpuFan.Pwm < 0 || updatedCpuFan.Pwm > 255 {
		t.Fatalf("updated PWM out of range: %d", updatedCpuFan.Pwm)
	}
}
