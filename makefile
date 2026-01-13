# ===============================
# CONFIG
# ===============================
APP_NAME := go-boilerplate
MIGRATE_BIN := migrate
MIGRATIONS_DIR := migrations
DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/go_boilerplate?sslmode=disable

# ===============================
# MIGRATION COMMANDS
# ===============================

.PHONY: migrate-create
migrate-create:
	@read -p "Migration name: " name; \
	$(MIGRATE_BIN) create -ext sql -dir $(MIGRATIONS_DIR) $$name

.PHONY: migrate-up
migrate-up:
	$(MIGRATE_BIN) -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) up

.PHONY: migrate-down
migrate-down:
	$(MIGRATE_BIN) -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) down 1

.PHONY: migrate-force
migrate-force:
	@read -p "Force version: " version; \
	$(MIGRATE_BIN) -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) force $$version

.PHONY: migrate-version
migrate-version:
	$(MIGRATE_BIN) -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) version

.PHONY: migrate-drop
migrate-drop:
	$(MIGRATE_BIN) -database "$(DATABASE_URL)" -path $(MIGRATIONS_DIR) drop -f

# ===============================
# DEV HELPERS
# ===============================

.PHONY: migrate-reset
migrate-reset:
	$(MAKE) migrate-drop
	$(MAKE) migrate-up
