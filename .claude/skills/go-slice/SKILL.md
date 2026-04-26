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
3. **Branch, commit, push, and open PR**:

```bash
git checkout -b slice-N/short-kebab-description
git add .
git commit -m "slice N: <description>"
git push -u origin slice-N/short-kebab-description
gh pr create \
  --title "slice N: <description>" \
  --body "$(cat <<'EOF'
## What was built
<one paragraph: what behaviour this slice adds and why it matters to the overall system>

## Changes made
- `path/to/file.go` — <one line: what this file does or what was changed>
- `path/to/file_test.go` — <one line: what the test covers>
(list every file created or modified)

## Architecture decisions
<explain any non-obvious decisions made in this slice — naming, structure, error handling approach — and reference ADR.md by number where relevant (e.g. "uses sqlx per ADR-0003")>

## Acceptance criteria covered
<paste the relevant numbered criteria from SPEC.md verbatim, e.g.:>
> AC #2: A request with a missing or invalid X-Hub-Signature-256 returns 401 before any DB or NATS interaction.

## How to test
\`\`\`bash
go test ./internal/handler/... -run TestWebhookHandler -v
# or: make test
\`\`\`
<add any curl commands or manual steps if applicable>

## What's next
Slice N+1 — <name of next slice> — depends on <what from this slice the next one builds on>
EOF
)" \
  --base master
```

Determine the slice number from the row just marked ✅ in `SLICES.md`. Branch name is kebab-case. Every section of the PR body must contain real content derived from the actual implementation — no placeholder text left unfilled.

**If `gh` is not available or the command fails**, fall back to:
```bash
git push -u origin slice-N/short-kebab-description
```
Then print: `PR: https://github.com/alibazlamit/githook/compare/slice-N/short-kebab-description`

Examples:

```
branch:  slice-1/health-endpoint
commit:  slice 1: GET /health handler, route registered, TestHealth passing

branch:  slice-3/signature-validation
commit:  slice 3: HMAC-SHA256 signature validation, returns 401 on invalid or missing sig

branch:  slice-5/webhook-ingest
commit:  slice 5: webhook ingest — DB insert then NATS publish, 409 on duplicate delivery_id
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
- [ ] Branch `slice-N/description` created and checked out
- [ ] `git commit` run with `slice N: <description>` message
- [ ] Branch pushed to origin with `git push -u origin slice-N/description`
- [ ] PR opened with `gh pr create` (or fallback URL printed if `gh` unavailable)
