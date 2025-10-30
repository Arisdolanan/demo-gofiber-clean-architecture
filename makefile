# GoFiber Clean Architecture Project Makefile

.PHONY: help info requirements clean build run worker dev test test-coverage fmt lint docs migrate-up migrate-down migrate-create migrate-status install-migrate

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
	$(info - migrate-up:     Apply database migrations)
	$(info - migrate-down:   Rollback database migrations)
	$(info - migrate-create: Create new migration file)
	$(info - migrate-status: Check migration status)
	$(info )
	$(info Usage: make <command>)

# Development and Build Commands
requirements:
	@echo "Tidying Go modules..."
	go mod tidy
	@echo "✅ Dependencies updated!"

clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/ coverage.out coverage.html
	go clean -cache
	@echo "✅ Clean completed!"

build:
	@echo "Building GoFiber applications..."
	@mkdir -p bin
	go build -o bin/web ./cmd/web
	go build -o bin/worker ./cmd/worker
	@echo "✅ Build completed! Binaries: bin/web, bin/worker"

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

production:
	@echo "Building for production..."
	@mkdir -p bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/web -ldflags="-s -w" ./cmd/web
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/worker -ldflags="-s -w" ./cmd/worker
	@echo "✅ Build completed! Binaries: bin/web, bin/worker

# Testing Commands
test:
	@echo "Running all tests..."
	go test ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

# Code Quality Commands
fmt:
	@echo "Formatting Go code..."
	go fmt ./...
	@echo "✅ Code formatted!"

lint:
	@echo "Linting Go code..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run
	@echo "✅ Linting completed!"

# Documentation
docs:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/web/main.go -o api/docs
	@echo "✅ Documentation generated!"

# Database Migration Commands
install-migrate:
	@echo "Installing golang-migrate tool..."
	@which migrate > /dev/null || (echo "Installing migrate..." && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest)
	@echo "✅ golang-migrate installed!"

check-migrate:
	@which migrate > /dev/null || (echo "❌ Error: migrate tool not found. Run 'make install-migrate' first." && exit 1)

migrate-up: check-migrate
	@echo "Applying database migrations..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gofiber?sslmode=disable" up
	@echo "✅ Migrations applied!"

migrate-down: check-migrate
	@echo "Rolling back database migrations..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gofiber?sslmode=disable" down 1
	@echo "✅ Migration rolled back!"

migrate-create: check-migrate
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name
	@echo "✅ Migration file created!"

migrate-status: check-migrate
	@echo "Checking migration status..."
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gofiber?sslmode=disable" version

# Docker Commands
docker-build: 
	@echo "Building Docker image..."
	 docker build -t demo-gofiber-app .
	@echo "✅ Docker image built!"

docker-run:
	@echo "Running Docker container..."
	docker run -d -p 3000:3000 demo-gofiber-app
	@echo "✅ Docker container running!"

docker-up:
	@echo "Starting Docker services..."
	docker-compose up -d
	@echo "✅ Docker services started!"

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down
	@echo "✅ Docker services stopped!"

docker-restart:
	@echo "Building Docker image app..."
	docker-compose build --no-cache app worker && docker-compose restart app worker
	@echo "✅ Docker services restart!"

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose logs -f

# Help command
help: info