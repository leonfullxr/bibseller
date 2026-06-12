SHELL := /bin/bash
-include .env
export

DATABASE_URL ?= postgres://postgres:dev@localhost:5432/bibseller?sslmode=disable

# Heavyweight tools run via `go run pkg@version`: pinned here, cached by Go,
# and kept OUT of backend/go.mod — sqlc's dependency tree would otherwise
# force go-directive bumps that break linters (see docs/ARCHITECTURE.md).
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.27.1
SQLC  := go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

.PHONY: help infra dev migrate migrate-down sqlc sqlc-check seed \
        test test-backend test-frontend lint lint-backend lint-frontend

help: ## list targets
	@grep -E '^[a-z-]+:.*##' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-15s %s\n", $$1, $$2}'

infra: ## start Postgres + MinIO + Mailpit
	docker compose up -d --wait

dev: infra ## full dev loop: infra + Go API (:8080) + SvelteKit (:5173)
	@trap 'kill 0' INT TERM; \
	( cd backend && { command -v air >/dev/null && air || go run ./cmd/api; } ) & \
	( cd frontend && npm run dev ) & \
	wait

migrate: ## apply goose migrations
	cd backend && $(GOOSE) -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down: ## roll back one migration
	cd backend && $(GOOSE) -dir db/migrations postgres "$(DATABASE_URL)" down

sqlc: ## regenerate type-safe query code
	cd backend && $(SQLC) generate

sqlc-check: ## fail if committed sqlc output is stale (no-op until M1 adds queries)
	@if ls backend/db/queries/*.sql >/dev/null 2>&1; then \
		cd backend && $(SQLC) generate && git diff --exit-code -- internal/platform/db/sqlcgen; \
	else echo "sqlc: no queries yet (M1, issue #2) — skipping"; fi

seed: ## wipe + load dev seed data (dev-only)
	cd backend && go run ./cmd/seed

test: test-backend test-frontend ## run all tests
test-backend:
	cd backend && go vet ./... && go test ./...
test-frontend:
	cd frontend && npm run test

lint: lint-backend lint-frontend ## lint both halves
lint-backend:
	@if command -v golangci-lint >/dev/null; then cd backend && golangci-lint run; \
	else echo "golangci-lint not installed — skipping (CI runs it)"; fi
lint-frontend:
	cd frontend && npm run lint
