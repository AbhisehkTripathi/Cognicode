.PHONY: build run test clean deps

# Build the application
build:
	go build -o bin/code-doc-tool cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Run with live reload (requires air)
dev:
	air

# Test the application
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf uploads/
	rm -rf output/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Install air for live reload
install-air:
	go install github.com/cosmtrek/air@latest

# Setup directories
setup:
	mkdir -p uploads output web/static web/templates
