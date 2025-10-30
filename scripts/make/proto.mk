# Protocol Buffers generation and management

# Service directories
SERVICES := auth users social chat notifications

.PHONY: proto-gen
proto-gen: ## Generate Go code from proto files for all services
	@echo "Generating proto files for all services..."
	@for service in $(SERVICES); do \
		echo "Generating $$service..."; \
		buf generate --template $$service/buf.gen.yaml; \
	done
	@echo "Proto generation completed ✓"

.PHONY: proto-gen-auth
proto-gen-auth: ## Generate proto files for auth service
	@echo "Generating proto files for auth service..."
	@buf generate --template auth/buf.gen.yaml
	@echo "Auth proto generation completed ✓"

.PHONY: proto-gen-users
proto-gen-users: ## Generate proto files for users service
	@echo "Generating proto files for users service..."
	@buf generate --template users/buf.gen.yaml
	@echo "Users proto generation completed ✓"

.PHONY: proto-gen-social
proto-gen-social: ## Generate proto files for social service
	@echo "Generating proto files for social service..."
	@buf generate --template social/buf.gen.yaml
	@echo "Social proto generation completed ✓"

.PHONY: proto-gen-chat
proto-gen-chat: ## Generate proto files for chat service
	@echo "Generating proto files for chat service..."
	@buf generate --template chat/buf.gen.yaml
	@echo "Chat proto generation completed ✓"

.PHONY: proto-gen-notifications
proto-gen-notifications: ## Generate proto files for notifications service
	@echo "Generating proto files for notifications service..."
	@buf generate --template notifications/buf.gen.yaml
	@echo "Notifications proto generation completed ✓"

.PHONY: proto-lint
proto-lint: ## Lint all proto files
	@echo "Linting proto files..."
	@buf lint
	@echo "Proto linting completed ✓"

.PHONY: proto-format
proto-format: ## Format all proto files
	@echo "Formatting proto files..."
	@buf format -w
	@echo "Proto formatting completed ✓"

.PHONY: proto-breaking
proto-breaking: ## Check for breaking changes in proto files
	@echo "Checking for breaking changes..."
	@buf breaking --against '.git#branch=main'
	@echo "Breaking change check completed ✓"

.PHONY: proto-clean
proto-clean: ## Clean generated proto files
	@echo "Cleaning generated proto files..."
	@for service in $(SERVICES); do \
		echo "Cleaning $$service/pkg/..."; \
		rm -rf $$service/pkg/*; \
		touch $$service/pkg/.gitkeep; \
	done
	@echo "Proto clean completed ✓"

