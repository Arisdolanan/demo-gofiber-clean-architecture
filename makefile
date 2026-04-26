# GoFiber Clean Architecture Project Makefile

.PHONY: help info requirements clean build run worker dev dev-worker test test-coverage fmt lint docs migrate-up migrate-down migrate-status migrate-up-dev install-migrate

# Default target
info:
	$(info ------------------------------------------)
	$(info -           GoFiber Project          -)
	$(info ------------------------------------------)
	$(info Available commands:)
	$(info - build:          Build web and worker binaries)
	$(info - run:            Run web application)
	$(info - worker:         Run worker application)
	$(info - dev:            Run in development mode with hot-reload)
	$(info - test:           Run all tests)
	$(info - test-coverage:  Run tests with coverage report)
	$(info - fmt:            Format Go code)
	$(info - lint:           Lint Go code)
	$(info - docs:           Generate Swagger documentation)
	$(info - requirements:   Tidy Go modules)
	$(info - clean:          Clean build artifacts)
	$(info - migrate-up:     Apply database migrations (no seed))
	$(info - migrate-down:   Rollback database migrations)
	$(info - migrate-up-dev: Apply migrations with seed data for development)
	$(info - migrate-create: Create new migration file)
	$(info - migrate-status: Check migration status)
	$(info )
	$(info Usage: make <command>)

# Development and Build Commands
requirements:
	@echo "Tidying Go modules..."
	go mod tidy
	@echo "Done: Dependencies updated!"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/ coverage.out coverage.html
	go clean -cache
	@echo "Done: Clean completed!"

build:
	@echo "Building GoFiber applications..."
	@mkdir -p bin
	go build -o bin/web ./cmd/web
	go build -o bin/worker ./cmd/worker
	@echo "Done: Build completed! Binaries: bin/web, bin/worker"

run:
	@echo "Starting web application..."
	go run ./cmd/web/main.go

worker:
	@echo "Starting worker application..."
	go run ./cmd/worker/main.go

dev:
	@echo "Starting development mode with hot-reload..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air

dev-worker:
	@echo "Starting worker in development mode with hot-reload..."
	@which air > /dev/null || (echo "Installing air..." && go install github.com/cosmtrek/air@latest)
	air -c .air-worker.toml

production:
	@echo "Building for production..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/web -ldflags="-s -w" ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/worker -ldflags="-s -w" ./cmd/worker
	@echo "Done: Build completed! Binaries: bin/web, bin/worker"

# Testing Commands
test:
	@echo "Running all tests..."
	go test ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Done: Coverage report generated: coverage.html"

# Code Quality Commands
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "Done: Code formatted!"

lint:
	@echo "Linting Go code..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run
	@echo "Done: Linting completed!"

# Documentation
docs:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/web/main.go -o api/docs
	@echo "Done: Documentation generated!"

# Database Migration Commands
install-migrate:
	@echo "Installing golang-migrate tool..."
	@which migrate > /dev/null || (echo "Installing migrate..." && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
	@echo "Done: golang-migrate installed!"

check-migrate:
	@which migrate > /dev/null || (echo "Error: migrate tool not found. Run 'make install-migrate' first." && exit 1)

migrate-up: check-migrate
	@echo "Applying database migrations (without seed data)..."
	@migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gofiber?sslmode=disable" up 11
	@echo "Done: Migrations applied successfully (production-ready, no seed data)!"

migrate-down: check-migrate
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gofiber?sslmode=disable" down 1
	@echo "Done: Migration rolled back!"

# Apply migrations with seed data (for development)
migrate-up-dev: check-migrate
	@echo "Applying ALL database migrations including seed data for development..."
	@migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gofiber?sslmode=disable" up
	@echo "Done: All migrations applied successfully (including seed data)!"

migrate-create: check-migrate
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
	@echo "Done: Migration file created!"

migrate-status: check-migrate
	@echo "Checking migration status..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gofiber?sslmode=disable" version

# Docker Commands
docker-build:
	@echo "Building Docker image..."
	docker build -t demo-gofiber-app .
	@echo "Done: Docker image built!"

docker-run:
	@echo "Running Docker container..."
	docker run -d -p 3000:3000 demo-gofiber-app
	@echo "Done: Docker container running!"

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d postgres pgadmin redis rabbitmq zookeeper kafka kafka-ui portainer
	@echo "Done: Docker services started!"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "Done: Docker services stopped!"

docker-restart:
	@echo "Building Docker image app..."
	docker-compose build --no-cache app worker && docker-compose restart app worker
	@echo "Done: Docker services restarted!"

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

# Help command
help: info
