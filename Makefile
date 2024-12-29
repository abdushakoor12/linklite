.PHONY: all build run clean generate db-create db-drop

# Default target
all: generate build run

# Generate templ templates
generate:
	@echo "Generating templates..."
	@$(HOME)/go/bin/templ generate

# Build the application
build:
	@echo "Building..."
	@go build -o bin/linklite

# Run the application
run:
	@echo "Running LinkLite..."
	@go run main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go install github.com/a-h/templ/cmd/templ@latest

# Development mode: generate templates and run
dev: generate run
