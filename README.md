# Pulse-Check

Pulse-Check is a lightweight, concurrent system monitoring agent written in Go. It collects CPU, Memory, and Disk metrics every 5 seconds and serves them via a fast, non-blocking HTTP endpoint.

## Features
- **Concurrent Polling:** Background worker uses Goroutines to collect stats.
- **Thread-Safe Data:** Uses `sync.RWMutex` to ensure fast, safe, non-blocking reads.
- **Minimal Dependencies:** Built heavily on the Go standard library, utilizing `gopsutil` for cross-platform metric gathering.
- **Hardened Production Ready:** Multi-stage Docker build resulting in a secure, distroless image containing only the compiled binary (<20MB).

## Getting Started

### Local Development
To run the agent locally without Docker, simply use:
```bash
make run
```
You can view the metrics by navigating to `http://localhost:8080/metrics`.
To check the server health, use `http://localhost:8080/health`.

### Docker Deployment
To build the hardened production container:
```bash
make docker
```
To run the container:
```bash
docker run -p 8080:8080 pulse-check:latest
```

Alternatively, use Docker Compose:
```bash
docker-compose up -d
```
