Slice 0 — project scaffolding: module init, package stubs, migration SQL — 2026-04-26
Slice 1 — GET /health: handler wired, route registered, TestHealth passing — 2026-04-26
Slice 3 — signature validation: HMAC-SHA256 check, missing/invalid sig returns 401, 4 test cases — 2026-04-26
Slice 2 — DB schema: make migrate target added, migration SQL verified by test — 2026-04-26
Slice 3 — signature validation: HMAC-SHA256 check, missing/invalid sig returns 401, 4 test cases — 2026-04-26

---
## Code Review — Slice 3 — 2026-04-26

### CRITICAL (must fix)
- [ ] `internal/handler/webhook.go:45-54` — `validSignature` performs HMAC-SHA256 computation directly inside the handler package. CLAUDE.md explicitly states "HTTP handlers must not contain business logic." HMAC validation is business logic; it must be extracted to the service layer (e.g., `internal/service/webhook.go`) and called via the `webhookIngestor` interface or a dedicated validator interface. The handler should only call the service and write the response.

### WARNINGS (should fix)
- [ ] `internal/handler/webhook.go:59` — return value of `json.NewEncoder(w).Encode(...)` in `writeError` is silently dropped. While recovery is impossible after `WriteHeader` has been called, CLAUDE.md convention says "never swallow errors." Assign to `_` with a comment, or log the error.
- [ ] `internal/handler/webhook_test.go` — `writeError` is tested only indirectly through `ServeHTTP`. Since `writeError` is unexported the convention does not strictly require its own test, but the `Content-Type` header it sets (`application/json`) is never asserted in any test case; add a header assertion to at least one table row.
- [ ] `REVIEW.md` (on branch) — the Slice 3 log entry was inserted between Slice 1 and Slice 2 entries (line 3), breaking chronological order. Minor, but should be appended at the end.

### OK
- `internal/handler/webhook.go` — no bare `return err`, no DB access, no new external dependencies (all added imports are stdlib). Error from `io.ReadAll` is correctly handled with a 500 response. `hmac.Equal` is used for constant-time comparison, preventing timing attacks — correct.
- `internal/handler/webhook_test.go` — table-driven test covers all four meaningful cases (missing header, wrong prefix, wrong secret, valid sig). Uses `httptest` correctly. Test file is present for the package as required.
- No SQL, `sqlx`, or raw DB calls anywhere in `internal/handler/`.
- No unwrapped `return err` statements anywhere in `internal/`.
- No new third-party dependencies introduced; `go.mod` unchanged.
- Test suite: all packages green (`ok github.com/ali/githook/internal/handler`, `ok github.com/ali/githook/migrations`).

---
