# Default target
.PHONY: default
default: help

include .env

PG_URL := $(MARIADB_DSN)

# Ensure the .env file exists
$(ENV_FILE):
	@echo "$(ENV_FILE) not found. Please create it and populate with required environment variables."
	@exit 1

# Print help message
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo ""
	@echo "generate-server : Generate server code from OpenAPI spec"
	@echo "run : run the application"
	@echo "lint : Run linters on the codebase"



.PHONY: generate-server
generate-server:
	@oapi-codegen -package handler -generate gin-server,types api/swagger.yml > internal/adapter/handler/server.gen.go
	@go mod tidy

.PHONY: run
run: generate-server
	@echo "Running the application..."
	@go run cmd/main.go
	@echo "Application is running on port 8080"

.PHONY: lint
lint:
	@golangci-lint run ./... --fix
	@fieldalignment -fix ./...