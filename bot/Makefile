.PHONY: help install run dev test lint format check clean logs setup

# Variables
PYTHON_VERSION = 3.12
VENV_PATH = $(shell poetry env info --path 2>/dev/null)

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Setup and Installation
setup: ## Setup the project (install poetry if needed and dependencies)
	@echo "Setting up the project..."
	@command -v poetry >/dev/null 2>&1 || { echo "Installing poetry..."; curl -sSL https://install.python-poetry.org | python3 -; }
	@poetry install
	@echo "Project setup complete!"

install: ## Install dependencies
	@echo "Installing dependencies..."
	@poetry install
	@echo "Dependencies installed!"

# Development
run: ## Run the bot
	@echo "Starting the bot..."
	@poetry run python src/main.py

dev: ## Run the bot in development mode with auto-reload
	@echo "Starting bot in development mode..."
	@poetry run python -u src/main.py

# Code Quality
lint: ## Run linter (ruff check)
	@echo "Running linter..."
	@poetry run ruff check .

lint-fix: ## Run linter with auto-fix
	@echo "Running linter with auto-fix..."
	@poetry run ruff check --fix .

format: ## Format code with ruff
	@echo "Formatting code..."
	@poetry run ruff format .

format-check: ## Check if code is formatted correctly
	@echo "Checking code formatting..."
	@poetry run ruff format --check .

check: ## Run all code quality checks
	@echo "Running all code quality checks..."
	@make lint
	@make format-check
	@echo "All checks passed!"

fix: ## Fix all auto-fixable issues
	@echo "Fixing all auto-fixable issues..."
	@make lint-fix
	@make format
	@echo "All fixes applied!"

# Testing
test: ## Run tests
	@echo "Running tests..."
	@poetry run pytest

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	@poetry run pytest -v

test-watch: ## Run tests in watch mode
	@echo "Running tests in watch mode..."
	@poetry run pytest-watch

# Utilities
clean: ## Clean cache and temporary files
	@echo "Cleaning cache and temporary files..."
	@find . -type d -name "__pycache__" -exec rm -rf {} + 2>/dev/null || true
	@find . -type d -name ".pytest_cache" -exec rm -rf {} + 2>/dev/null || true
	@find . -type d -name ".ruff_cache" -exec rm -rf {} + 2>/dev/null || true
	@find . -name "*.pyc" -delete 2>/dev/null || true
	@echo "Cleanup complete!"

logs: ## Show bot logs
	@echo "Showing bot logs..."
	@tail -f /tmp/bot.log 2>/dev/null || echo "No log file found at /tmp/bot.log"

logs-follow: ## Follow bot logs in real-time
	@echo "Following bot logs..."
	@tail -f /tmp/bot.log 2>/dev/null || echo "No log file found at /tmp/bot.log"

# Environment
env-info: ## Show environment information
	@echo "Environment Information:"
	@echo "Poetry version: $(shell poetry --version 2>/dev/null || echo 'Not installed')"
	@echo "Python version: $(shell python --version 2>/dev/null || echo 'Not found')"
	@echo "Virtual env path: $(VENV_PATH)"
	@echo "Project dependencies:"
	@poetry show --tree 2>/dev/null || echo "Run 'make install' first"

shell: ## Open a shell in the virtual environment
	@echo "Opening shell in virtual environment..."
	@poetry shell

# Git helpers
git-hooks: ## Install git hooks using the pre-commit framework
	@echo "Installing git hooks via pre-commit..."
	@poetry run pre-commit install
	@echo "Git hooks installed!"

# All-in-one commands
start: ## Quick start (install + run)
	@make install
	@make run

dev-setup: ## Complete development setup
	@make setup
	@make git-hooks
	@echo "Development environment ready!"

