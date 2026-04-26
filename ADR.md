# Architecture Decision Records

| ADR | Title | Status | Date |
|-----|-------|--------|------|
| [0001](#adr-0001-nats-jetstream-over-redis-streams) | NATS JetStream over Redis Streams | accepted | 2026-04-25 |
| [0002](#adr-0002-x-github-delivery-as-idempotency-key) | X-GitHub-Delivery as idempotency key | accepted | 2026-04-25 |
| [0003](#adr-0003-sqlx-over-pgx-directly) | sqlx over pgx directly | accepted | 2026-04-25 |
| [0004](#adr-0004-stdlib-nethttp-over-a-framework) | stdlib net/http over a framework | accepted | 2026-04-25 |

---

## ADR-0001: NATS JetStream over Redis Streams

**Date**: 2026-04-25
**Status**: accepted

### Context

githook ingests GitHub webhooks and must reliably hand them off to worker goroutines for processing. The broker must support at-least-once delivery, explicit acknowledgement, and automatic redelivery on worker failure. Redis is already a common dependency in many stacks; NATS JetStream is purpose-built for durable messaging. The decision was made before any infrastructure was provisioned.

### Decision

We use NATS JetStream as the message bus between the HTTP ingestion layer and the worker pool.

### Alternatives Considered

**Redis Streams**
- **Pros**: Often already in the stack; low operational overhead if Redis is already present; familiar to most Go developers
- **Cons**: Consumer group management is manual; no built-in per-consumer ack timeout or redelivery policy without extra tooling; streams are a secondary feature of Redis, not its core abstraction
- **Why not**: Redelivery and ack-timeout semantics require explicit XCLAIM/XAUTOCLAIM calls that push complexity into application code. JetStream handles this at the broker level.

**Apache Kafka**
- **Pros**: Industry standard for high-throughput event streaming; rich ecosystem
- **Cons**: Heavy operationally (ZooKeeper or KRaft, topic partitioning, consumer groups, offset management); gross overkill for a single-topic webhook ingestion service
- **Why not**: Complexity-to-value ratio is poor for this use case. Kafka is appropriate at orders-of-magnitude higher throughput.

### Consequences

**Positive**
- Redelivery, ack timeout, and max-delivery-attempt policies are declared in the JetStream stream config — no application-level retry logic needed
- Push consumers simplify worker code: receive a message, process it, ack or nack
- NATS server is a single statically-linked binary with no external dependencies

**Negative**
- Adds NATS as an infrastructure dependency that must be run locally and in CI
- Team members unfamiliar with NATS must learn JetStream stream/consumer concepts

**Risks**
- JetStream persistence is file-based; disk pressure on the NATS server could cause message loss. Mitigated by monitoring NATS storage metrics and setting appropriate retention limits.

---

## ADR-0002: X-GitHub-Delivery as Idempotency Key

**Date**: 2026-04-25
**Status**: accepted

### Context

GitHub may redeliver webhooks on network failure or timeout. Processing the same event twice could produce duplicate DB writes or double-trigger downstream logic. We need an idempotency key that is stable across redeliveries, available before the payload is parsed, and guaranteed unique per logical delivery attempt by GitHub.

### Decision

We use the `X-GitHub-Delivery` request header as the idempotency key, stored as `delivery_id` in `webhook_events` with a UNIQUE constraint. A 409 is returned if the key already exists.

### Alternatives Considered

**SHA-256 hash of the request body**
- **Pros**: No dependency on a GitHub-specific header; works for any webhook source
- **Cons**: Two distinct deliveries of identical payloads (e.g., a push to the same commit twice) would collide; computing the hash requires buffering the full body before validation
- **Why not**: Hash collisions on legitimate distinct events are a correctness bug, not just a performance issue.

**Server-generated UUID stored in a short-lived cache**
- **Pros**: Fully self-contained; no GitHub coupling
- **Cons**: Requires a shared cache with TTL management; UUID is assigned after receipt, so retries from GitHub arrive with no shared key — deduplication is impossible
- **Why not**: Breaks the fundamental requirement: the key must be the same across GitHub's redeliveries of the same event.

### Consequences

**Positive**
- The UNIQUE constraint on `delivery_id` makes deduplication a single DB insert with conflict detection — no separate lookup needed
- `X-GitHub-Delivery` is present in all GitHub webhook payloads and is documented as stable across redeliveries
- The delivery UUID is visible in GitHub's webhook delivery logs, aiding debugging

**Negative**
- Couples the idempotency mechanism to GitHub's header convention; a different webhook source would need its own strategy
- If GitHub ever changes the semantics of this header, deduplication breaks silently

**Risks**
- The UNIQUE index must be present before the service handles traffic. Mitigated by enforcing schema migrations as a deploy prerequisite.

---

## ADR-0003: sqlx over pgx Directly

**Date**: 2026-04-25
**Status**: accepted

### Context

The service needs to read and write PostgreSQL. Options range from raw `database/sql` (verbose struct scanning), `sqlx` (thin scanning layer over `database/sql`), `pgx` used directly (bypasses `database/sql` for lower-level control), and ORMs (explicitly excluded by project conventions). The query volume is low and the schema is simple; raw performance is not a differentiator.

### Decision

We use `sqlx` (`github.com/jmoiron/sqlx`) with the `pgx` stdlib adapter (`github.com/jackc/pgx/v5/stdlib`) as the underlying driver.

### Alternatives Considered

**pgx directly (without database/sql)**
- **Pros**: Lower-level access to PostgreSQL-specific types; pgxpool handles connection pooling natively; slightly lower overhead for bulk operations
- **Cons**: Row scanning into structs requires manual field-by-field assignment or reflection code that `sqlx` already provides; pgx's API diverges from `database/sql`, making future driver swaps harder
- **Why not**: The benefit (access to pgx-native types) is unused in this project. The cost (manual scanning boilerplate) is paid on every query.

**Raw database/sql**
- **Pros**: Zero dependencies beyond the stdlib and driver; maximum transparency
- **Cons**: Scanning result rows into structs requires explicit `rows.Scan(&field1, &field2, ...)` calls in strict column order, which is brittle and verbose
- **Why not**: `sqlx` adds exactly the struct scanning ergonomics we'd otherwise hand-roll, with no SQL hiding and no magic.

### Consequences

**Positive**
- `db.StructScan`, `db.Select`, and `db.Get` eliminate repetitive scanning boilerplate while keeping SQL explicit and readable
- Named parameters via `db.NamedExec` map directly to struct fields, reducing positional-argument bugs
- Using the pgx stdlib adapter means we can switch drivers without changing application code

**Negative**
- `sqlx` is an additional dependency (though stable and widely adopted)
- Struct tags (`db:"column_name"`) must be kept in sync with schema column names manually

---

## ADR-0004: stdlib net/http over a Framework

**Date**: 2026-04-25
**Status**: accepted

### Context

githook exposes three HTTP endpoints: `POST /webhook/github`, `GET /health`, and `GET /metrics`. The routing requirements are minimal — no path parameters, no middleware chains beyond signature validation. Go 1.22 added method-based routing and enhanced `ServeMux` patterns to the standard library, closing the gap with lightweight routers like chi.

### Decision

We use the Go standard library `net/http` package with no third-party HTTP framework.

### Alternatives Considered

**Gin**
- **Pros**: Fast router; large ecosystem; familiar to many Go developers; built-in JSON binding and validation helpers
- **Cons**: Adds a dependency with its own release cadence and CVE surface; Gin's context type (`*gin.Context`) leaks into handler signatures, coupling all handlers to Gin
- **Why not**: None of Gin's features are needed for three endpoints with no path parameters. The dependency cost exceeds the benefit.

**Echo**
- **Pros**: Clean API; good middleware support
- **Cons**: Same framework-coupling problem as Gin; Echo's middleware ordering is a common source of subtle bugs
- **Why not**: Same reasoning as Gin.

**chi**
- **Pros**: Composable, idiomatic; uses stdlib `net/http` handlers natively; much lighter than Gin or Echo
- **Cons**: Still an external dependency; adds a router abstraction on top of what stdlib now provides natively in Go 1.22+
- **Why not**: Go 1.22 `ServeMux` covers our routing needs without any external package.

### Consequences

**Positive**
- Zero additional dependencies for HTTP routing
- Handler signatures use `http.ResponseWriter` and `*http.Request` — no framework coupling, fully testable with `net/http/httptest`
- No framework upgrade path to manage

**Negative**
- Request body binding and error response helpers must be written by hand (small, but not zero cost)
- If the API grows significantly (path parameters, versioning, complex middleware), migrating to chi will require touching all handler registrations

**Risks**
- Scope creep: as endpoints are added, the temptation to introduce a framework grows. Revisit this ADR if the endpoint count exceeds ~10 or path parameters become necessary.
