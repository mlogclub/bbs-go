# Build the application
all: build

build:
	@echo "Building..."
	@go build -v -o bbs-go main.go

# Run the application
run:
	@go run main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./...

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bbs-go

generator:
	@go run cmd/generator/generator.go

.PHONY: all build run test clean
