#!/usr/bin/env bash
# End-to-end smoke: boots API + frontend against a migrated+seeded local
# Postgres and asserts the behaviors that define the product — above all,
# the transfer-policy gating matrix. Wipes dev data (runs `make seed`).
#
# Usage: make smoke   (requires Postgres reachable; `make infra` or local PG)
set -u

API=http://localhost:8080
WEB=http://localhost:5173
FAILURES=0

say() { printf '%s\n' "$*"; }
pass() { say "  ✓ $1"; }
fail() { say "  ✗ $1"; FAILURES=$((FAILURES + 1)); }

expect_status() { # url want desc
	local got
	got=$(curl -s -o /dev/null -w '%{http_code}' "$1")
	[ "$got" = "$2" ] && pass "$3" || fail "$3 (got HTTP $got, want $2)"
}
expect_contains() { # url pattern desc
	curl -s "$1" | grep -q "$2" && pass "$3" || fail "$3 (pattern not found: $2)"
}
expect_not_contains() { # url pattern desc
	curl -s "$1" | grep -q "$2" && fail "$3 (forbidden pattern found: $2)" || pass "$3"
}

say "── prepare database (migrate + seed)"
if ! make -s migrate >/dev/null 2>&1 || ! make -s seed >/dev/null 2>&1; then
	say "✗ Postgres not reachable or migrate/seed failed — run 'make infra' first"
	exit 2
fi

say "── boot stack"
(cd backend && go build -o /tmp/bibseller-smoke-api ./cmd/api) || exit 2
/tmp/bibseller-smoke-api >/tmp/bibseller-smoke-api.log 2>&1 &
API_PID=$!
(cd frontend && npm run dev -- --strictPort >/tmp/bibseller-smoke-web.log 2>&1) &
WEB_PID=$!
trap 'kill $API_PID $WEB_PID 2>/dev/null; wait 2>/dev/null' EXIT

wait_for() {
	for _ in $(seq 1 60); do curl -sf "$1" >/dev/null 2>&1 && return 0; sleep 0.5; done
	return 1
}
wait_for "$API/api/healthz" || { say "✗ API never became healthy (see /tmp/bibseller-smoke-api.log)"; exit 2; }
wait_for "$WEB/" || { say "✗ frontend never came up (see /tmp/bibseller-smoke-web.log)"; exit 2; }

say "── api basics"
expect_status "$API/api/readyz" 200 "readyz reports ready (DB reachable)"
expect_status "$WEB/api/healthz" 200 "healthz reachable through the dev proxy"
expect_status "$API/api/v1/races?country=ESP" 400 "invalid filter params are rejected"
expect_status "$API/api/v1/races/madrid-marathon-2027" 404 "draft races are hidden (404)"

say "── catalog SSR"
expect_contains "$WEB/races" "Munich Marathon 2026" "browse page renders seeded races"
expect_contains "$WEB/races?q=valencia" "Valencia" "full-text search works through the stack"

say "── the policy matrix (the product's core guarantee)"
expect_contains "$WEB/races/munich-marathon-2026" "allows bib resale" "platform_sale: resale callout renders"
expect_contains "$WEB/races/valencia-marathon-2026" "Official transfer process" "official_only: link-out renders"
expect_contains "$WEB/races/berlin-marathon-2026" "restricts bib transfers" "connect_only: strong disclaimer renders"
expect_contains "$WEB/races/sevilla-marathon-2027" "not verified yet" "unknown: unverified warning renders"

MUNICH_LISTING=$(curl -s "$API/api/v1/races/munich-marathon-2026/listings" |
	python3 -c 'import json,sys; print(json.load(sys.stdin)["items"][0]["id"])' 2>/dev/null)
BERLIN_LISTING=$(curl -s "$API/api/v1/races/berlin-marathon-2026/listings" |
	python3 -c 'import json,sys; print(json.load(sys.stdin)["items"][0]["id"])' 2>/dev/null)
if [ -n "$MUNICH_LISTING" ] && [ -n "$BERLIN_LISTING" ]; then
	expect_contains "$WEB/listings/$MUNICH_LISTING" "Buy securely" "platform_sale listing shows the buy CTA"
	expect_not_contains "$WEB/listings/$BERLIN_LISTING" "Buy securely" "NO buy affordance outside platform_sale"
else
	fail "could not resolve listing ids from the API"
fi

say "──"
if [ "$FAILURES" -gt 0 ]; then
	say "SMOKE FAILED: $FAILURES assertion(s) red"
	exit 1
fi
say "SMOKE PASSED: all assertions green"
