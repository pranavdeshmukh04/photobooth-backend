.PHONY: run build test test-coverage lint fmt clean docker-build docker-run help

# Variables
APP_NAME=photobooth-server
BINARY_PATH=./cmd/server
DOCKER_IMAGE=photobooth-backend

# Run the application
run:
	go run $(BINARY_PATH)/main.go

# Build binary
build:
	go build -ldflags="-s -w" -o $(APP_NAME) $(BINARY_PATH)

# Run tests
test:
	go test ./...

# Run tests with coverage
test-coverage:
	go test -cover ./...

# Run linter
lint:
	golangci-lint run ./...

# Format code
fmt:
	go fmt ./...
	goimports -w .

# Clean build artifacts
clean:
	rm -f $(APP_NAME)
	go clean ./...

# Build Docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Install dependencies
deps:
	go mod download
	go mod tidy

# Show help
help:
	@echo "Available commands:"
	@echo "  make run              - Run the application"
	@echo "  make build            - Build binary"
	@echo "  make test             - Run tests"
	@echo "  make test-coverage    - Run tests with coverage"
	@echo "  make lint             - Run linter"
	@echo "  make fmt              - Format code"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make docker-build     - Build Docker image"
	@echo "  make docker-run       - Run Docker container"
	@echo "  make deps             - Install dependencies"
	@echo "  make help             - Show this help"
