package service

import "context"

type WebhookService interface {
	Ingest(ctx context.Context, deliveryID, eventType string, payload []byte) error
}
