# Telegram Semantic Search Bot Makefile

# Variables
BINARY_NAME=telegram-search-bot
MAIN_PATH=./main.go
BUILD_DIR=./bin

# Default target
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run         - Run the bot in development mode"
	@echo "  make build       - Build the binary"
	@echo "  make clean       - Clean build artifacts"
	@echo "  make deps        - Download and tidy dependencies"
	@echo "  make test        - Run tests"
	@echo "  make fmt         - Format code"
	@echo "  make lint        - Run linter (requires golangci-lint)"
	@echo "  make check       - Run fmt, lint, and test"
	@echo "  make setup       - Initial setup (deps + build)"

# Development
.PHONY: run
run: deps
	@echo "Running bot in development mode..."
	@echo "Make sure TELEGRAM_TOKEN is set in your environment"
	go run $(MAIN_PATH)

.PHONY: build
build: deps fmt
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Binary created: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Code quality
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

.PHONY: lint
lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run

.PHONY: test
test:
	@echo "Running tests..."
	go test -v ./...

.PHONY: check
check: fmt lint test
	@echo "All checks passed!"

# Setup and cleanup
.PHONY: setup
setup: deps build
	@echo "Setting up .env file..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from template"; \
		echo ""; \
		echo "⚠️  IMPORTANT: Edit .env file and set your TELEGRAM_TOKEN"; \
		echo "   1. Get token from @BotFather on Telegram"; \
		echo "   2. Edit .env file: TELEGRAM_TOKEN=your_actual_token"; \
		echo ""; \
	else \
		echo ".env file already exists"; \
	fi
	@echo "Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "1. Edit .env file and set your TELEGRAM_TOKEN"
	@echo "2. Run the bot: make run"
	@echo "3. Add bot to a Telegram chat and send /start"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f messages.db
	go clean

# Database management
.PHONY: clean-db
clean-db:
	@echo "Removing database file..."
	rm -f messages.db

.PHONY: backup-db
backup-db:
	@if [ -f messages.db ]; then \
		cp messages.db messages_backup_$(shell date +%Y%m%d_%H%M%S).db; \
		echo "Database backed up"; \
	else \
		echo "No database file found"; \
	fi

# Development helpers
.PHONY: dev
dev: clean-db run

.PHONY: logs
logs:
	@echo "Following bot logs (if running with systemd or docker)..."
	@echo "For local development, logs appear in the terminal where you ran 'make run'"

# Docker support (optional)
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run:
	@echo "Running in Docker..."
	@echo "Make sure to set TELEGRAM_TOKEN environment variable"
	docker run --rm -e TELEGRAM_TOKEN=$(TELEGRAM_TOKEN) $(BINARY_NAME)

# Release
.PHONY: release
release: clean check build
	@echo "Release build complete!"
	@echo "Binary available at: $(BUILD_DIR)/$(BINARY_NAME)"