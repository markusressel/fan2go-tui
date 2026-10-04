package demo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	host       string
	port       int
	simulator  *Simulator
	tickPeriod time.Duration

	listener   net.Listener
	httpServer *http.Server

	simCtx    context.Context
	simCancel context.CancelFunc
	simWg     sync.WaitGroup
}

func NewServer(host string, port int, simulator *Simulator) *Server {
	if simulator == nil {
		simulator = NewSimulator()
	}

	return &Server{
		host:       host,
		port:       port,
		simulator:  simulator,
		tickPeriod: 250 * time.Millisecond,
	}
}

func (s *Server) SetTickPeriod(period time.Duration) {
	if period > 0 {
		s.tickPeriod = period
	}
}

func (s *Server) Simulator() *Simulator {
	return s.simulator
}

func (s *Server) Port() int {
	if s.listener == nil {
		return s.port
	}
	return s.listener.Addr().(*net.TCPAddr).Port
}

func (s *Server) Addr() string {
	if s.listener == nil {
		return net.JoinHostPort(s.host, strconv.Itoa(s.port))
	}
	return s.listener.Addr().String()
}

func (s *Server) createMux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /fan", s.handleFans)
	mux.HandleFunc("GET /fan/{id}", s.handleFan)

	mux.HandleFunc("GET /curve", s.handleCurves)
	mux.HandleFunc("GET /curve/{id}", s.handleCurve)

	mux.HandleFunc("GET /sensor", s.handleSensors)
	mux.HandleFunc("GET /sensor/{id}", s.handleSensor)

	mux.HandleFunc("GET /", s.handleIndex)

	return mux
}

func (s *Server) Start() error {
	addr := net.JoinHostPort(s.host, strconv.Itoa(s.port))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	s.httpServer = &http.Server{
		Handler: s.createMux(),
	}

	s.simCtx, s.simCancel = context.WithCancel(context.Background())
	s.simWg.Add(1)
	go s.runSimulation(s.simCtx)

	go func() {
		_ = s.httpServer.Serve(listener)
	}()

	return nil
}

func (s *Server) runSimulation(ctx context.Context) {
	defer s.simWg.Done()

	ticker := time.NewTicker(s.tickPeriod)
	defer ticker.Stop()

	lastTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			dt := now.Sub(lastTime).Seconds()
			lastTime = now
			if dt > 1.0 {
				dt = 1.0
			}
			s.simulator.Tick(dt)
		}
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.simCancel != nil {
		s.simCancel()
		s.simWg.Wait()
	}

	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}

	return nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (s *Server) handleFans(w http.ResponseWriter, r *http.Request) {
	s.simulator.mutex.RLock()
	defer s.simulator.mutex.RUnlock()
	writeJSON(w, http.StatusOK, s.simulator.fans)
}

func (s *Server) handleFan(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	if rawID == "" {
		s.handleFans(w, r)
		return
	}

	id, err := url.PathUnescape(rawID)
	if err != nil {
		id = rawID
	}

	s.simulator.mutex.RLock()
	fan, ok := s.simulator.fans[id]
	s.simulator.mutex.RUnlock()

	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "fan not found"})
		return
	}

	writeJSON(w, http.StatusOK, fan)
}

func (s *Server) handleCurves(w http.ResponseWriter, r *http.Request) {
	s.simulator.mutex.RLock()
	defer s.simulator.mutex.RUnlock()
	writeJSON(w, http.StatusOK, s.simulator.curves)
}

func (s *Server) handleCurve(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	if rawID == "" {
		s.handleCurves(w, r)
		return
	}

	id, err := url.PathUnescape(rawID)
	if err != nil {
		id = rawID
	}

	s.simulator.mutex.RLock()
	curve, ok := s.simulator.curves[id]
	s.simulator.mutex.RUnlock()

	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "curve not found"})
		return
	}

	writeJSON(w, http.StatusOK, curve)
}

func (s *Server) handleSensors(w http.ResponseWriter, r *http.Request) {
	s.simulator.mutex.RLock()
	defer s.simulator.mutex.RUnlock()
	writeJSON(w, http.StatusOK, s.simulator.sensors)
}

func (s *Server) handleSensor(w http.ResponseWriter, r *http.Request) {
	rawID := r.PathValue("id")
	if rawID == "" {
		s.handleSensors(w, r)
		return
	}

	id, err := url.PathUnescape(rawID)
	if err != nil {
		id = rawID
	}

	s.simulator.mutex.RLock()
	sensor, ok := s.simulator.sensors[id]
	s.simulator.mutex.RUnlock()

	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "sensor not found"})
		return
	}

	writeJSON(w, http.StatusOK, sensor)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"service": "fan2go-demo-server",
		"status":  "running",
		"endpoints": []string{
			"/fan",
			"/fan/{id}",
			"/curve",
			"/curve/{id}",
			"/sensor",
			"/sensor/{id}",
		},
	})
}

var ErrServerNotRunning = errors.New("server is not running")
