-include .env
export

MIGRATIONS_DIR := $(shell pwd)/migrations

.PHONY: init-swagger run-main init-services migrate-database init run dev help run-wire

init-swagger:
	swag init -d "./internal/router,./internal/controller,./internal/model,./internal/dto" -g "router.go" -o ./docs

run-wire:
	wire ./internal/di/wire.go
	
run-main:
	go run ./cmd/api/main.go

init-services:
	docker compose up -d
	
migrate-up:
	docker run --rm -v $(MIGRATIONS_DIR):/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(MONGO_MIGRATION_URL)" up

migrate-down:
	docker run --rm -v $(MIGRATIONS_DIR):/migrations \
		--network host \
		migrate/migrate \
		-path=/migrations \
		-database "$(MONGO_MIGRATION_URL)" down 1
	
init: init-services migrate-up init-swagger
	@echo "Initialization complete!"

run: init-swagger run-main

dev: init run