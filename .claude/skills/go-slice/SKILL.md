---
name: go-slice
description: Use when implementing any slice of the githook Go service — repository, handler, service, config, worker, or cron layers.
---

# go-slice

## Overview

Implements one vertical slice of the githook service end-to-end. A slice delivers exactly one behaviour, nothing more.

## Before Writing Any Code

Read these three files in order:

1. `CLAUDE.md` — stack, conventions, current task
2. `SPEC.md` — endpoints, schema, acceptance criteria
3. `SLICES.md` — which slice is active and what "done" means

## Rules

**Scope**
- Implement only what the current slice requires
- Do not touch files outside the slice boundary
- Do not add features "while you're there"

**Error handling**
- Every error must be wrapped: `fmt.Errorf("functionName: %w", err)`
- Never swallow errors with `_` or empty catch blocks

**Tests**
- Every new public function gets a corresponding test in a `_test.go` file in the same package
- When a test fails, fix the production code — never modify the test

**Dependencies**
- Do not add any package not already in `go.mod`
- If a slice genuinely requires a new dependency, stop and ask before proceeding

## Completing a Slice

When the slice is done and `make check` passes:

1. **Update `SLICES.md`** — change `⬜ Todo` to `✅ Done` for the completed slice
2. **Append to `REVIEW.md`** — one line: `Slice N — <what was built> — <date>`
3. **Commit** — stage all changes and commit:

```bash
git add .
git commit -m "slice N: <description of what was actually implemented>"
```

Determine the slice number from the row just marked ✅ in `SLICES.md`. The description must name the behaviour added, not just say "implement slice". Examples:

```
slice 1: GET /health handler, route registered, TestHealth passing
slice 3: HMAC-SHA256 signature validation, returns 401 on invalid or missing sig
slice 5: webhook ingest — DB insert then NATS publish, 409 on duplicate delivery_id
```

`REVIEW.md` format:
```
Slice 0 — project scaffolding: module init, package stubs, migration — 2026-04-26
Slice 1 — GET /health: handler wired, route registered, test passing — 2026-04-26
```

## Definition of Done

A slice is done when ALL of the following are true:

- [ ] `make check` passes (vet + lint + test)
- [ ] Behaviour matches the SPEC.md acceptance criterion for this slice exactly
- [ ] Every new public function has a test in `_test.go`
- [ ] `SLICES.md` updated
- [ ] One line appended to `REVIEW.md`
- [ ] `git commit` run with `slice N: <description>` message
