# Build the application
all: build

build: build-spa
	@echo "Building..."
	@go build -v -o bbs-go main.go

buildlinux: build-spa
	@echo "Building..."
	@GOOS=linux GOARCH=amd64 go build

build-spa:
	@echo "Building SPA..."
	@cd web && pnpm build:spa

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

generate-permissions:
	@go run ./cmd/generator/permissions

.PHONY: all build buildlinux build-spa run test clean generator generate-permissions
