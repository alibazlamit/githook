package repository

import (
	"context"
	"time"
)

const (
	StatusReceived  = "received"
	StatusProcessed = "processed"
	StatusFailed    = "failed"
)

type WebhookEvent struct {
	ID          string     `db:"id"`
	DeliveryID  string     `db:"delivery_id"`
	EventType   string     `db:"event_type"`
	Payload     []byte     `db:"payload"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	ProcessedAt *time.Time `db:"processed_at"`
}

type WebhookRepository interface {
	Insert(ctx context.Context, event *WebhookEvent) error
	UpdateStatus(ctx context.Context, deliveryID, status string, processedAt *time.Time) error
	FindStuck(ctx context.Context, olderThan time.Duration) ([]*WebhookEvent, error)
}
