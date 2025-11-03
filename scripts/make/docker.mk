# Docker and Docker Compose commands

DOCKER_COMPOSE := docker compose
DOCKER_COMPOSE_FILE := deployments/docker-compose.yml

.PHONY: docker-build
docker-build: ## Build all Docker images
	@echo "Building Docker images..."
	@cd deployments && $(DOCKER_COMPOSE) build
	@echo "Docker images built successfully ✓"

.PHONY: docker-up
docker-up: ## Start all services with docker-compose
	@echo "Starting all services..."
	@cd deployments && $(DOCKER_COMPOSE) up -d
	@echo "All services started ✓"
	@echo ""
	@echo "🌐 Gateway API: http://localhost:8080"
	@echo "📚 API Documentation (Swagger): http://localhost:8081"
	@echo ""
	@echo "Individual gRPC services (for debugging):"
	@echo "  - Auth Service: localhost:9001"
	@echo "  - Users Service: localhost:9002"
	@echo "  - Chat Service: localhost:9003"
	@echo "  - Social Service: localhost:9004"
	@echo "  - Notifications Service: localhost:9005"

.PHONY: docker-down
docker-down: ## Stop all services
	@echo "Stopping all services..."
	@cd deployments && $(DOCKER_COMPOSE) down -v
	@echo "All services stopped ✓"

.PHONY: docker-restart
docker-restart: docker-down docker-up ## Restart all services

.PHONY: docker-rebuild
docker-rebuild: docker-down docker-build docker-up ## Clean rebuild and restart all services

