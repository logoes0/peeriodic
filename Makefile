# Peeriodic - Real-Time Collaborative Text Editor
# Makefile for development and deployment

.PHONY: help install build test clean run-be run-fe dev

# Default target
help:
	@echo "Peeriodic - Real-Time Collaborative Text Editor"
	@echo ""
	@echo "Available commands:"
	@echo "  install    - Install all dependencies (backend and frontend)"
	@echo "  build      - Build both backend and frontend"
	@echo "  test       - Run tests for both backend and frontend"
	@echo "  clean      - Clean build artifacts"
	@echo "  run-be     - Start backend server"
	@echo "  run-fe     - Start frontend development server"
	@echo "  dev        - Start both backend and frontend in development mode"
	@echo "  mod        - Tidy Go modules"
	@echo "  docker     - Build and run with Docker"

# Install dependencies
install:
	@echo "Installing backend dependencies..."
	cd backend && go mod tidy
	@echo "Installing frontend dependencies..."
	cd frontend/client && npm install

# Build both backend and frontend
build:
	@echo "Building backend..."
	cd backend && go build -o bin/server main.go
	@echo "Building frontend..."
	cd frontend/client && npm run build

# Run tests
test:
	@echo "Running backend tests..."
	cd backend && go test ./...
	@echo "Running frontend tests..."
	cd frontend/client && npm test

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf backend/bin
	rm -rf frontend/client/build
	rm -rf frontend/client/node_modules

# Start backend server
run-be:
	@echo "Starting backend server..."
	cd backend && go run main.go

# Start frontend development server
run-fe:
	@echo "Starting frontend development server..."
	cd frontend/client && npm start

# Start both backend and frontend in development mode
dev:
	@echo "Starting development environment..."
	@make run-be & make run-fe

# Tidy Go modules
mod:
	@echo "Tidying Go modules..."
	cd backend && go mod tidy

# Docker commands
docker-build:
	@echo "Building Docker images..."
	docker-compose build

docker-up:
	@echo "Starting services with Docker..."
	docker-compose up -d

docker-down:
	@echo "Stopping Docker services..."
	docker-compose down

docker: docker-build docker-up