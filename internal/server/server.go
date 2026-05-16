package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AwsPotato/pulse-check/internal/monitor"
)

// Server handles HTTP requests
type Server struct {
	monitor *monitor.Monitor
}

// NewServer creates a new instance of the Server with its dependencies
func NewServer(m *monitor.Monitor) *Server {
	return &Server{
		monitor: m,
	}
}

// RegisterRoutes sets up the routing for the server
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/metrics", s.handleMetrics)
	mux.HandleFunc("/health", s.handleHealth)
}

func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := s.monitor.GetStats()
	w.Header().Set("Content-Type", "application/json")
	
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
