.PHONY: build clean run test help

# Binary name
BINARY_NAME=corepower

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Use these for running with credentials
USERNAME?=""
PASSWORD?=""

# Default target
.DEFAULT_GOAL := help

help: ## Display available commands
	@echo "CorePower Reservation System Commands:"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(GOBIN)/$(BINARY_NAME) .

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(GOBIN)
	@go clean

run: build ## Run the application (use: make run USERNAME="email" PASSWORD="pass")
	@if [ $(USERNAME) = "" ] || [ $(PASSWORD) = "" ]; then \
		echo "Error: USERNAME and PASSWORD are required"; \
		echo "Usage: make run USERNAME=\"your.email@example.com\" PASSWORD=\"your_password\""; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME)..."
	@$(GOBIN)/$(BINARY_NAME) -username=$(USERNAME) -password=$(PASSWORD)

lint: ## Run linters
	@echo "Running linters..."
	@go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
	else \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
