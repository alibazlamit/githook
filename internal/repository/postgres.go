package repository

import (
	"context"
	"time"
)

// postgresWebhookRepository is the Postgres-backed implementation of WebhookRepository.
// Real DB wiring is completed in slice 5; stubs return safe zero values until then.
type postgresWebhookRepository struct{}

// NewPostgresWebhookRepository constructs a postgresWebhookRepository.
// The db parameter will be wired in slice 5; it is intentionally unused here.
func NewPostgresWebhookRepository() WebhookRepository {
	return &postgresWebhookRepository{}
}

func (r *postgresWebhookRepository) Insert(_ context.Context, _ *WebhookEvent) error {
	// implemented in slice 5
	return nil
}

func (r *postgresWebhookRepository) UpdateStatus(_ context.Context, _, _ string, _ *time.Time) error {
	// implemented in slice 5
	return nil
}

func (r *postgresWebhookRepository) FindStuck(_ context.Context, _ time.Duration) ([]*WebhookEvent, error) {
	// implemented in slice 5
	return nil, nil
}

func (r *postgresWebhookRepository) DeliveryExists(_ context.Context, _ string) (bool, error) {
	// stub — real SQL query added in slice 5
	return false, nil
}
