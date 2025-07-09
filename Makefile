VERSION ?= 1.0.0
OUTPUT_DIR := dist

.PHONY: build clean test install help

help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build binaries for all platforms
	@./scripts/build.sh $(VERSION)

clean: ## Clean build artifacts
	@rm -rf $(OUTPUT_DIR)
	@echo "Cleaned build artifacts"

test: ## Run tests
	@go test -v ./...

install: ## Install dependencies
	@go mod download
	@go mod tidy

dev: ## Build for current platform only
	@echo "Building for development..."
	@go build -ldflags="-X github.com/snupai/cngt-cli/internal/version.Version=$(VERSION) -X github.com/snupai/cngt-cli/internal/version.GitCommit=$$(git rev-parse HEAD) -X github.com/snupai/cngt-cli/internal/version.BuildTime=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o cngt-cli ./cmd/main.go
	@echo "Built: cngt-cli"

run: ## Run the application
	@go run ./cmd/main.go

fmt: ## Format code
	@go fmt ./...

lint: ## Run linter
	@golangci-lint run

deps: ## Download dependencies
	@go mod download

pre-release: ## Run pre-release validation checks
	@./scripts/pre-release-check.sh

release: ## Create a new release (usage: make release VERSION=1.0.1)
	@if [ -z "$(VERSION)" ]; then echo "Error: VERSION is required. Usage: make release VERSION=1.0.1"; exit 1; fi
	@./scripts/tag-release.sh $(VERSION)