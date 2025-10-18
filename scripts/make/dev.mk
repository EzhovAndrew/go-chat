# Development commands

.PHONY: lint
lint: ## Run golangci-lint on all services
	@echo "Running golangci-lint..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd $$service && golangci-lint run ./... && cd ..; \
	done
	@echo "Linting completed ✓"

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running golangci-lint with auto-fix..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd $$service && golangci-lint run --fix ./... && cd ..; \
	done
	@echo "Linting with auto-fix completed ✓"

.PHONY: test
test: ## Run tests for all services
	@echo "Running tests..."
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd $$service && go test -v -race -cover ./... && cd ..; \
	done
	@echo "Tests completed ✓"

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		cd $$service && go test -v -race -coverprofile=coverage.out -covermode=atomic ./... && cd ..; \
	done
	@echo "Tests with coverage completed ✓"

.PHONY: build
build: ## Build all services
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && go build -o bin/$$service ./cmd/... && cd ..; \
	done
	@echo "Build completed ✓"

.PHONY: clean
clean: proto-clean ## Clean all build artifacts
	@echo "Cleaning build artifacts..."
	@for service in $(SERVICES); do \
		echo "Cleaning $$service..."; \
		rm -rf $$service/bin $$service/coverage.out; \
	done
	@echo "Clean completed ✓"

.PHONY: mod-tidy
mod-tidy: ## Run go mod tidy for all services
	@echo "Running go mod tidy..."
	@for service in $(SERVICES); do \
		echo "Tidying $$service..."; \
		cd $$service && go mod tidy && cd ..; \
	done
	@echo "Go mod tidy completed ✓"

.PHONY: mod-download
mod-download: ## Download dependencies for all services
	@echo "Downloading dependencies..."
	@for service in $(SERVICES); do \
		echo "Downloading dependencies for $$service..."; \
		cd $$service && go mod download && cd ..; \
	done
	@echo "Dependencies downloaded ✓"

