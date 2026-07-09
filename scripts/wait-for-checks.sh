#!/usr/bin/env bash
# Wait for a commit's required CI checks to succeed, so `make promote` does not
# push to the check-protected production ref while main's checks are still
# running (#213). Polls the check-runs API for each required check's latest run:
#   - all success            -> exit 0 (promote proceeds)
#   - any concluded non-success -> exit 1 (fail fast, do not deploy)
#   - still running / not yet created -> keep polling until the deadline
#
# Usage: wait-for-checks.sh <commit-sha>
# Escape hatch: PROMOTE_SKIP_CHECKS=1 skips the wait (gh down, known-green, etc.)
set -euo pipefail

sha="${1:?usage: wait-for-checks.sh <commit-sha>}"
required="backend frontend smoke"
timeout_s="${PROMOTE_CHECK_TIMEOUT:-900}" # 15 min
interval_s=20

if [ "${PROMOTE_SKIP_CHECKS:-}" = "1" ]; then
	echo "PROMOTE_SKIP_CHECKS=1 - skipping the required-checks wait for ${sha:0:7}" >&2
	exit 0
fi

deadline=$(($(date +%s) + timeout_s))
echo "Waiting for required checks (${required// /, }) on ${sha:0:7}..." >&2

while :; do
	# Latest run per check name (reruns produce several; take the most recent by
	# started_at). Tolerate a transient gh/network blip by retrying next tick.
	summary=$(gh api "repos/{owner}/{repo}/commits/$sha/check-runs" --paginate 2>/dev/null |
		jq -r --argjson want "$(printf '%s\n' $required | jq -R . | jq -s .)" '
			[.check_runs[] | select(.name as $n | $want | index($n))]
			| group_by(.name)
			| map(sort_by(.started_at) | last | {name, status, conclusion})
			| .[] | "\(.name)\t\(.status)\t\(.conclusion // "")"' 2>/dev/null || true)

	done_ok=0
	failed=""
	pending=""
	for name in $required; do
		line=$(printf '%s\n' "$summary" | awk -F'\t' -v n="$name" '$1==n {print; exit}')
		if [ -z "$line" ]; then
			pending="$pending $name(absent)"
			continue
		fi
		status=$(printf '%s' "$line" | cut -f2)
		conclusion=$(printf '%s' "$line" | cut -f3)
		if [ "$status" != "completed" ]; then
			pending="$pending $name($status)"
		elif [ "$conclusion" = "success" ]; then
			done_ok=$((done_ok + 1))
		else
			failed="$failed $name($conclusion)"
		fi
	done

	if [ -n "$failed" ]; then
		echo "ERROR: required check(s) did not pass on ${sha:0:7}:$failed" >&2
		echo "Not promoting. Inspect: gh run list --commit $sha" >&2
		exit 1
	fi
	if [ "$done_ok" -eq "$(printf '%s\n' $required | wc -w)" ]; then
		echo "All required checks green on ${sha:0:7}." >&2
		exit 0
	fi
	if [ "$(date +%s)" -ge "$deadline" ]; then
		echo "ERROR: timed out after ${timeout_s}s waiting for checks on ${sha:0:7}:$pending" >&2
		echo "Not promoting. Re-run when green, or set PROMOTE_SKIP_CHECKS=1 to override." >&2
		exit 1
	fi
	echo "  still waiting:$pending (${done_ok}/$(printf '%s\n' $required | wc -w) green)" >&2
	sleep "$interval_s"
done
