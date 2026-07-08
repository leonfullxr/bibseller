SHELL := /bin/bash
-include .env
export

# 54320: dev Postgres publishes off the default port so dev tooling can never
# collide with the prod stack's 127.0.0.1:5432 on the self-host box (#159).
DATABASE_URL ?= postgres://postgres:dev@localhost:54320/bibseller?sslmode=disable

# Heavyweight tools run via `go run pkg@version`: pinned here, cached by Go,
# and kept OUT of backend/go.mod - sqlc's dependency tree would otherwise
# force go-directive bumps that break linters (see docs/ARCHITECTURE.md).
GOOSE := go run github.com/pressly/goose/v3/cmd/goose@v3.27.1
SQLC  := go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1

.PHONY: help infra dev migrate migrate-down sqlc sqlc-check seed \
        test test-backend test-frontend lint lint-backend lint-frontend \
        verify smoke prod-up prod-build prod-down prod-logs prod-migrate prod-backup \
        staging-up staging-down staging-logs staging-migrate staging-seed \
        prod-backup-offsite prod-restore-drill promote

help: ## list targets
	@grep -E '^[a-z-]+:.*##' $(MAKEFILE_LIST) | awk -F':.*## ' '{printf "  %-15s %s\n", $$1, $$2}'

infra: ## start Postgres + MinIO + Mailpit
	docker compose up -d --wait
	@# Stamp the server as dev infrastructure: `make seed` refuses to wipe a
	@# database without the marker, and prod/staging never get one (#159).
	@docker compose exec -T db psql -q -U postgres -d bibseller \
		-c "CREATE TABLE IF NOT EXISTS dev_marker (stamped_at timestamptz NOT NULL DEFAULT now())"

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

sqlc-check: ## fail if committed sqlc output is stale
	cd backend && $(SQLC) generate && git diff --exit-code -- internal/platform/db/sqlcgen

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
	else echo "golangci-lint not installed - skipping (CI runs it)"; fi
lint-frontend:
	cd frontend && npm run lint

verify: ## pre-commit gate: lint + typecheck + tests + sqlc drift (docs/CONTEXT.md -> D6)
	$(MAKE) lint
	cd frontend && npm run check
	$(MAKE) test
	$(MAKE) sqlc-check
	@echo "VERIFY PASSED"

smoke: ## end-to-end assertions against the seeded local stack (wipes dev data)
	./scripts/smoke.sh

# --- Production: self-host (docs/DEPLOYMENT.md -> Model A) -------------------
# Secrets live in deploy/.env.prod (gitignored; copy from .env.prod.example).
PROD_COMPOSE := docker compose --env-file deploy/.env.prod -f deploy/compose.prod.yml

prod-up: ## build + start the self-host prod stack (needs deploy/.env.prod)
	@test -f deploy/.env.prod || { echo "ERROR: deploy/.env.prod missing (copy from deploy/.env.prod.example)"; exit 1; }
	@if grep -qE 'replace-with|\.example' deploy/.env.prod; then \
	  echo "ERROR: deploy/.env.prod still has placeholder values to fill:"; \
	  grep -nE 'replace-with|\.example' deploy/.env.prod | sed -E 's/=.*/=<PLACEHOLDER>/'; \
	  exit 1; \
	fi
	$(PROD_COMPOSE) up -d --build

prod-build: ## build prod images without starting them (mixed expand/contract deploys - docs/DEPLOYMENT.md)
	$(PROD_COMPOSE) build

prod-down: ## back up offsite if the stack is up, then stop the prod stack (volumes kept)
	@if $(PROD_COMPOSE) ps --services --status running | grep -qx db; then \
	  echo "prod db is up - taking an offsite backup before stopping (#12)"; \
	  $(MAKE) prod-backup-offsite || echo "WARNING: pre-stop offsite backup failed; stopping anyway (the kept db volume still holds the data)" >&2; \
	else \
	  echo "prod db not running - nothing to back up; stopping"; \
	fi
	$(PROD_COMPOSE) down

prod-logs: ## tail prod logs
	$(PROD_COMPOSE) logs -f --tail=100

prod-migrate: ## apply goose migrations to the prod DB (via the one-shot migrate container)
	$(PROD_COMPOSE) run --rm migrate

prod-backup: ## pg_dump the prod DB to ./backups (then copy OFFSITE)
	@mkdir -p backups
	$(PROD_COMPOSE) exec -T db sh -c 'pg_dump --no-owner -U "$$POSTGRES_USER" "$$POSTGRES_DB"' \
	  | gzip > backups/db-$$(date +%F-%H%M).sql.gz
	@echo "Wrote backups/. Copy it OFFSITE - the laptop is a single point of failure."

prod-backup-offsite: ## full offsite backup: dump + MinIO mirror to R2 (the cron job)
	./scripts/offsite-backup.sh

prod-restore-drill: ## restore the latest R2 dump into a throwaway DB and verify it
	./scripts/restore-drill.sh

promote: ## fast-forward production to the tested origin/main and push (run from the prod clone)
	@test "$$(git branch --show-current)" = production || { echo "ERROR: run from the prod clone, on the production branch"; exit 1; }
	git fetch origin
	@git merge-base --is-ancestor HEAD origin/main || { echo "ERROR: production is not behind origin/main - nothing to promote, or it diverged"; exit 1; }
	git merge --ff-only origin/main
	git push
	@echo "Promoted production -> $$(git rev-parse --short HEAD). Deploy in the order the release's migrations require: docs/DEPLOYMENT.md 'Migration ordering: expand/contract'."

# --- Staging: ephemeral parallel stack on the same box (docs/DEPLOYMENT.md) --
# The prod compose file + Caddyfile, isolated by a different compose project name
# (own network + pgdata/miniodata volumes) and deploy/.env.staging. Tracks `main`
# (prod tracks `production`); bring it up to test a release, then `staging-down`.
STAGING_COMPOSE := docker compose -p bibseller-staging --env-file deploy/.env.staging -f deploy/compose.prod.yml

staging-up: ## build + start the ephemeral staging stack (needs deploy/.env.staging)
	@test -f deploy/.env.staging || { echo "ERROR: deploy/.env.staging missing (copy from deploy/.env.staging.example)"; exit 1; }
	@if grep -qE 'replace-with|\.example' deploy/.env.staging; then \
	  echo "ERROR: deploy/.env.staging still has placeholder values to fill:"; \
	  grep -nE 'replace-with|\.example' deploy/.env.staging | sed -E 's/=.*/=<PLACEHOLDER>/'; \
	  exit 1; \
	fi
	$(STAGING_COMPOSE) up -d --build

staging-down: ## stop the staging stack (volumes are kept)
	$(STAGING_COMPOSE) down

staging-logs: ## tail staging logs
	$(STAGING_COMPOSE) logs -f --tail=100

staging-migrate: ## apply goose migrations to the staging DB
	$(STAGING_COMPOSE) run --rm migrate

staging-seed: ## load preview races into the staging DB (run after staging-migrate)
	$(STAGING_COMPOSE) exec -T db sh -c 'psql -U "$$POSTGRES_USER" -d "$$POSTGRES_DB"' < deploy/seed-preview-races.sql
