# Migration commands
MIGRATIONS_DIR := db/migrations
PG_DSN ?= postgres://postgres:mysecretpassword@localhost:44228/postgres?sslmode=disable

GOOSE_BIN=$(LOCAL_BIN)/goose
$(GOOSE_BIN):
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.19.2

.PHONY: migrate-up
migrate-up: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) -allow-missing postgres "$(PG_DSN)" up

.PHONY: migrate-down
migrate-down: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "$(PG_DSN)" down

.PHONY: migrate-reset
migrate-reset: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "$(PG_DSN)" reset

.PHONY: migrate-generate
migrate-generate: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) create $(name) sql

.PHONY: migrate-status
migrate-status: $(GOOSE_BIN)
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) postgres "$(PG_DSN)" status
