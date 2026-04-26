# Slices — githook

Tracks vertical slice implementation progress. Each slice delivers one end-to-end behaviour, is independently testable, and maps to SPEC.md acceptance criteria.

## Progress

| # | Slice | Status | AC |
|---|-------|--------|-----|
| 0 | Project scaffolding — module, structure, migrations | ✅ Done | — |
| 1 | `GET /health` — server starts, health endpoint returns 200 | ✅ Done | — |
| 2 | DB schema — migration applied, table exists | ✅ Done | — |
| 3 | Signature validation — invalid sig returns 401 | ⬜ Todo | AC #2 |
| 4 | Idempotency — duplicate delivery returns 409 | ⬜ Todo | AC #1 |
| 5 | Happy path ingest — valid event stored + published to NATS | ⬜ Todo | AC #3, #4 |
| 6 | NATS worker — consumes messages, updates row status | ⬜ Todo | AC #5 |
| 7 | Recovery cron — republishes stuck `received` events | ⬜ Todo | AC #6, #7 |
| 8 | `GET /metrics` — Prometheus endpoint live | ⬜ Todo | — |

## Definition of Done

A slice is **done** when:
1. `make check` passes (vet + lint + test all green)
2. The behaviour matches the SPEC.md acceptance criterion exactly
3. Every new public function has a test in `_test.go`

**One slice per session.** Update `## Current Task` in `CLAUDE.md` before each session.
