---
name: review-pr
description: Use after opening a slice PR to review it against SPEC.md and CLAUDE.md. Merges automatically if no CRITICALs found.
---

# review-pr

## Overview

Reviews the current open PR against project conventions and acceptance criteria. Auto-merges with squash if clean; blocks on CRITICALs.

## Steps

Work through these in order. Do not skip any.

### 1. Fetch PR details

```bash
gh pr view --json number,title,body,files,commits
```

Note the PR number, title, and which files changed.

### 2. Read conventions and spec

Read `SPEC.md` and `CLAUDE.md` — note the acceptance criteria relevant to this slice.

### 3. Check every changed file

For each `.go` file in the PR diff, check:

**CRITICAL (blocks merge):**
- `return err` without `fmt.Errorf("functionName: %w", err)` wrapping
- Business logic inside `internal/handler/` (DB calls, HMAC, NATS publishing)
- DB access (`SELECT`, `INSERT`, `UPDATE`, `db.`, `sqlx.`) outside `internal/repository/`
- Exported function without a corresponding test in `_test.go`
- New import not present in `go.mod`

**WARNINGS (note but do not block):**
- Error message starting with capital letter or ending with punctuation
- Missing `context.Context` threading in functions that hit DB or network
- Magic strings/numbers that should be constants
- Functions longer than ~40 lines

### 4. Check acceptance criteria coverage

List which criteria from SPEC.md this PR satisfies. Flag any criterion that the PR claims to cover but the implementation does not fully satisfy.

### 5. Run tests

```bash
cd d:/githook && go test ./... 2>&1
```

Any `FAIL` or `panic` line is a CRITICAL finding.

### 6. Write findings to REVIEW.md

Append a dated block:

```
---
## Code Review — Slice N — YYYY-MM-DD

### CRITICAL (must fix)
- [ ] `path/file.go:line` — description

### WARNINGS (should fix)
- [ ] `path/file.go:line` — description

### OK
- what was clean

---
```

### 7. Decide: merge or block

**If CRITICAL section is empty (`None`):**

```bash
gh pr merge --squash --delete-branch
```

Then confirm: "PR merged and branch deleted."

**If CRITICALs exist:**

List each CRITICAL item and stop. Do not merge. Wait for the author to fix and push before re-running this skill.
