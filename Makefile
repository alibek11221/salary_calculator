.PHONY: all help build b run-local rl clean c fmt f fmt-strict fs fmt-check fc deps-update du deps-download dd tools t mocks mk sqlc-gen sqg swg swagger run-docker rd clean-docker cd restart-docker rsd full-reset fr run-migrations rm rollback-migrations rbm create-migration mc migrate-reset mres migrate-fresh mf migrate-status ms br

# ==============================================================================
# Variables
# ==============================================================================

APP_NAME=salary_calculator
BUILD_DIR=bin
DOCKER_COMPOSE=docker compose
GO=go
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
TOOLS_DIR=$(CURDIR)/tools/bin
VENDOR_DIR=$(CURDIR)/vendor

# Database migration variables
MIGRATION_DIR=migrations
DB_STRING="user=$(DB_USERNAME) password=$(DB_PASSWORD) dbname=$(DB_DATABASE) sslmode=disable host=$(DB_HOST) port=$(DB_PORT)"

# Load environment variables from .env file
-include .env

# ==============================================================================
# Help
# ==============================================================================

.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^[a-zA-Z\-_0-9]+ [a-zA-Z\-_0-9 ]*:.*## / { \
		helpMessage = match($$0, /## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$0, 0, index($$0, ":")-1); \
			helpMessage = substr($$0, RSTART + 3, RLENGTH); \
			printf "  \033[36m%-20s\033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	/^## [a-zA-Z &]+/ { \
		printf "\n\033[1m%s\033[0m\n", substr($$0, 4); \
	}' $(MAKEFILE_LIST)

# ==============================================================================
# Build & Run
# ==============================================================================

## Build & Run

build b: ## Build the application
	@rm -rf $(BUILD_DIR)
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go

run-local rl: ## Run the application locally
	@go run ./cmd/main.go

clean c: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp/cache

br: build ## Build and run the application
	@echo "Starting application locally..."
	@./$(BUILD_DIR)/$(APP_NAME)

# ==============================================================================
# Code Formatting
# ==============================================================================

## Formatting

fmt f: fmt-strict ## Format code (alias for fmt-strict)

fmt-strict fs: ## Strict formatting with gofumpt and goimports
	@echo "Formatting Go code with gofumpt + goimports..."
	@if [ ! -f "$(TOOLS_DIR)/gofumpt" ] || [ ! -f "$(TOOLS_DIR)/goimports" ]; then \
		$(MAKE) tools; \
	fi
	@$(TOOLS_DIR)/gofumpt -w $(GO_FILES)
	@$(TOOLS_DIR)/goimports -w -local salary_calculator $(GO_FILES)
	@echo "Strict formatting done."

fmt-check fc: ## Check formatting without changes
	@if [ ! -f "$(TOOLS_DIR)/gofumpt" ] || [ ! -f "$(TOOLS_DIR)/goimports" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Checking formatting (gofumpt + goimports)..."
	@GF_OUT=$$($(TOOLS_DIR)/gofumpt -l $(GO_FILES)); \
	GI_OUT=$$($(TOOLS_DIR)/goimports -l -local salary_calculator $(GO_FILES)); \
	if [ -n "$$GF_OUT$$GI_OUT" ]; then \
	  echo "Following files need formatting:"; \
	  echo "$$GF_OUT"; \
	  echo "$$GI_OUT"; \
	  exit 1; \
	else \
	  echo "Formatting is OK"; \
	fi

# ==============================================================================
# Dependencies
# ==============================================================================

## Dependencies

deps-update du: ## Update and tidy Go modules
	@go mod tidy
	@go mod verify
	@go mod vendor

deps-download dd: ## Download Go modules
	@go mod download

# ==============================================================================
# Tools & Code Generation
# ==============================================================================

## Tools & Generation

tools t: ## Install development tools (goose, sqlc, mockgen, swag, gofumpt, goimports)
	@echo "Deleting tools dir"
	@rm -rf $(TOOLS_DIR)
	@echo "Installing development tools"
	@mkdir -p $(TOOLS_DIR)
	@GOBIN=$(TOOLS_DIR) go install github.com/pressly/goose/v3/cmd/goose@latest
	@GOBIN=$(TOOLS_DIR) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@GOBIN=$(TOOLS_DIR) go install github.com/golang/mock/mockgen@latest
	@GOBIN=$(TOOLS_DIR) go install github.com/swaggo/swag/cmd/swag@latest
	@GOBIN=$(TOOLS_DIR) go install mvdan.cc/gofumpt@latest
	@GOBIN=$(TOOLS_DIR) go install golang.org/x/tools/cmd/goimports@latest
	@echo "Tools installation completed"

check-env ce: ## Check if .env file exists
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		echo "Please copy .env.example to .env and configure it"; \
		exit 1; \
	fi

mocks mk: ## Generate mocks
	@if [ ! -f "$(TOOLS_DIR)/mockgen" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Generating mocks..."
	@PATH=$(TOOLS_DIR):$(PATH) go generate ./...

sqlc-gen sqg: check-env ## Generate SQL code
	@if [ ! -f "$(TOOLS_DIR)/sqlc" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Generating SQL code..."
	@$(TOOLS_DIR)/sqlc generate
	@if [ -d "internal/generated" ]; then \
		git add internal/generated/; \
		echo "Added generated files to git"; \
	fi

swagger swg: ## Generate Swagger documentation
	@if [ ! -f "$(TOOLS_DIR)/swag" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Generating Swagger documentation..."
	@$(TOOLS_DIR)/swag init -g cmd/main.go --parseDependency --parseInternal

# ==============================================================================
# Docker
# ==============================================================================

## Docker

run-docker rd: ## Start Docker services
	@echo "Starting docker services..."
	@$(DOCKER_COMPOSE) up -d

clean-docker cd: ## Clean up Docker services
	@echo "Cleaning up Docker services..."
	@$(DOCKER_COMPOSE) down --volumes --remove-orphans

restart-docker rsd: ## Restart Docker services
	@echo "Restarting Docker services..."
	@$(DOCKER_COMPOSE) down
	@$(DOCKER_COMPOSE) up -d

full-reset fr: ## Full reset (Docker + Migrations)
	$(MAKE) clean-docker
	$(MAKE) restart-docker
	@sleep 3
	$(MAKE) run-migrations

# ==============================================================================
# Database Migrations
# ==============================================================================

## Migrations

run-migrations rm: check-env ## Apply database migrations
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Applying migrations..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) up

rollback-migrations rbm: check-env ## Rollback last migration
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Rolling back last migration..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) down

create-migration mc: ## Create a new migration (requires name=...)
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@if [ -z "$(name)" ]; then \
		echo "Error: migration name not specified. Use 'make create-migration name=<migration_name>'"; \
		exit 1; \
	fi
	@echo "Creating migration $(name)..."
	@OUTPUT=$$($(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) create $(name) sql 2>&1); \
	echo "$$OUTPUT"; \
	MIGRATION_FILE=$$(find $(MIGRATION_DIR) -name "*.sql" -type f -printf '%T@ %p\n' | sort -n | tail -1 | cut -d' ' -f2-); \
	if [ -n "$$MIGRATION_FILE" ] && [ -f "$$MIGRATION_FILE" ]; then \
		git add "$$MIGRATION_FILE"; \
		echo "Added $$MIGRATION_FILE to git"; \
	fi

migrate-reset mres: check-env ## Rollback all migrations
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Rolling back all migrations..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) reset

migrate-fresh mf: check-env ## Reset and reapply all migrations
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Rolling back all migrations..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) reset
	@echo "Applying migrations..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) up

migrate-status ms: check-env ## Check migration status
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Checking migration status..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) status
