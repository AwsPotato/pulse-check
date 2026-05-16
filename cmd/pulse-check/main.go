package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AwsPotato/pulse-check/internal/monitor"
	"github.com/AwsPotato/pulse-check/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pollIntervalStr := os.Getenv("POLL_INTERVAL")
	if pollIntervalStr == "" {
		pollIntervalStr = "5s"
	}

	pollInterval, err := time.ParseDuration(pollIntervalStr)
	if err != nil {
		log.Printf("Invalid POLL_INTERVAL '%s', defaulting to 5s", pollIntervalStr)
		pollInterval = 5 * time.Second
	}

	// Initialize the monitor service
	mon := monitor.NewMonitor()
	
	// Start the background worker
	mon.Start(pollInterval)

	// Initialize the server and inject the monitor dependency
	srv := server.NewServer(mon, port)

	// Create a channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start the HTTP server in a goroutine
	go func() {
		log.Printf("Starting Pulse-Check server on port :%s", port)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Block until a signal is received
	<-stop
	log.Println("Shutting down gracefully...")

	// Create a context with a timeout for the server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shut down the server and the monitor
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	mon.Stop()
	log.Println("Shutdown complete. Exiting.")
}
