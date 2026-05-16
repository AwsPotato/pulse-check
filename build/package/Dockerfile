# Stage 1: Build the Go binary
FROM golang:1.21-alpine AS builder

# Install git if needed by dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy the entire project source
COPY . .

# Build the statically linked binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o pulse-check ./cmd/pulse-check/main.go

# Stage 2: Create the final production image
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/pulse-check .

# Expose the application port
EXPOSE 8080

# Run as a non-root user for security
USER nonroot:nonroot

# Entry point for the container
ENTRYPOINT ["/pulse-check"]
