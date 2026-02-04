# =========================
# CONFIG
# =========================
SHELL := /bin/bash
ENV_FILE := .env
ENT_DIR := ent
MIGRATIONS_DIR := migrations

# default env
ATLAS_ENV ?= local

# =========================
# HELP
# =========================
.PHONY: help
help:
	@echo ""
	@echo "Available commands:"
	@echo "  make gen               Generate ent client"
	@echo "  make migrate-hash      Generate atlas.sum (once)"
	@echo "  make migrate-diff name=init"
	@echo "  make migrate-apply     Apply migrations"
	@echo "  make setup             First time project setup"
	@echo "  make migrate-local"
	@echo "  make migrate-staging"
	@echo "  make migrate-prod"
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
# ENT
# =========================
.PHONY: gen
gen:
	@echo ">> Generating Ent client..."
	@$(LOAD_ENV) go generate ./$(ENT_DIR)

# =========================
# ATLAS
# =========================
.PHONY: migrate-hash
migrate-hash:
	@echo ">> Hashing migration directory..."
	@$(LOAD_ENV) atlas migrate hash

.PHONY: migrate-diff
migrate-diff:
ifndef name
	$(error Usage: make migrate-diff name=your_migration_name)
endif
	@echo ">> Creating migration diff: $(name)"
	@$(LOAD_ENV) atlas migrate diff $(name) --env $(ATLAS_ENV)

.PHONY: migrate-apply
migrate-apply:
	@echo ">> Applying migrations (env=$(ATLAS_ENV))"
	@$(LOAD_ENV) atlas migrate apply --env $(ATLAS_ENV)

# =========================
# SHORTCUTS
# =========================
.PHONY: migrate-local migrate-staging migrate-prod
migrate-local:
	@$(MAKE) migrate-apply ATLAS_ENV=local

migrate-staging:
	@$(MAKE) migrate-apply ATLAS_ENV=staging

migrate-prod:
	@$(MAKE) migrate-apply ATLAS_ENV=production

# =========================
# FIRST TIME SETUP
# =========================
.PHONY: setup
setup:
	@echo ">> First time setup"
	@$(MAKE) gen
	@$(MAKE) migrate-hash
	@$(MAKE) migrate-apply
