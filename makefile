# =========================
# CONFIG
# =========================
SHELL := /bin/bash
ENV_FILE := .env
MIGRATIONS_DIR := migrations

# =========================
# HELP
# =========================
.PHONY: help
help:
	@echo ""
	@echo "Available commands:"
	@echo "  make setup             First time project setup"
	@echo "  make dev               Start development server"
	@echo ""

# =========================
# ENV LOADER
# =========================
define LOAD_ENV
	set -a; \
	if [ -f $(ENV_FILE) ]; then \
		. $(ENV_FILE); \
	fi; \
	set +a;
endef

# =========================
# DEVELOPMENT
# =========================
.PHONY: dev
dev:
	@echo ">> Starting development server..."
	@$(LOAD_ENV) air

# =========================
# FIRST TIME SETUP
# =========================
.PHONY: setup
setup:
	@echo ">> First time setup"
	@echo ">> Installing dependencies..."
	@go mod tidy
	@echo ">> Setup complete! Run 'make dev' to start the server."
