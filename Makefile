# Variables
PROJECT_NAME=task_tracker
DOCKER_COMPOSE=docker-compose -f docker-compose.yml

# Default target
all: build

# Build the Docker images
build:
	@echo "Building Docker images..."
	$(DOCKER_COMPOSE) build

# Start the application
up:
	@echo "Starting the application..."
	$(DOCKER_COMPOSE) up --build


# Start the application
upd:
	@echo "Starting the application..."
	$(DOCKER_COMPOSE) up -d --build


# Stop the application
down:
	@echo "Stopping the application..."
	$(DOCKER_COMPOSE) down

# Restart the application
restart: down up

# View logs for all services
logs:
	@echo "Viewing logs..."
	$(DOCKER_COMPOSE) logs -f

# View logs for a specific service (e.g., make log service=app)
log:
	@echo "Viewing logs for $(service)..."
	$(DOCKER_COMPOSE) logs -f $(service)


# Run lint
lint:
	@echo "Running lint..."
	golangci-lint run ./...

# Run formatters
fmt:
	@echo "Formatting..."
	golangci-lint fmt ./... -v

# Run tests
test: lint
	@echo "Running tests..."
	go test ./...

# Clean up Docker resources (containers, volumes, networks)
clean:
	@echo "Cleaning up Docker resources..."
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

# Remove all Docker images
clean-images:
	@echo "Removing Docker images..."
	docker rmi -f $$(docker images -q)

# Open a shell in the app container
shell:
	@echo "Opening a shell in the app container..."
	docker exec -it $(PROJECT_NAME)_app sh

# Open a shell in the db container
shell-db:
	@echo "Opening a shell in the db container..."
	docker exec -it $(PROJECT_NAME)_db sh

# Open a shell in the loki container
shell-loki:
	@echo "Opening a shell in the loki container..."
	docker exec -it $(PROJECT_NAME)_loki sh

# Open a shell in the promtail container
shell-promtail:
	@echo "Opening a shell in the promtail container..."
	docker exec -it $(PROJECT_NAME)_promtail sh

# Open a shell in the grafana container
shell-grafana:
	@echo "Opening a shell in the grafana container..."
	docker exec -it $(PROJECT_NAME)_grafana sh

# Help command to list all available commands
help:
	@echo "Available commands:"
	@echo "  make build        - Build Docker images"
	@echo "  make up           - Start the application"
	@echo "  make down         - Stop the application"
	@echo "  make restart      - Restart the application"
	@echo "  make logs         - View logs for all services"
	@echo "  make log service=<service> - View logs for a specific service"
	@echo "  make lint         - Run lint"
	@echo "  make fmt          - Format go source code"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean up Docker resources"
	@echo "  make clean-images - Remove all Docker images"
	@echo "  make shell        - Open a shell in the app container"
	@echo "  make shell-db     - Open a shell in the db container"
	@echo "  make shell-loki   - Open a shell in the loki container"
	@echo "  make shell-promtail - Open a shell in the promtail container"
	@echo "  make shell-grafana - Open a shell in the grafana container"
	@echo "  make help         - Show this help message"

# Default target
.PHONY: all build up down restart logs log test clean clean-images shell shell-db shell-loki shell-promtail shell-grafana help