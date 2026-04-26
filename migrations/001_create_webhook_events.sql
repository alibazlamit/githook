CREATE TABLE webhook_events (
    id           UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    delivery_id  TEXT        NOT NULL UNIQUE,
    event_type   TEXT        NOT NULL,
    payload      JSONB       NOT NULL,
    status       TEXT        NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at TIMESTAMPTZ
);

-- supports the recovery cron query: WHERE status = 'received' AND created_at < NOW() - INTERVAL '10 minutes'
CREATE INDEX idx_webhook_events_status_created ON webhook_events(status, created_at);
