.PHONY: build test run docker-build docker-run setup

# Go related variables
BINARY_NAME=search-service
WORKER_NAME=search-worker



build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME) cmd/api/main.go
	go build -o bin/$(WORKER_NAME) cmd/worker/main.go

test:
	@echo "Running tests..."
	go test -v ./...

run:
	@echo "Running API service..."
	./bin/$(BINARY_NAME)

run-worker:
	@echo "Running worker..."
	./bin/$(WORKER_NAME)



docker-run:
	@echo "Running Docker container..."
	docker-compose up

setup:
	@echo "Setting up development environment..."
	go mod download
	go mod tidy
	cp config/config.example.yaml config/config.yaml

clean:
	@echo "Cleaning..."
	rm -rf bin/
	go clean 