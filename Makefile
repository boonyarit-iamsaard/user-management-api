.PHONY: run build clean

# Run the application
run:
	go run cmd/main.go

# Build the application
build:
	go build -o bin/user-management-api cmd/main.go

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run tests
test:
	go test ./...

# Download dependencies
deps:
	go mod download
	go mod tidy

# Lint code
lint:
	golangci-lint run
