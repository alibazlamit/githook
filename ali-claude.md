# Ali's Claude Code Playbook

A personal reference for starting and developing any software project with Claude Code professionally. Use Part 1 when bootstrapping a new project. Reference Part 2 during development.

---

## Part 1 — Project Bootstrap Checklist

Create these files in order before writing a single line of implementation code.

---

### 1. CLAUDE.md

**Why**: Claude's operating manual for the project. Every session loads this file first. Without it, Claude makes assumptions about your stack and conventions — with it, every session starts aligned.

**Must contain**: stack, build commands, hard constraints, error handling convention, architecture boundaries, current task placeholder.

````markdown
# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

[One paragraph: what this service does and why it exists]

## Stack

- **[Language + version]**
- **[Primary database]** via [driver] — [constraint, e.g., raw SQL only, no ORM]
- **[Message broker / queue]** — [purpose]
- **[Metrics / observability]**

## Data Flow

```
[ASCII diagram of the happy path through your system]
[e.g.: Client → HTTP handler → validate → DB → queue → worker → DB]
```

## Commands

```bash
make build    # compile
make test     # run all tests
make lint     # run linter
make vet      # static analysis
make check    # vet + lint + test

[lang] test ./... -run TestFoo  # run a single test by name
```

## Conventions

### Error Handling
All errors wrapped with `fmt.Errorf("functionName: %w", err)`. Never swallow errors.

### Database
- Raw SQL only — no ORM, no query builders
- Parameterized placeholders only — never string concatenation
- Repository layer is the only code that touches the DB directly

### Architecture Boundaries
- HTTP handlers: validate input, call services, write responses — no business logic
- Repository layer: all DB access
- Service/worker layer: sits between handlers and repositories

### Testing
- Every public function has a corresponding test in a `_test.go` file in the same package
- When a test fails, fix the production code — never modify tests to make them pass
- Use table-driven tests for functions with multiple input/output cases

### Dependencies
Do not add any dependency not already in the module file without asking first.

## Do Not

- Use an ORM
- Add unapproved dependencies
- Implement beyond what is asked in each step
- Put business logic in HTTP handlers

## Current Task

[ Update this before each prompt ]
````

---

### 2. SPEC.md

**Why**: The contract between intent and implementation. Claude refers to this when implementing every feature. Without it, behaviour drifts — endpoints return the wrong status codes, DB columns get the wrong types, edge cases get missed. Write this before any code exists.

**Must contain**: all endpoints with status codes and body shapes, full DB schema with column types and constraints, acceptance criteria written as concrete testable conditions (not vague requirements).

````markdown
# SPEC.md — [Service Name] API Contract & Data Model

## Endpoints

### POST /[resource]

[One sentence description]

**Request headers:**

| Header | Required | Description |
|---|---|---|
| `Content-Type` | yes | `application/json` |
| `[Custom-Header]` | yes | [purpose] |

**Responses:**

| Status | Condition |
|---|---|
| `200 OK` | [exact success condition] |
| `400 Bad Request` | [missing header or malformed body] |
| `401 Unauthorized` | [auth failure condition] |
| `409 Conflict` | [duplicate/idempotency condition] |
| `500 Internal Server Error` | [infrastructure failure] |

**Success response body:**
```json
{ "status": "accepted" }
```

**Error response body:**
```json
{ "error": "<human-readable message>" }
```

---

### GET /health

**Response:** `200 OK`
```json
{ "status": "ok" }
```

---

### GET /metrics

Prometheus metrics in text format.

**Response:** `200 OK`, `Content-Type: text/plain; version=0.0.4`

---

## Database Schema

### Table: `[table_name]`

```sql
CREATE TABLE [table_name] (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    [field]     TEXT        NOT NULL UNIQUE,
    status      TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ
);
```

**`status` values:**

| Value | Meaning |
|---|---|
| `pending` | [initial state after creation] |
| `processed` | [terminal success state] |
| `failed` | [terminal failure state] |

---

## Acceptance Criteria

| # | Criterion |
|---|---|
| 1 | [Concrete, testable condition — names the input and expected output exactly] |
| 2 | [Concrete, testable condition] |
| 3 | [Concrete, testable condition] |
````

---

### 3. ADR.md

**Why**: Records *why* you made architectural choices, not just what they are. Six months later, or when Claude is asked to change something, this prevents wheel-reinvention and re-litigating settled decisions.

**Must contain**: the decision, 2+ alternatives with *specific* rejection reasons (not "it's worse"), honest consequences including negative ones.

````markdown
# Architecture Decision Records

| ADR | Title | Status | Date |
|-----|-------|--------|------|
| [0001](#adr-0001-decision-title) | [Decision Title] | accepted | YYYY-MM-DD |

---

## ADR-0001: [Decision Title]

**Date**: YYYY-MM-DD
**Status**: accepted

### Context

[2-4 sentences: what problem forced this decision, what constraints existed, when the decision was made]

### Decision

[1-2 sentences: what was chosen and at what scope]

### Alternatives Considered

**[Alternative 1 Name]**
- **Pros**: [genuine benefits]
- **Cons**: [genuine drawbacks]
- **Why not**: [the specific reason this was rejected for this project]

**[Alternative 2 Name]**
- **Pros**: [genuine benefits]
- **Cons**: [genuine drawbacks]
- **Why not**: [the specific reason this was rejected for this project]

### Consequences

**Positive**
- [concrete benefit]

**Negative**
- [honest trade-off]

**Risks**
- [risk and how it is mitigated]
````

---

### 4. Makefile

**Why**: Single source of truth for build, test, and lint. Claude runs `make check` as its self-verification step before claiming work is done. Every project should have a `check` target that is the definition of "ready to merge."

```makefile
.PHONY: build test lint vet check

build:
	go build ./...

test:
	go test ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

check: vet lint test
```

*Adapt to your stack — swap in `npm run build`, `cargo build`, `pytest`, `bundle exec rspec`, etc.*

---

### 5. Linting Config

**Why**: Defines "clean" mechanically so Claude doesn't guess. Issues caught here never reach code review.

**Go — `.golangci.yml`:**

```yaml
linters:
  enable:
    - errcheck       # unhandled errors
    - gosimple       # simplification suggestions
    - govet          # suspicious constructs
    - ineffassign    # unused assignments
    - staticcheck    # static analysis
    - unused         # unused code
    - gofmt          # formatting
    - goimports      # import ordering
    - revive         # opinionated linter
    - noctx          # http requests without context
    - bodyclose      # unclosed http response bodies
    - sqlclosecheck  # unclosed sql rows/statements

linters-settings:
  revive:
    rules:
      - name: exported
      - name: error-return
      - name: error-strings
      - name: error-naming
      - name: var-naming

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - exhaustruct
```

**TypeScript — `eslint.config.js`:**

```js
export default [
  { rules: { "no-unused-vars": "error", "eqeqeq": "error", "no-console": "warn" } }
];
```

**Python — `pyproject.toml` (ruff):**

```toml
[tool.ruff]
select = ["E", "F", "I", "N", "UP"]
line-length = 88
```

---

### 6. Module / Package Init

**Why**: Locks the module path and dependency graph before Claude touches dependencies. Do this before adding any library.

```bash
# Go
go mod init github.com/[org]/[repo]

# Node
npm init -y

# Python
python -m venv .venv && pip install -e ".[dev]"

# Rust
cargo init
```

---

## Part 2 — Development Workflow

### Vertical Slices

A vertical slice is a thin, end-to-end working increment of exactly one behaviour. It cuts through every layer of the stack (HTTP → service → DB) and is independently testable and deployable.

**How to size a slice:**
- One user-facing behaviour, fully working
- Small enough to complete in one Claude session (roughly 1–3 hours of prompts)
- Has a clear "done" condition traceable to a row in SPEC.md's acceptance criteria

**Good slice:** `POST /webhook/github returns 401 for invalid signatures — handler + validation + test`

**Bad slice:** `Implement the HTTP layer` — too wide, no single behaviour is complete at the end

**Slice ordering rule:** Build the slice that unblocks the most other work first.

For backend HTTP services, a reliable default order:
1. `GET /health` — smoke test that the server starts
2. DB schema + migration — everything else depends on this
3. Happy path write endpoint — core value delivered
4. Idempotency / deduplication — correctness on the write path
5. Auth / validation — security layer
6. Read endpoints — observability and querying
7. Error cases and edge conditions

---

### Session Structure

Each Claude session is a single focused unit of work on one slice.

**Opening a session:**
1. Update `## Current Task` in `CLAUDE.md` to name the slice: `Implementing slice 3 — HMAC signature validation`
2. Give Claude a one-sentence context in your first prompt: *"We're implementing slice 3 — HMAC-SHA256 signature validation on POST /webhook/github. See SPEC.md criterion 2 and CLAUDE.md for conventions."*

**During a session:**
- One slice per session — resist "while we're here" scope creep
- After each significant change, run `make check` — fix failures before moving forward
- When Claude says something is done, ask it to run the tests before accepting the claim

**Closing a session:**
1. Run `make check` — all targets green before stopping
2. Reset `## Current Task` to `[ Update this before each prompt ]`
3. Commit with a message naming the slice: `feat: HMAC signature validation for POST /webhook/github`

---

### Planning Before Coding

**Write a plan first when:**
- The slice touches more than 2 files
- You are unsure of the right implementation shape
- The feature has non-obvious ordering constraints or edge cases

**How to ask for a plan:**
> *"Before writing any code, write a step-by-step implementation plan for [slice]. List every file to create or modify, the order of changes, and flag any decision points where you are uncertain."*

Review the plan against CLAUDE.md (conventions), SPEC.md (behaviour), and ADR.md (architecture choices). Push back on anything that deviates. Then say "go ahead."

**Skip the plan when:**
- The slice is a single function with a clear signature
- You are adding a test for an already-specified, already-implemented behaviour

---

### The Review Loop

Before accepting any Claude output as done:

1. **Run `make check`** — do not accept code that fails lint or tests
2. **Read the diff** — Claude is fast but not infallible; read every changed line
3. **Check against SPEC.md** — does the behaviour match the acceptance criteria exactly?
4. **Fix code, not tests** — if a test fails, the production code is wrong; never modify a test to make it pass
5. **Be specific when pushing back** — *"The signature check happens after the DB insert — per SPEC.md criterion 2 it must happen before any DB or NATS interaction."* Vague feedback produces vague fixes.
