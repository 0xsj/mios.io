.PHONY: build run test clean migrate-up migrate-down migrate-create sqlc docker-up docker-down docker-build docker-logs docker-ps docker-exec docker-restart lint mock help format openapi serve-docs

BINARY_NAME=gin-sqlc-app
VERSION=0.1.0
BUILD_DIR=./bin
ENV_FILE=./dev.env
DB_DSN=postgres://devuser:devpass@localhost:5432/devdb?sslmode=disable
MIGRATION_PATH=./db/migration
DOCKER_COMPOSE_FILE=docker-compose.yml
CONTAINER_NAME=gin-sqlc-app

# Default target
all: clean build

## build: Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) -v ./main.go

## run: Build and run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

## dev: Run the application with hot-reload using air
dev:
	@if command -v air > /dev/null; then \
		echo "Running with air for hot-reload..."; \
		air; \
	else \
		echo "air not found. Installing air..."; \
		go install github.com/air-verse/air@latest; \
		air; \
	fi

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## test-coverage: Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out

## clean: Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

## migrate-up: Run all database migrations
migrate-up:
	@echo "Running migrations up..."
	@if command -v migrate > /dev/null; then \
		migrate -path $(MIGRATION_PATH) -database "$(DB_DSN)" -verbose up; \
	else \
		echo "migrate tool not found. Use 'make install-tools' to install it."; \
		exit 1; \
	fi

## migrate-down: Revert last database migration
migrate-down:
	@echo "Running migrations down..."
	@if command -v migrate > /dev/null; then \
		migrate -path $(MIGRATION_PATH) -database "$(DB_DSN)" -verbose down 1; \
	else \
		echo "migrate tool not found. Use 'make install-tools' to install it."; \
		exit 1; \
	fi

## migrate-down-all: Revert all database migrations
migrate-down-all:
	@echo "Reverting all migrations..."
	@if command -v migrate > /dev/null; then \
		migrate -path $(MIGRATION_PATH) -database "$(DB_DSN)" -verbose down; \
	else \
		echo "migrate tool not found. Use 'make install-tools' to install it."; \
		exit 1; \
	fi

## migrate-create name=migration_name: Create a new migration file
migrate-create:
	@echo "Creating migration $(name)..."
	@if command -v migrate > /dev/null; then \
		migrate create -ext sql -dir $(MIGRATION_PATH) -seq $(name); \
	else \
		echo "migrate tool not found. Use 'make install-tools' to install it."; \
		exit 1; \
	fi

## docker-migrate-up: Run database migrations inside the Docker environment
docker-migrate-up:
	@echo "Running migrations up in Docker environment..."
	@docker-compose exec postgres sh -c "PGPASSWORD=devpass psql -U devuser -d devdb -h localhost -f /migrations/up.sql"

## sqlc: Generate database code using SQLC
sqlc:
	@echo "Generating SQLC code..."
	@if command -v sqlc > /dev/null; then \
		sqlc generate; \
	else \
		echo "sqlc not found. Use 'make install-tools' to install it."; \
		exit 1; \
	fi

## docker-build: Build Docker images
docker-build:
	@echo "Building Docker images..."
	@docker-compose build

## docker-up: Start docker containers
docker-up:
	@echo "Starting docker containers..."
	@docker-compose up -d

## docker-down: Stop docker containers
docker-down:
	@echo "Stopping docker containers..."
	@docker-compose down

## docker-down-volumes: Stop docker containers and remove volumes
docker-down-volumes:
	@echo "Stopping docker containers and removing volumes..."
	@docker-compose down -v

## docker-logs: Show logs from all containers or specific container (c=container_name)
docker-logs:
	@if [ -z "$(c)" ]; then \
		echo "Showing logs from all containers..."; \
		docker-compose logs -f; \
	else \
		echo "Showing logs from $(c)..."; \
		docker-compose logs -f $(c); \
	fi

## docker-ps: List running containers
docker-ps:
	@echo "Listing containers..."
	@docker-compose ps

## docker-exec: Execute command in container (c=container_name, default=app)
docker-exec:
	@container=$${c:-$(CONTAINER_NAME)}; \
	echo "Executing shell in $$container..."; \
	docker-compose exec $$container sh

## docker-restart: Restart specific container or all (c=container_name)
docker-restart:
	@if [ -z "$(c)" ]; then \
		echo "Restarting all containers..."; \
		docker-compose restart; \
	else \
		echo "Restarting $(c)..."; \
		docker-compose restart $(c); \
	fi

## docker-prune: Remove unused Docker data (images, containers, volumes)
docker-prune:
	@echo "Pruning unused Docker data..."
	@docker system prune -f

## lint: Lint the code
lint:
	@echo "Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found. Use 'make install-tools' to install it."; \
		exit 1; \
	fi

## mock: Generate mocks for testing
mock:
	@echo "Generating mocks..."
	@if command -v mockgen > /dev/null; then \
		go generate ./...; \
	else \
		echo "mockgen not found. Use 'make install-tools' to install it."; \
		exit 1; \
	fi

## install-tools: Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/air-verse/air@latest

## format: Format code
format:
	@echo "Formatting code..."
	@if command -v gofumpt > /dev/null; then \
		gofumpt -l -w .; \
	else \
		go fmt ./...; \
	fi

## openapi: Generate OpenAPI documentation
openapi:
	@echo "Generating OpenAPI documentation..."
	@go run cmd/openapi/main.go

## serve-docs: Serve OpenAPI documentation
serve-docs:
	@echo "Serving OpenAPI documentation..."
	@docker run -p 8085:8080 -e SWAGGER_JSON=/docs/openapi.json -v $(PWD)/docs:/docs swaggerapi/swagger-ui

## docker-compose: Create docker-compose.yml file
docker-compose:
	@echo "Creating docker-compose.yml file..."
	@cat > docker-compose.yml << EOF
services:
  postgres:
    image: postgres:14
    container_name: gin-sqlc-postgres
    environment:
      POSTGRES_DB: devdb
      POSTGRES_USER: devuser
      POSTGRES_PASSWORD: devpass
    ports:
      - "5432:5432"
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./db/migration:/migrations
    networks:
      - app-network
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U devuser -d devdb"]
      interval: 10s
      timeout: 5s
      retries: 5
  
  redis:
    image: redis:alpine
    container_name: gin-sqlc-redis
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
  
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: gin-sqlc-app
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    environment:
      # Existing environment variables
      - ENVIRONMENT=development
      - HOST=0.0.0.0
      - PORT=8080
      - DB_USERNAME=devuser
      - DB_PASSWORD=devpass
      - DB_HOSTNAME=postgres
      - DB_PORT=5432
      - DB_NAME=devdb
      - JWT_SECRET=askimaskimaskimasecurelongersecret1234
      - TOKEN_HOUR_LIFESPAN=24
      - API_SECRET=jagiya
      - VERSION=1
      - GIN_MODE=release
      # Redis environment variables
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - REDIS_PASSWORD=
      - REDIS_DB=0
    networks:
      - app-network
    restart: on-failure

networks:
  app-network:
    driver: bridge

volumes:
  postgres-data:
  redis-data:
EOF
	@echo "docker-compose.yml created successfully"

## docker-file: Create Dockerfile
docker-file:
	@echo "Creating Dockerfile..."
	@cat > Dockerfile << EOF
# Build stage
FROM golang:1.24-alpine AS builder
# Set working directory
WORKDIR /app
# Install build dependencies
RUN apk add --no-cache gcc musl-dev make git
# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
# Download all dependencies
RUN go mod download
# Copy the source code
COPY . .
# Build the application
RUN make build
# Final stage
FROM alpine:latest
# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata
# Set working directory
WORKDIR /app
# Copy the binary from the builder stage
COPY --from=builder /app/bin/gin-sqlc-app .
# Copy the environment file
COPY dev.env .
# Expose the application port
EXPOSE 8080
# Run the application
CMD ["./gin-sqlc-app"]
EOF
	@echo "Dockerfile created successfully"

## docker-all: Quick setup for Docker development (create files, build, and start)
docker-all: docker-compose docker-file docker-build docker-up
	@echo "Docker environment set up and running. Use 'make docker-logs' to see container logs."

## help: Display help information
help:
	@echo "Available targets:"
	@grep -E '^## [a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/## //' | sort