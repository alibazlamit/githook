# SPEC.md — githook API Contract & Data Model

## Endpoints

### POST /webhook/github

Receives a GitHub webhook event.

**Request headers:**

| Header | Required | Description |
|---|---|---|
| `X-Hub-Signature-256` | yes | HMAC-SHA256 of the raw body, prefixed `sha256=` |
| `X-GitHub-Delivery` | yes | UUID identifying this delivery (idempotency key) |
| `X-GitHub-Event` | yes | Event type (e.g. `push`, `pull_request`) |
| `Content-Type` | yes | Must be `application/json` |

**Responses:**

| Status | Condition |
|---|---|
| `200 OK` | Event accepted, deduplicated, published to NATS, and recorded in Postgres |
| `400 Bad Request` | Missing required header or malformed JSON body |
| `401 Unauthorized` | `X-Hub-Signature-256` does not match |
| `409 Conflict` | `X-GitHub-Delivery` already exists in `webhook_events` |
| `500 Internal Server Error` | Postgres or NATS failure |

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

Liveness check.

**Response:** `200 OK`

```json
{ "status": "ok" }
```

---

### GET /metrics

Prometheus metrics exposition in text format.

**Response:** `200 OK`, `Content-Type: text/plain; version=0.0.4`

Standard Go runtime metrics plus application-level counters (to be defined as implementation progresses).

---

## Database Schema

### Table: `webhook_events`

```sql
CREATE TABLE webhook_events (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    delivery_id  TEXT        NOT NULL UNIQUE,
    event_type   TEXT        NOT NULL,
    payload      JSONB       NOT NULL,
    status       TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ
);
```

**`status` values:**

| Value | Meaning |
|---|---|
| `received` | Stored after initial ingestion; not yet processed by a worker |
| `processed` | Worker completed successfully |
| `failed` | Worker exhausted all NATS redelivery attempts |

---

## Acceptance Criteria

| # | Criterion |
|---|---|
| 1 | A duplicate `X-GitHub-Delivery` returns `409` without publishing to NATS or creating a new DB row |
| 2 | A request with a missing or invalid `X-Hub-Signature-256` returns `401` before any DB or NATS interaction |
| 3 | A valid event is persisted to `webhook_events` with `status = received` before NATS publish is attempted. The DB row is the durable source of truth — an event that exists in the DB is never permanently lost regardless of what happens to NATS. |
| 4 | If NATS publish fails after a successful Postgres insert, the handler returns `500`. The row remains with `status = received` and will be recovered by the cron job (criterion 6). |
| 5 | Worker failures result in NATS JetStream redelivery; the row status is only updated to `processed` or `failed` once the worker reaches a terminal outcome. |
| 6 | A recovery cron job runs every 5 minutes. It scans for rows where `status = received` AND `created_at < NOW() - INTERVAL '10 minutes'` and re-publishes them to NATS JetStream. This guarantees eventual processing of any event stored but never successfully published. |
| 7 | Workers are idempotent: on receiving a NATS message, the worker checks the current `status` of the row before processing. If `status` is already `processed`, the worker acks the message and returns without reprocessing. This handles re-publication by the cron job of an already-processed event. |
