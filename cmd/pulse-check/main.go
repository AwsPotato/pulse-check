package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AwsPotato/pulse-check/internal/monitor"
	"github.com/AwsPotato/pulse-check/internal/server"
)

func main() {
	// Initialize the monitor service
	mon := monitor.NewMonitor()
	
	// Start the background worker (polls every 5 seconds)
	mon.Start(5 * time.Second)
	defer mon.Stop()

	// Initialize the server and inject the monitor dependency
	srv := server.NewServer(mon)

	// Set up the HTTP multiplexer
	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	port := ":8080"
	log.Printf("Starting Pulse-Check server on port %s", port)
	
	// Start the HTTP server
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
