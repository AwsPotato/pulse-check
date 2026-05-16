.PHONY: build docker run clean

APP_NAME = pulse-check
CMD_DIR = ./cmd/pulse-check
BUILD_DIR = build/package

build:
	@echo "Building $(APP_NAME) for current OS..."
	go build -o $(APP_NAME) $(CMD_DIR)/main.go

docker:
	@echo "Building distroless Docker image..."
	docker build -t $(APP_NAME):latest -f $(BUILD_DIR)/Dockerfile .

run:
	@echo "Running $(APP_NAME) locally..."
	go run $(CMD_DIR)/main.go

clean:
	@echo "Cleaning up..."
	@if [ -f $(APP_NAME) ]; then rm $(APP_NAME); fi
