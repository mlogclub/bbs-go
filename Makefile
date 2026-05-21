APP := bbs-go
MAIN := ./main.go
WEB_DIR := web
SPA_INDEX := $(WEB_DIR)/build/spa/index.html

GO ?= go
PNPM ?= pnpm
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

.DEFAULT_GOAL := help

.PHONY: all
all: build

.PHONY: help
help:
	@echo "Usage:"
	@echo "  make build          Build SPA and embed it into the Go binary"
	@echo "  make build-go       Build Go binary, building SPA first only when missing"
	@echo "  make build-linux    Build linux/amd64 binary with embedded SPA"
	@echo "  make release        Build release binaries for linux, macOS, and Windows"
	@echo "  make run            Build SPA, then run the Go server"
	@echo "  make run-go         Run the Go server, building SPA first only when missing"
	@echo "  make test           Run Go tests"
	@echo "  make check          Run Go tests and web checks"
	@echo "  make clean          Remove Go binaries"
	@echo "  make clean-web      Remove web build output"
	@echo "  make web-build-spa  Build SPA output for embedding"
	@echo "  make web-build-ssr  Build SSR output"
	@echo "  make web-dev        Start the web dev server"

.PHONY: build
build: web-build-spa
	@$(MAKE) build-go

.PHONY: build-go
build-go: ensure-spa
	@echo "Building $(APP)..."
	@$(GO) build -v -o $(APP) $(MAIN)

.PHONY: build-linux
build-linux: web-build-spa
	@echo "Building $(APP)-linux-amd64..."
	@GOOS=linux GOARCH=amd64 $(GO) build -v -o $(APP)-linux-amd64 $(MAIN)

.PHONY: release
release: web-build-spa
	@echo "Building release binaries..."
	@GOOS=linux GOARCH=amd64 $(GO) build -v -o $(APP)-linux-amd64 $(MAIN)
	@GOOS=darwin GOARCH=amd64 $(GO) build -v -o $(APP)-macos-amd64 $(MAIN)
	@GOOS=darwin GOARCH=arm64 $(GO) build -v -o $(APP)-macos-arm64 $(MAIN)
	@GOOS=windows GOARCH=amd64 $(GO) build -v -o $(APP)-windows-amd64.exe $(MAIN)

.PHONY: run
run: web-build-spa
	@$(GO) run $(MAIN)

.PHONY: run-go
run-go: ensure-spa
	@$(GO) run $(MAIN)

.PHONY: test
test: ensure-spa
	@echo "Running Go tests..."
	@$(GO) test ./...

.PHONY: check
check: test web-typecheck web-lint

.PHONY: clean
clean:
	@echo "Cleaning Go binaries..."
	@rm -f $(APP) $(APP)-linux-* $(APP)-macos-* $(APP)-windows-*.exe

.PHONY: clean-web
clean-web:
	@echo "Cleaning web build output..."
	@rm -rf $(WEB_DIR)/build

.PHONY: web-install
web-install:
	@cd $(WEB_DIR) && $(PNPM) install --frozen-lockfile

.PHONY: web-dev
web-dev:
	@cd $(WEB_DIR) && $(PNPM) dev

.PHONY: web-build-spa
web-build-spa:
	@echo "Building SPA..."
	@cd $(WEB_DIR) && $(PNPM) build:spa

.PHONY: ensure-spa
ensure-spa:
	@if [ ! -f "$(SPA_INDEX)" ]; then \
		echo "SPA build output is missing; building SPA..."; \
		$(MAKE) web-build-spa; \
	fi

.PHONY: build-spa
build-spa: web-build-spa

.PHONY: web-build-ssr
web-build-ssr:
	@cd $(WEB_DIR) && $(PNPM) build:ssr

.PHONY: web-typecheck
web-typecheck:
	@cd $(WEB_DIR) && $(PNPM) typecheck

.PHONY: web-lint
web-lint:
	@cd $(WEB_DIR) && $(PNPM) lint

.PHONY: generator
generator:
	@$(GO) run cmd/generator/generator.go

.PHONY: generate-permissions
generate-permissions:
	@$(GO) run ./cmd/generator/permissions
