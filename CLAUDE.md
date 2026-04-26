# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**githook** is a GitHub webhook ingestion service. It receives webhook events over HTTP, validates them, deduplicates via Postgres, publishes to NATS JetStream, and processes them with worker goroutines that write results back to Postgres.

## Stack

- **Go 1.23**
- **NATS JetStream** — event bus between HTTP ingestion and workers
- **PostgreSQL** via `sqlx` (raw SQL only, no ORM)
- **Prometheus** — metrics

## Data Flow

```
GitHub → HTTP handler → validate X-Hub-Signature-256
                      → check X-GitHub-Delivery idempotency in Postgres
                      → publish to NATS JetStream
                          → worker goroutines consume
                          → write results to Postgres
```

## Commands

```bash
make build          # compile all packages
make test           # run all tests
make lint           # run golangci-lint
make vet            # static analysis
make check          # vet + lint + test

go test ./... -run TestFoo  # run a single test by name
```

## Conventions

### Error Handling
All errors must be wrapped with context using `fmt.Errorf("functionName: %w", err)`. Never swallow errors.

### Database
- Raw SQL via `sqlx` only — no ORM, no query builders.
- All queries must use parameterized placeholders (`$1`, `$2`, …), never string concatenation.
- The **repository layer** is the only code allowed to touch the DB directly.

### Architecture Boundaries
- **HTTP handlers** must not contain business logic — they validate input, call services, and write responses.
- **Repository layer** owns all DB access.
- **Service/worker layer** sits between handlers/workers and repositories.

### Testing
- Every public function must have a corresponding test in a `_test.go` file in the same package.
- When a test fails, fix the production code — never modify tests to make them pass.
- Use table-driven tests for functions with multiple input/output cases.

### Dependencies
Do not add any dependency not already in `go.mod` without asking first.

### Commits
Each slice ends with a git commit. Format: `slice N: description` (e.g. `slice 1: GET /health handler and test`). After marking a slice done in SLICES.md, always remind the user to commit before moving to the next slice.

## Do Not

- Use an ORM.
- Add dependencies not in `go.mod` without explicit approval.
- Implement anything beyond what is asked for in each step.
- Put business logic in HTTP handlers.

## Current Task

[ Update this before each prompt ]
