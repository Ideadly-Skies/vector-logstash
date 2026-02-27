.PHONY: help build build-sender build-sender-advanced run run-advanced test clean fmt lint deps vector-start

# Variables
BINARY_NAME=sender
BINARY_ADVANCED=sender-advanced
VECTOR_CONFIG=configs/vector.yaml
GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
GO_FMT=$(GO_CMD) fmt
GO_VET=$(GO_CMD) vet

help: ## Show this help message
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download dependencies
	$(GO_GET) -u ./...
	$(GO_CMD) mod tidy
	$(GO_CMD) mod download

build: build-sender build-sender-advanced ## Build all binaries

build-sender: ## Build basic sender
	$(GO_BUILD) -o bin/$(BINARY_NAME) ./cmd/sender

build-sender-advanced: ## Build advanced sender
	$(GO_BUILD) -o bin/$(BINARY_ADVANCED) ./cmd/sender-advanced

run: build-sender ## Run basic sender
	./bin/$(BINARY_NAME)

run-advanced: build-sender-advanced ## Run advanced sender
	./bin/$(BINARY_ADVANCED)

test: ## Run tests
	$(GO_TEST) -v -race -coverprofile=coverage.out ./...

coverage: test ## Generate coverage report
	$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

fmt: ## Format code
	$(GO_FMT) ./...
	@echo "Code formatted"

lint: ## Run linter
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: brew install golangci-lint"; \
		$(GO_VET) ./...; \
	fi

clean: ## Clean build artifacts
	$(GO_CLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf logs/

vector-start: ## Start Vector
	@if ! command -v vector > /dev/null; then \
		echo "Vector not installed. Install with: brew install vector"; \
		exit 1; \
	fi
	vector --config $(VECTOR_CONFIG)

vector-validate: ## Validate Vector configuration
	vector validate $(VECTOR_CONFIG)

install: build ## Install binaries to $GOPATH/bin
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/
	cp bin/$(BINARY_ADVANCED) $(GOPATH)/bin/

all: clean deps fmt lint test build ## Run all checks and build
