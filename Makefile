.PHONY: test

# Define Docker Compose command
DOCKER_COMPOSE := docker-compose

# Docker Compose file
COMPOSE_FILE := ./deployment/docker-compose.yml

# Targets

# Build Docker containers
build:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) build

# Run Docker containers
up:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) up -d

# Stop Docker containers
down:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down

# Restart Docker containers
restart: down up

# View logs of Docker containers
logs:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) logs -f

# Remove Docker containers and associated volumes
clean:
	$(DOCKER_COMPOSE) -f $(COMPOSE_FILE) down -v

# Run Test
test:
	go test -v -count=1 ./...

# Benchmark Test
benchmark: up test clean