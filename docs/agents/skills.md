# Available Skills

Catalog of Claude Code skills available in this environment and when to reach for each on Bibseller. These are user-level plugins, not repo files - this doc just tracks what's available and how it maps to our workflow. Re-check against the live skill list occasionally; entries here can go stale if skills are added, removed, or renamed upstream.

## Bibseller-specific

| Skill | What it does |
| --- | --- |
| `/ship-to-prod` | Promotes a tested PR (or `main`) to the live site: squash-merge to `main`, then `make promote` + `make prod-migrate` + `make prod-up`. Not for staging (`make staging-up`). |

## Issue tracker & planning

| Skill | What it does |
| --- | --- |
| `/to-prd` | Turns the current conversation into a PRD and publishes it as a GitHub issue. |
| `/to-issues` | Breaks a plan/PRD into independently-grabbable issues (tracer-bullet vertical slices). |
| `/triage` | Moves issues through the five-role triage state machine (see `docs/agents/triage-labels.md`). |
| `/grill-me` | Interviews you against a plan/design until every branch is resolved. |
| `/grill-with-docs` | Same, but also updates `CONTEXT.md`/ADRs inline as decisions crystallize - see `docs/agents/domain.md`. |
| `/handoff` | Compacts the current conversation into a handoff doc for another agent. |

## Quality gates

Map to CLAUDE.md's "Goal-Driven Execution" (`make verify` / `make smoke` before behavior-touching PRs).

| Skill | What it does |
| --- | --- |
| `/code-review` | Reviews the current diff for correctness bugs and simplification/efficiency cleanups. `--comment` posts inline PR comments, `--fix` applies findings. `ultra` runs a multi-agent cloud review. |
| `/security-review` | Security review of the pending changes on the current branch. |
| `/simplify` | Applies reuse/simplification/efficiency cleanups to changed code - quality only, no bug hunting. |
| `/verify` | Runs the app and confirms a change actually works, not just that tests pass. |
| `/tdd` | Red-green-refactor loop for building features or fixing bugs test-first. |
| `/run` | Launches/screenshots the app to see a change working end to end. |
| `/diagnose` | Disciplined reproduce -> minimise -> hypothesise -> instrument -> fix loop for hard bugs and perf regressions. |

## Simplicity enforcement (ponytail suite)

Maps to CLAUDE.md's "Simplicity First" section - same YAGNI/stdlib-first bias, enforced as a mode or a scan instead of prose.

| Skill | What it does |
| --- | --- |
| `/ponytail` | Forces the laziest solution that actually works: stdlib/native/existing-dependency before new code, no unrequested abstractions. Persists as a mode (`lite`/`full`/`ultra`) until turned off. |
| `/ponytail-review` | Reviews a diff for over-engineering only: reinvented stdlib, unneeded deps, speculative abstractions. Complements `/code-review`, doesn't replace it. |
| `/ponytail-audit` | Whole-repo scan for over-engineering - ranked list of what to delete/simplify/replace with stdlib. One-shot report, no fixes applied. |
| `/ponytail-debt` | Harvests `ponytail:` shortcut comments left in the code into a debt ledger, so deliberate simplifications get tracked instead of forgotten. |
| `/ponytail-help` | Quick-reference card for ponytail modes and commands. |

## Codebase & review

| Skill | What it does |
| --- | --- |
| `/init` | Initializes or refreshes a CLAUDE.md from the current codebase. |
| `/improve-codebase-architecture` | Finds deepening/refactoring opportunities, informed by `CONTEXT.md` and `docs/adr/`. |
| `/review` | Reviews a GitHub PR by number (vs. `/code-review` for the local working diff). |

## Meta / tooling

| Skill | What it does |
| --- | --- |
| `/find-skills` | Discovers and installs other skills. |
| `/write-a-skill` | Authors a new skill with proper structure. |
| `/fewer-permission-prompts` | Scans transcripts for common read-only tool calls, adds an allowlist to `.claude/settings.json`. |
| `/loop` | Runs a prompt/command on a recurring interval (self-paced if no interval given). |
| `/schedule` | Creates/manages scheduled cloud agents (cron-based routines). |
| `/update-config` | Configures the Claude Code harness itself - hooks, permissions, env vars, settings.json. |
| `/keybindings-help` | Customizes keyboard shortcuts / chord bindings. |
| `/claude-api` | Reference for the Claude API / Anthropic SDK (models, pricing, tool use, streaming, caching). |
| `/obsidian-vault` | Searches/creates/organizes notes in an Obsidian vault - unrelated to this repo's own docs. |
| `/prototype` | Builds a throwaway prototype (terminal app or multiple UI variations) to sanity-check a design before committing to it. |
| `/caveman` | Ultra-compressed communication mode - terse prose, same technical content. |
