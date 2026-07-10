BINARY := go-mod-analyzer
BUILD_DIR := bin

# main
MAIN_PKG := ./cmd/go-mod-analyzer

# Версия Go и флаги линтера
GO := go
LDFLAGS := -s -w

.PHONY: all build run test lint clean fmt vet docker coverage help

all: build

build: ## Build the binary
	@mkdir -p $(BUILD_DIR)
	$(GO) build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY) $(MAIN_PKG)

run: ## Run the program (passes current dir as argument)
	$(GO) run $(MAIN_PKG) .

test: ## Remove build artifacts and coverage files
	$(GO) test -race -cover ./...

fmt: ## Format source code
	$(GO) fmt ./...

vet: ## Run go vet
	$(GO) vet ./...

lint: ## Run golangci-lint (ensure it is installed)
	golangci-lint run ./...

coverage: ## Run tests with race detector and coverage
	$(GO) test -coverprofile=coverage.txt ./...
	$(GO) tool cover -html=coverage.txt -o coverage.html

clean: ## Remove build artifacts and coverage files
	rm -rf $(BUILD_DIR) coverage.txt coverage.html

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'
