.PHONY: help install up down logs backend-test frontend-test wire swagger backend-run db-start db-stop env-check

# Load .env file
include .env
export

help:
	@echo "Available commands:"
	@echo "  make env-check      - Check if .env file exists"
	@echo "  make install        - Install all dependencies"
	@echo "  make up             - Start all services with docker-compose"
	@echo "  make down           - Stop all services"
	@echo "  make logs           - Show docker-compose logs"
	@echo "  make db-start       - Start only PostgreSQL"
	@echo "  make db-stop        - Stop PostgreSQL"
	@echo "  make wire           - Generate Wire dependency injection code"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo "  make backend-run    - Run backend locally (requires PostgreSQL)"
	@echo "  make backend-test   - Run backend tests"
	@echo "  make frontend-test  - Run frontend tests"

env-check:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found!"; \
		echo "Please copy .env.example to .env and configure it."; \
		exit 1; \
	fi
	@echo ".env file found ✓"

install:
	@echo "Installing backend dependencies..."
	cd backend && go mod download
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

backend-test:
	cd backend && go test ./... -v -cover

frontend-test:
	cd frontend && npm run test

db-start:
	@echo "Starting PostgreSQL..."
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5
	@echo "PostgreSQL is ready ✓"

db-stop:
	@echo "Stopping PostgreSQL..."
	docker-compose stop postgres

wire:
	@echo "Generating Wire dependency injection code..."
	cd backend && go install github.com/google/wire/cmd/wire@latest && wire

swagger:
	@echo "Generating Swagger documentation..."
	cd backend && go install github.com/swaggo/swag/cmd/swag@latest && swag init -g cmd/api/main.go -o docs

backend-run:
	@echo "Running backend server..."
	cd backend && go run cmd/api/main.go
