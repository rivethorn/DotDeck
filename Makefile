# App metadata
APP_NAME := deck
PKG      := .
BIN_DIR  := bin
GO_FILES := $(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Go build settings
GO       := go
GOFLAGS  :=
LDFLAGS  := -s -w

# Git info (optional, inject into build)
VERSION  := $(shell git describe --tags --always --dirty)
COMMIT   := $(shell git rev-parse --short HEAD)
BUILDTIME:= $(shell date +'%Y-%m-%dT%H:%M:%S%z')

# Default target
all: build

## Build binary
build: $(GO_FILES)
	@echo "🔨 Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS) -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildTime=$(BUILDTIME)'" -o $(BIN_DIR)/$(APP_NAME) $(PKG)
	@echo "✅ Built $(BIN_DIR)/$(APP_NAME)"

## Install binary globally
install: build
	@echo "📦 Installing to $$GOPATH/bin (or Go bin dir)..."
	@$(GO) install $(GOFLAGS) -ldflags "$(LDFLAGS) -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.buildTime=$(BUILDTIME)'" $(PKG)

## Run tests
test:
	@echo "🧪 Running tests..."
	@$(GO) test ./... -v

## Run lint (requires golangci-lint installed)
lint:
	@echo "🔍 Linting..."
	@golangci-lint run

## Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BIN_DIR)

## Cross-compile (example for Linux amd64)
build-linux:
	@echo "🌐 Cross compiling for linux/amd64..."
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(APP_NAME)-linux $(PKG)

build-windows:
	@echo "🌐 Cross compiling for windows/amd64..."
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(APP_NAME)-windows.exe $(PKG)

## Show help
help:
    @echo "Available make targets:"
    @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
        sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

.PHONY: all build install test lint clean build-linux help