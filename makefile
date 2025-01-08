# Makefile for running Goose migrations

# Variables
DB_DRIVER ?= postgres
DB_DSN ?= "user=postgres dbname=exchange_rates sslmode=disable password=postgres host=localhost port=5432"
MIGRATIONS_DIR ?= ./migrations

# Default target
start: migrate run

# Run Goose migrations
migrate:
	@echo "Running Goose migrations..."
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) up
	@echo "Migrations completed."

# Rollback the latest migration
rollback:
	@echo "Rolling back the latest migration..."
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) down
	@echo "Rollback completed."

# Reset the database (rollback all migrations)
reset:
	@echo "Resetting the database..."
	@goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(DB_DSN) reset
	@echo "Database reset completed."

run:
	go run cmd/main.go


.PHONY: start run migrate rollback reset 