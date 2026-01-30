.PHONY: all build test lint clean run docker-* mock swagger-setup tools run-migrations mocks help

# ==============================================================================
# Variables
# ==============================================================================

APP_NAME=salary_calculator
BUILD_DIR=bin
DOCKER_COMPOSE=docker compose
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
GOPATH=$(shell go env GOPATH)
TOOLS_DIR=$(CURDIR)/tools/bin
VENDOR_DIR=$(CURDIR)/vendor

# Database migration variables
MIGRATION_DIR=migrations
DB_STRING="user=$(DB_USERNAME) password=$(DB_PASSWORD) dbname=$(DB_DATABASE) sslmode=disable host=$(DB_HOST) port=$(DB_PORT)"

# Add tools/bin to PATH
export PATH := $(TOOLS_DIR):$(PATH)

# Load environment variables from .env file
include .env
export

# section: Help
# ==============================================================================
# Help
# ==============================================================================

.DEFAULT_GOAL := help

help: ## Show this help message
	@awk 'BEGIN { \
	  printf "\033[1;34m%-20s %-15s %s\033[0m\n", "Command", "Alias", "Description"; \
	  printf "\033[1;34m%-20s %-15s %s\033[0m\n", "-------", "-----", "-----------"; \
	} \
	/^# section: / { \
	  current_section = substr($$0, 12); \
	} \
	/^[a-zA-Z_-]+:.*## .*$$/ { \
	  split($$0, parts, ": .*## "); \
	  target = parts[1]; \
	  desc = parts[2]; \
	  if (desc ~ /^Alias for /) { \
	    split(desc, alias_parts, "Alias for "); \
	    main_target = alias_parts[2]; \
	    aliases[main_target] = target; \
	  } else { \
	    targets[target] = desc; \
	    sections[target] = current_section; \
	    section_order[current_section] = section_order[current_section] count + 1; \
	    command_order[current_section, ++section_count[current_section]] = target; \
	  } \
	} \
	END { \
	  for (sect in section_order) { \
	    printf "\n\033[1;36m%s\033[0m\n", sect; \
	    printf "\033[1;36m%s\033[0m\n", "================================"; \
	    for (i = 1; i <= section_count[sect]; i++) { \
	      t = command_order[sect, i]; \
	      alias = (t in aliases) ? aliases[t] : "-"; \
	      printf "\033[33m%-20s\033[0m \033[32m%-15s\033[0m %s\n", t, alias, targets[t]; \
	    } \
	  } \
	}' $(MAKEFILE_LIST)

# section: Build & Run
# ==============================================================================
# Build & Run
# ==============================================================================

build: ## Build the application
	@rm -rf $(BUILD_DIR)
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/main.go

run-local: ## Run the application locally
	@go run ./cmd/main.go

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -rf tmp/cache

br: build ## Build and run the application
	@echo "Starting application locally..."
	@./$(BUILD_DIR)/$(APP_NAME)

# section: Code Formatting
# ==============================================================================
# Code Formatting
# ==============================================================================

fmt-strict: ## Strict formatting with gofumpt and goimports
	@echo "Formatting Go code with gofumpt + goimports..."
	@command -v gofumpt >/dev/null 2>&1 || go install mvdan.cc/gofumpt@latest
	@command -v goimports >/dev/null 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
	@gofumpt -w $(GO_FILES)
	@goimports -w -local salary_calculator $(GO_FILES)
	@echo "Strict formatting done."

fmt-check: ## Check formatting without changes
	@command -v gofumpt >/dev/null 2>&1 || go install mvdan.cc/gofumpt@latest
	@command -v goimports >/dev/null 2>&1 || go install golang.org/x/tools/cmd/goimports@latest
	@echo "Checking formatting (gofumpt + goimports)..."
	@GF_OUT=$$(gofumpt -l $(GO_FILES)); \
	GI_OUT=$$(goimports -l -local salary_calculator $(GO_FILES)); \
	if [ -n "$$GF_OUT$$GI_OUT" ]; then \
	  echo "Following files need formatting:"; \
	  echo "$$GF_OUT"; \
	  echo "$$GI_OUT"; \
	  exit 1; \
	else \
	  echo "Formatting is OK"; \
	fi

# section: Dependencies
# ==============================================================================
# Dependencies
# ==============================================================================

deps-update: ## Update and tidy Go modules
	@go mod tidy
	@go mod verify
	@go mod vendor

deps-download: ## Download Go modules
	@go mod download

# section: Tools & Code Generation
# ==============================================================================
# Tools & Code Generation
# ==============================================================================

tools: ## Install development tools (goose, sqlc, mockgen)
	@echo "Deleting tools dir"
	@rm -rf $(TOOLS_DIR)
	@echo "Installing development tools"
	@mkdir -p $(TOOLS_DIR)
	@echo "Installing goose"
	@GOBIN=$(TOOLS_DIR) go install github.com/pressly/goose/v3/cmd/goose@latest
	@echo "Installing sqlc"
	@GOBIN=$(TOOLS_DIR) go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@echo "Installing mockgen"
	@GOBIN=$(TOOLS_DIR) go install github.com/golang/mock/mockgen@latest
	@echo "Tools installation completed"

check-env: ## Check if .env file exists
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found"; \
		echo "Please copy .env.example to .env and configure it"; \
		exit 1; \
	fi

mocks: ## Generate mocks
	@if [ ! -f "$(TOOLS_DIR)/mockgen" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Generating mocks..."
	@PATH=$(TOOLS_DIR):$(PATH) go generate ./...

sqlc-gen: check-env ## Generate SQL code
	@if [ ! -f "$(TOOLS_DIR)/sqlc" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Generating SQL code..."
	@$(TOOLS_DIR)/sqlc generate
	@if [ -d "internal/generated" ]; then \
		git add internal/generated/; \
		echo "Added generated files to git"; \
	fi

# section: Docker
# ==============================================================================
# Docker
# ==============================================================================

run-docker: ## Start Docker services
	@echo "Starting docker services..."
	@$(DOCKER_COMPOSE) up -d

clean-docker: ## Clean up Docker containers, volumes, and networks
	@echo "Cleaning up Docker services..."
	@$(DOCKER_COMPOSE) down --volumes --remove-orphans

restart-docker: ## Restart Docker services (stop and start)
	@echo "Restarting Docker services..."
	@$(DOCKER_COMPOSE) down
	@$(DOCKER_COMPOSE) up -d

full-reset: ## Full reset: clean Docker, restart services, and apply fresh migrations
	$(MAKE) clean-docker
	$(MAKE) restart-docker
	@sleep 3
	$(MAKE) run-migrations

# section: Database Migrations
# ==============================================================================
# Database Migrations
# ==============================================================================

run-migrations: check-env ## Apply database migrations
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Applying migrations..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) up

rollback-migrations: check-env ## Rollback last migration
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Rolling back last migration..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) down

create-migration: ## Create a new migration (requires name=...)
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

migrate-reset: check-env
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Rolling back all migrations..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) reset
	@echo "Migrations reset completed"

migrate-fresh: check-env ## Reset and reapply all migrations
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Rolling back all migrations..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) reset
	@echo "Migrations reset completed"
	@echo "Applying migrations..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) up
	@echo "Migrations applied"

migrate-status: check-env ## Check migration status
	@if [ ! -f "$(TOOLS_DIR)/goose" ]; then \
		$(MAKE) tools; \
	fi
	@echo "Checking migration status..."
	@$(TOOLS_DIR)/goose -dir $(MIGRATION_DIR) postgres $(DB_STRING) status

# section: Aliases
# ==============================================================================
# Aliases
# ==============================================================================

# Build & Run aliases
b: build ## Alias for build
c: clean ## Alias for clean
rl: run-local ## Alias for run-local

# Code formatting aliases
f: fmt ## Alias for fmt
fs: fmt-strict ## Alias for fmt-strict
fc: fmt-check ## Alias for fmt-check

# Dependencies aliases
du: deps-update ## Alias for deps-update
dd: deps-download ## Alias for deps-download

# Tools & Code Generation aliases
t: tools ## Alias for tools
ce: check-env ## Alias for check-env
mk: mocks ## Alias for mocks

# Docker aliases
rd: run-docker ## Alias for run-docker
cd: clean-docker ## Alias for clean-docker
rsd: restart-docker ## Alias for restart-docker
fr: full-reset ## Alias for full-reset

# Database migration aliases
mc: create-migration ## Alias for create-migration
rm: run-migrations ## Alias for run-migrations
rbm: rollback-migrations ## Alias for rollback-migrations
mf: migrate-fresh ## Alias for migrate-fresh
ms: migrate-status ## Alias for migrate-status
mres: migrate-reset

# SQL aliases
sqg: sqlc-gen ## Alias for sqlc-gen
