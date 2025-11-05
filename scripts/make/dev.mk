# Development commands

.PHONY: lint
lint: ## Run golangci-lint on all services
	@echo "Running golangci-lint..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd $$service && golangci-lint run ./... && cd .. || exit 1; \
	done
	@echo "Linting completed ✓"

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running golangci-lint with auto-fix..."
	@for service in $(SERVICES); do \
		echo "Linting $$service..."; \
		cd $$service && golangci-lint run --fix ./... && cd .. || exit 1; \
	done
	@echo "Linting with auto-fix completed ✓"

.PHONY: test
test: ## Run tests for all services
	@echo "Running tests..."
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		go test -v -race -cover ./$$service/...; \
	done
	@echo "Tests completed ✓"

.PHONY: test-integration
test-integration: ## Run integration tests (requires running services)
	@echo "Running integration tests..."
	@./integration_tests/validation/run_tests.sh
	@echo "Integration tests completed ✓"

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@for service in $(SERVICES); do \
		echo "Testing $$service..."; \
		go test -v -race -coverprofile=$$service/coverage.out -covermode=atomic ./$$service/...; \
	done
	@echo "Tests with coverage completed ✓"

.PHONY: build
build: ## Build all services
	@echo "Building all services..."
	@for service in $(SERVICES); do \
		echo "Building $$service..."; \
		cd $$service && go build -o bin/$$service ./cmd/main.go && cd .. || exit 1; \
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
		cd $$service && go mod tidy && cd .. || exit 1; \
	done
	@echo "Go mod tidy completed ✓"

.PHONY: mod-download
mod-download: ## Download dependencies for all services
	@echo "Downloading dependencies..."
	@for service in $(SERVICES); do \
		echo "Downloading dependencies for $$service..."; \
		cd $$service && go mod download && cd .. || exit 1; \
	done
	@echo "Dependencies downloaded ✓"

