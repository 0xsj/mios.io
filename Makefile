.PHONY: build run test clean migrate-up migrate-down migrate-create sqlc docker-up docker-down lint mock help

BINARY_NAME=mios.io-app
VERSION=0.1.0
BUILD_DIR=./bin
ENV_FILE=./dev.env
DB_DSN=postgres://devuser:devpass@localhost:5432/devdb?sslmode=disable
MIGRATION_PATH=./db/migration

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

## sqlc: Generate database code using SQLC
sqlc:
	@echo "Generating SQLC code..."
	@if command -v sqlc > /dev/null; then \
		sqlc generate; \
	else \
		echo "sqlc not found. Use 'make install-tools' to install it."; \
		exit 1; \
	fi

## docker-up: Start docker containers
docker-up:
	@echo "Starting docker containers..."
	@docker-compose up -d

## docker-down: Stop docker containers
docker-down:
	@echo "Stopping docker containers..."
	@docker-compose down

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

## help: Display help information
help:
	@echo "Available targets:"
	@grep -E '^## [a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/## //' | sort

## help: format code
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